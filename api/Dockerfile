# Stage 1: Build the Go application
FROM golang:alpine as builder

# Create a directory for building the application
RUN mkdir /build

# Copy the source code into the build directory
ADD . /build/

# Set the working directory to the build directory
WORKDIR /build

# Build the Go application and generate the 'main' binary
RUN go build -o main .

# Stage 2: Create a minimal runtime environment
FROM alpine

# Create a non-root user 'appuser' for running the application
RUN adduser -S -D -H -h /app appuser

# Switch to the 'appuser' for improved security
USER appuser

# Copy the source code from the host into the container's '/app' directory
COPY . /app

# Copy the 'main' binary from the builder stage to the '/app' directory
COPY --from=builder /build/main /app/

# Set the working directory to '/app'
WORKDIR /app

# Expose port 3000 (metadata hint)
EXPOSE 3000

# Define the command to run the application
CMD ["./main"]
