'use client';

import { useChat } from '@ai-sdk/react';

export default function Chat() {
  const { messages, input, handleInputChange, handleSubmit } = useChat();

  return (
    <div className="flex flex-col w-full max-w-2xl py-24 mx-auto stretch">
      <h1 className="text-2xl font-bold mb-8 text-center">MCP Client with Gemini</h1>
      
      {messages.map(m => (
        <div key={m.id} className="whitespace-pre-wrap mb-6 bg-gray-50 p-4 rounded-lg shadow-sm border border-gray-100">
          <strong className="text-blue-600 block mb-1">{m.role === 'user' ? '🧑 Bạn:' : '🤖 AI:'}</strong>
          <div className="text-gray-800">{m.content}</div>
          
          {/* Render tool calls */}
          {m.toolInvocations?.map((toolInvocation: any) => {
            const toolCallId = toolInvocation.toolCallId;
            return (
              <div key={toolCallId} className="text-gray-500 italic mt-2 text-sm border-l-2 border-blue-300 pl-2 bg-blue-50 p-2 rounded">
                ⚙️ Gọi công cụ: <span className="font-semibold">{toolInvocation.toolName}</span>
                {toolInvocation.state === 'result' ? (
                  <span className="text-green-600 font-semibold ml-2">✓ Xong</span>
                ) : (
                  <span className="text-blue-500 font-semibold ml-2 animate-pulse">...đang xử lý</span>
                )}
              </div>
            );
          })}
        </div>
      ))}

      <form 
        onSubmit={handleSubmit} 
        className="fixed bottom-0 w-full max-w-2xl p-2 mb-8 bg-white border border-gray-300 rounded shadow-xl flex items-center"
      >
        <input
          className="flex-grow p-2 focus:outline-none text-black"
          value={input}
          placeholder="Hỏi tôi bất cứ điều gì..."
          onChange={handleInputChange}
        />
        <button type="submit" className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded ml-2 transition-colors">
          Gửi
        </button>
      </form>
    </div>
  );
}
