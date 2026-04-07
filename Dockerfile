# Build container
#FROM alpine:latest
FROM golang:1.26.1-alpine
RUN apk -U upgrade --no-cache && apk add --no-cache git && rm -rf /var/cache/apk/* && mkdir /build
WORKDIR /build
RUN git clone https://github.com/johannes-kuhfuss/jobsvc.git
WORKDIR /build/job-service
RUN go build -o /build/job-service/job-service /build/job-service/main.go
# Run container
FROM alpine:3.23.1
RUN apk -U upgrade --no-cache && rm -rf /var/cache/apk/* && mkdir /app
WORKDIR /app
COPY --from=0 /build/job-service/job-service /app/job-service
COPY --from=0 /build/job-service/templates /app/templates
RUN addgroup -g 101 servicegroup && adduser -s /sbin/nologin -G servicegroup -D -H -u 101 serviceuser
USER serviceuser
HEALTHCHECK --interval=120s --timeout=5s CMD wget -q --spider http://localhost:8080/ || exit 1
ENTRYPOINT ["/app/job-service"]
EXPOSE 8080/tcp