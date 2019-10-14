# Start from golang base image
FROM golang:1.13-alpine as builder

LABEL maintainer="Alexandros Papapostolou <alpapapo@gmail.com>"

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
RUN mkdir /app
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Compile the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
RUN mkdir sample
COPY --from=builder /app/sample/data.json sample/

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]