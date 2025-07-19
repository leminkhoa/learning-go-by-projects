# Use Go 1.24 with Alpine
FROM golang:1.24.5-alpine3.22 AS build

# Install necessary dependency
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory inside the container
WORKDIR /go/src/github.com/leminkhoa/go-grpc-graphql-microservice

# Copy Go module files and dependencies
COPY go.mod go.sum ./

# Copy project source files
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order


# Final stage: lightweight image for running the app
FROM alpine:3.22

# Set working directory for the app
WORKDIR /usr/bin

# Copy build from `build` to current 
COPY --from=build /go/bin .

# Expose app to port 8080
EXPOSE 8080
CMD ["app"]