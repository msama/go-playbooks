# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

MAINTAINER Michele Sama

WORKDIR /go/src/aws-ebs-nginx
ADD . /go/src/aws-ebs-nginx
RUN go build -o aws-ebs-nginx main.go

ENTRYPOINT ./aws-ebs-nginx
EXPOSE 3000