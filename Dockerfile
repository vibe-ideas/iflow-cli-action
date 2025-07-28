# Use Node.js 22 as base for runtime with latest Node.js pre-installed
FROM node:22-slim AS runtime-base

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y \
        bash \
        curl \
        git \
        ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Pre-install iFlow CLI directly from the tgz package
RUN npm install -g https://cloud.iflow.cn/iflow-cli/iflow-iflow-cli-0.0.2.tgz

# Use official Go 1.24.4 image for building
FROM golang:1.24.4-bullseye AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y git ca-certificates curl

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iflow-action main.go

# Final stage - copy Go binary to Node.js runtime
FROM runtime-base

# Create a non-root user
RUN groupadd -g 1001 iflow && \
    useradd -r -u 1001 -g iflow iflow

# Set working directory
WORKDIR /github/workspace

# Copy the binary from builder stage
COPY --from=builder /app/iflow-action /usr/local/bin/iflow-action

# Make sure binary is executable
RUN chmod +x /usr/local/bin/iflow-action

# Switch to non-root user
USER iflow

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/iflow-action"]