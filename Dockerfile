FROM golang:1.13 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /workspace

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY main.go main.go
COPY kafka/ kafka/
COPY template/ template/

# Build the application
RUN go build -a -o compiler main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/compiler .
USER nonroot:nonroot

ENTRYPOINT ["/compiler"]