FROM golang:1.23-alpine

WORKDIR /app

# Install dogo and necessary build tools
RUN apk add --no-cache git \
    && go install github.com/liudng/dogo@latest

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy dogo config
COPY dogo.json* ./

EXPOSE 8080

CMD ["dogo"]