FROM golang:alpine
RUN apk add --no-cache make build-base
WORKDIR /go/src/app
COPY . .
RUN ./deploy_frontend.sh
RUN go get -d -v ./...
RUN go install -v ./... 
EXPOSE 8000 10000-10007/udp 
CMD ["RETHi-Comm"]

# docker pull golang:alpine
# docker buildx build -t amyangxyz111/rethi-comm --platform linux/arm64,linux/amd64 .
# docker run -d -p8000:8000 -p127.0.0.1:10000-10007:10000-10007/udp --name comm amyangxyz111/rethi-comm