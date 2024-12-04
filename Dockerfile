FROM golang:1.23-alpine

WORKDIR /app

# Install dogo and necessary build tools
RUN apk add --no-cache git \
    && go install github.com/liudng/dogo@latest \
    && apk del git

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

RUN go mod tidy

# Copy source code
COPY . .

# Copy dogo config
COPY dogo.json* ./

EXPOSE 8080

CMD ["dogo"]