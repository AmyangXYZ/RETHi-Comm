FROM golang:alpine AS backend-builder
RUN apk add --no-cache make build-base git
WORKDIR /go/src/app
COPY ./* .
RUN go mod download 
RUN go install ./... .

FROM node:lts-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
ENV NODE_OPTIONS="--openssl-legacy-provider"
RUN npm install
COPY frontend/ .
RUN npm run build

FROM alpine
WORKDIR /app
COPY --from=frontend-builder /app/dist ./
COPY --from=backend-builder /go/bin/rethi-comm ./
COPY topos.json ./
RUN mkdir templates &&\
    mv ./index.html ./templates/index.html &&\
    sed -i 's/^/{{define "index"}}/' ./templates/index.html &&\
    sed -i 's/$/{{end}}/' ./templates/index.html 
EXPOSE 8000 10000-10007/udp 
CMD ["./rethi-comm"]

# docker buildx build -t amyangxyz111/rethi-comm --platform linux/arm64,linux/amd64 .