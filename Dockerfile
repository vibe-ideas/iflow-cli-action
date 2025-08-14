# Use Ubuntu 22.04 as base image
FROM ubuntu:22.04 AS runtime-base

# Set noninteractive installation mode to avoid prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install runtime dependencies including Node.js
RUN apt-get update && \
    apt-get install -y \
        bash \
        curl \
        git \
        procps \
        ca-certificates \
        software-properties-common && \
    # Install Node.js 22
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    # Clean up
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

# Install Go for github-mcp-server
# Download and install Go
ENV GO_VERSION=1.23.2
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

# https://www.npmjs.com/package/@iflow-ai/iflow-cli
# Pre-install iFlow CLI using npm package
RUN npm install -g @iflow-ai/iflow-cli

# Install github-mcp-server CLI tool
RUN /usr/local/go/bin/go install github.com/github/github-mcp-server/cmd/github-mcp-server@latest

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

# Final stage - copy Go binary to Ubuntu runtime
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

# Ensure Go is in PATH for the runtime user
ENV PATH="/usr/local/go/bin:$PATH"

# Switch to non-root user
USER iflow

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/iflow-action"]