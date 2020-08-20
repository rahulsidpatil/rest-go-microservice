# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

#adding this line so that we don't miss importing mysql driver 
RUN go get "github.com/go-sql-driver/mysql"

# Build the Go app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o rest-go-microservice cmd/main.go

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/rest-go-microservice .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./rest-go-microservice"]