FROM golang:1.19.3-buster
ENV LANG C.UTF-8
ENV PATH=${PATH}:${GOPATH}/bin
RUN apt update && apt install -y \
    build-essential \
    ca-certificates \
    git \
    protobuf-compiler \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*
RUN go install github.com/golang/protobuf/protoc-gen-go@latest
RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
ADD https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-linux-x86_64 /usr/local/bin
RUN mv /usr/local/bin/protoc-gen-grpc-web-1.4.2-linux-x86_64 /usr/local/bin/protoc-gen-grpc-web
RUN chmod +x /usr/local/bin/protoc-gen-grpc-web