# Stage 1: Build Angular App
FROM node:22-alpine as angular_builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm install

COPY ./frontend .

RUN npm run build

# Stage 2: Build Golang App
FROM golang:1.23-alpine as go_builder

ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy Golang source code
COPY ./backend .

# Copy Angular build artifacts from the previous stage
COPY --from=angular_builder /app/dist/ ./backend/static

# Build the Golang app
RUN go build -o main .

# Stage 3: Run the Golang App
FROM alpine:latest

WORKDIR /app

# Copy the built Golang binary from the previous stage
COPY --from=go_builder /app/main .

# Expose the port your Golang app listens on (e.g., 8080)
EXPOSE 8080

# Command to run the Golang app
CMD ["./main"]
