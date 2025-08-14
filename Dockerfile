# Use Node.js 22 as base for runtime with latest Node.js pre-installed
FROM node:22-slim AS runtime-base

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y \
        bash \
        curl \
        git \
        procps \
        ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install GitHub CLI (gh)
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null && \
    apt-get update && \
    apt-get install -y gh && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# https://www.npmjs.com/package/@iflow-ai/iflow-cli
# Pre-install iFlow CLI using npm package
RUN npm install -g @iflow-ai/iflow-cli

# Install github-mcp-server CLI tool
RUN go install github.com/github/github-mcp-server/cmd/github-mcp-server@latest

# Use official Go 1.24.4 image for building
FROM golang:1.24.4-bullseye AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y git ca-certificates curl

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./
COPY cmd/ ./cmd/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iflow-action .

# Final stage - copy Go binary to Node.js runtime
FROM runtime-base

# Create a non-root user with proper home directory
RUN groupadd -g 1001 iflow && \
    useradd -r -u 1001 -g iflow -m -d /home/iflow iflow

# Set working directory
WORKDIR /github/workspace

# Copy the binary from builder stage
COPY --from=builder /app/iflow-action /usr/local/bin/iflow-action

# Make sure binary is executable
RUN chmod +x /usr/local/bin/iflow-action

# Create .iflow directory for the non-root user and set permissions
RUN mkdir -p /home/iflow/.iflow && \
    chown -R iflow:iflow /home/iflow/.iflow

# Switch to non-root user
USER iflow

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/iflow-action"]