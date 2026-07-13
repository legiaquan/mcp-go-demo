FROM golang:1.25-alpine AS builder
WORKDIR /app

# Cài đặt các dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the Go application
COPY . .
RUN go build -o mcp-server ./cmd/mcp-server

# Môi trường chạy với Node.js để dùng MCP Inspector
FROM node:20-alpine
WORKDIR /app

# Copy file binary từ bước build
COPY --from=builder /app/mcp-server /app/mcp-server

# Mở port mặc định của Inspector
EXPOSE 6274

# Biến môi trường để web server bind ra ngoài container
ENV HOST=0.0.0.0

# Chạy inspector và bind vào mcp-server
CMD ["npx", "-y", "@modelcontextprotocol/inspector", "/app/mcp-server"]
