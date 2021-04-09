FROM golang:1.16.3-alpine3.13 as builder
COPY . /app
WORKDIR /app
RUN go build

FROM alpine:3.13
COPY --from=builder /app/sysctl_exporter /bin/sysctl_exporter
ENTRYPOINT [ "/bin/sysctl_exporter" ]