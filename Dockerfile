FROM golang:1.23-alpine AS golang-builder

WORKDIR /go/src/github.com/imtaiwanese18741130/bobcorn
COPY . .
RUN go mod tidy
RUN go build -o /go/bin/bobcorn .

FROM alpine:latest
COPY --from=golang-builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=golang-builder /go/bin/bobcorn /var/www/app/
COPY --from=golang-builder /go/src/github.com/imtaiwanese18741130/bobcorn/assets /var/www/app/assets
COPY --from=golang-builder /go/src/github.com/imtaiwanese18741130/bobcorn/templates /var/www/app/templates
WORKDIR /var/www/app

EXPOSE 8080
ENV APP_PORT 8080

ENTRYPOINT ["./bobcorn"]
