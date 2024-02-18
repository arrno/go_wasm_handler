FROM golang:1.22

WORKDIR /handler

COPY handler/go.mod ./go.mod
COPY handler/go.sum ./go.sum
COPY handler/gcp.go ./gcp.go
COPY handler/main.go ./main.go
COPY handler/util.go ./util.go
COPY handler/keys ./keys
COPY wasm ./wasm

# Build
RUN go build -o main .
EXPOSE 8080

# Run
ENV GOPATH=/handler
CMD ["/handler/main", ":8080"]