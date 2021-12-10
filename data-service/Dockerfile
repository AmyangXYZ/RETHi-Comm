FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go mod download &&\
    go build -o data_service

FROM alpine
COPY --from=builder /go/src/app/data_service ./
EXPOSE 8000 10000-10007/udp 
CMD ["./data_service"]

# build locally: docker build -t data-service . 
# build and push for arm64 and amd64: docker buildx build -t 'your_docker_hub_repo' --platform linux/arm64,linux/amd64 --push .