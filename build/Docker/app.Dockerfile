# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.17-alpine base image
FROM golang:1.17-alpine

# Add Maintainer Info
LABEL maintainer="Oleg Belov <oleg.s.belov@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/web/

# Expose port 4000 to the outside world
EXPOSE 4000

# Run the executable
CMD ["./main", "-db-url=postgres://my_user:myUser!2021@postgresdb:5432/snippetbox?sslmode=disable"]
