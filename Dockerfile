FROM golang as builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY . /go/src/backend
WORKDIR /go/src/backend
RUN go build -ldflags "-s -w"  -o backend cmd/main.go

FROM alpine
COPY --from=builder /go/src/backend/backend /backend

RUN apk --no-cache add curl

EXPOSE 8080/tcp
ENTRYPOINT ["/backend"]
