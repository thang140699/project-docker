# syntax=docker/dockerfile:1
FROM docker/getting-started:latest

WORKDIR /go/src/github.com/alexellis/href-counter/
COPY main.go ./
RUN go get -d -v golang.org/x/net/html \
  && CGO_ENABLED=0 go build -a -installsuffix cgo -o main .


FROM authen

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/alexellis/href-counter/app ./
CMD ["./controller"]