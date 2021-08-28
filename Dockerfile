FROM golang
WORKDIR /go/src/github.com/sl4mmer/crawlAndBench
COPY . ./
RUN go mod vendor
EXPOSE 8080
CMD ["go","run","cmd/main.go"]