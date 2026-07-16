import { google } from '@ai-sdk/google';
import { streamText, jsonSchema } from 'ai';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { StdioClientTransport } from '@modelcontextprotocol/sdk/client/stdio.js';
import path from 'path';

export const maxDuration = 60; // Allows longer execution time for tool calls

// Singleton client to avoid spawning a new Go process on every request
let mcpClient: Client | null = null;
let mcpTransport: StdioClientTransport | null = null;

async function getMcpClient() {
  if (mcpClient) {
    return mcpClient;
  }

  console.log('Initializing new MCP Client (spawning Go server)...');
  
  const mcpServerPath = path.resolve(process.cwd(), '../mcp-server');

  mcpTransport = new StdioClientTransport({
    command: 'go',
    args: ['run', './cmd/mcp-server/main.go'],
    options: {
      cwd: mcpServerPath,
      env: {
        ...process.env,
        // Add any DB env variables if your Go server needs them from the Node environment,
        // or ensure they are present in the system environment.
      }
    }
  });

  mcpClient = new Client(
    { name: 'nextjs-mcp-client', version: '1.0.0' },
    { capabilities: { tools: {} } }
  );

  await mcpClient.connect(mcpTransport);
  console.log('MCP Client connected to Go Server');
  return mcpClient;
}

export async function POST(req: Request) {
  try {
    const { messages } = await req.json();

    const client = await getMcpClient();
    
    // Fetch all tools from the Go Server
    const toolsResponse = await client.listTools();
    
    // Convert MCP tools format to Vercel AI SDK format
    const aiTools: Record<string, any> = {};
    for (const tool of toolsResponse.tools) {
      aiTools[tool.name] = {
        description: tool.description || `Execute ${tool.name}`,
        // Use Vercel AI SDK's built-in jsonSchema utility to wrap the MCP JSON Schema
        parameters: jsonSchema(tool.inputSchema),
        // The execute function will be called by the LLM when it wants to use the tool
        execute: async (args: any) => {
          console.log(`LLM called tool: ${tool.name} with args:`, args);
          const result = await client.callTool({
            name: tool.name,
            arguments: args,
          });
          console.log(`Tool ${tool.name} returned:`, result);
          return result.content;
        }
      };
    }

    // Call Gemini with the tools
    const result = streamText({
      model: google('gemini-1.5-pro'), // Or gemini-1.5-flash
      system: "You are a helpful AI assistant. You can use the provided tools to fetch data or perform actions. Always use tools when appropriate to get accurate and up-to-date information.",
      messages,
      tools: aiTools,
      maxSteps: 5, // Allow the model to call multiple tools in a row
    });

    return result.toDataStreamResponse();
  } catch (error) {
    console.error('Error in chat route:', error);
    return new Response(JSON.stringify({ error: 'Failed to process request' }), { status: 500 });
  }
}
