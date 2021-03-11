FROM golang:1.15-alpine as builder
RUN adduser -D -g '' appuser
RUN apk update && apk add --no-cache make git ca-certificates && update-ca-certificates
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get github.com/go-mail/mail
RUN go get github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -extldflags -s" -o ./app ./main.go

FROM scratch
ENV GOOGLE_APPLICATION_CREDENTIALS=/auth.json
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/app/app /
USER appuser
ENTRYPOINT ["/app"]
