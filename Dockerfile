FROM golang:1.15-alpine as builder
RUN adduser -D -g '' appuser
RUN apk update && apk add --no-cache make git ca-certificates && update-ca-certificates
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get github.com/go-mail/mail
RUN go get github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -extldflags -s" -o ./app ./main.go

#FROM scratch
FROM ubuntu:latest
ENV GOOGLE_APPLICATION_CREDENTIALS=/auth.json
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/app/js /js
COPY --from=builder /go/src/app/css /css
COPY --from=builder /go/src/app/img /img
COPY --from=builder /go/src/app/fonts /fonts
COPY --from=builder /go/src/app/app /app
COPY --from=builder /go/src/app/index.html /index.html
COPY --from=builder /go/src/app/home.html /home.html
COPY --from=builder /go/src/app/confirmation.html /confirmation.html
USER appuser
ENTRYPOINT ["/app"]
