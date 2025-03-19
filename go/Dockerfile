FROM golang:1.24.1-alpine3.21 as builder

RUN apk add git

LABEL maintainer="siransbach"

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o taskmanagementapi .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

RUN mkdir /app

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/taskmanagementapi .

# Expose port 8000
EXPOSE 8000

# Run Executable
CMD ["./taskmanagementapi"]
