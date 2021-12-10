FROM golang:alpine AS builder
RUN apk add --no-cache make build-base
WORKDIR /go/src/app
COPY . .
RUN go mod download &&\
    go build -o RETHi-Comm . 

FROM alpine
COPY --from=builder /go/src/app/RETHi-Comm ./
COPY dist.tar.xz comm.db flex_config.json ./
RUN mkdir templates &&\
    tar -xvf dist.tar.xz &&\
    mv ./dist/static ./static &&\
    mv ./dist/index.html ./templates/index.html &&\
    rm -rf dist dist.tar.xz &&\
    sed -i 's/^/{{define "index"}}/' ./templates/index.html &&\
    sed -i 's/$/{{end}}/' ./templates/index.html 
EXPOSE 8000 10000-10007/udp 
CMD ["./RETHi-Comm"]

# docker pull golang:alpine
# docker buildx build -t amyangxyz111/rethi-comm --platform linux/arm64,linux/amd64 .
# docker run -d -p8000:8000 -p"127.0.0.1:10000-10007:10000-10007/udp" --name comm amyangxyz111/rethi-comm