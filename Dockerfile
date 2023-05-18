FROM golang:1.14.6-alpine3.12 as builder

COPY go.mod go.sum /go/src/github.com/elman23/articleapi/
WORKDIR /go/src/github.com/elman23/articleapi
RUN go mod download
COPY . /go/src/github.com/elman23/articleapi
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/articleapi github.com/elman23/articleapi
RUN GOOS=linux go build -o build/articleapi /go/src/github.com/elman23/articleapi/main.go

FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/elman23/articleapi/build/articleapi /usr/bin/articleapi

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/articleapi"]
