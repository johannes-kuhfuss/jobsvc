# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

# Setup ENV
WORKDIR /app

# Download prereqs
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy sources
COPY . .

# Build
RUN go build -o /jobsvc

##
## Deploy
##
FROM alpine:latest

WORKDIR /

COPY --from=build /jobsvc /jobsvc

RUN chmod +x /jobsvc

EXPOSE 8080

USER nobody:nogroup

ENTRYPOINT ["/jobsvc"]
