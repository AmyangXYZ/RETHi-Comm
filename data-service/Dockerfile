FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go install ./... .

FROM alpine
COPY --from=builder /go/bin/data-service ./
EXPOSE 8000 10000-10007/udp 
CMD ["./data-service"]

# build locally: docker build -t data-service . 
# build and push for arm64 and amd64: docker buildx build -t 'your_docker_hub_repo' --platform linux/arm64,linux/amd64 --push .