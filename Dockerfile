FROM golang:1.12.7-alpine AS build

COPY . /go/src/github.com/abousselmi/go-jenkins-exporter
WORKDIR /go/src/github.com/abousselmi/go-jenkins-exporter
RUN apk update && apk -U add git \
	&& GO111MODULE=on go get -v \
	&& CGO_ENABLED=0 go build -a installsuffix nocgo -o /go/bin/go-jenkins-exporter .

FROM scratch

LABEL description="A simple jenkins exporter for prometheus, written in Go."
COPY --from=build /go/bin/go-jenkins-exporter /app/go-jenkins-exporter
WORKDIR /app
EXPOSE 5000
CMD ["./go-jenkins-exporter"]
