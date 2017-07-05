FROM golang:alpine as build
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build

FROM golang:alpine
COPY --from=build /go/src/app/authobot /bin/authobot
ENTRYPOINT ["/bin/authobot"]