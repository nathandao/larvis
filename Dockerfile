# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.20-alpine as build

WORKDIR /tmp/larvis

COPY *.go ./
COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go build -o build/larvis *.go

# Run stage
FROM alpine:3

COPY --from=build /tmp/larvis/build/larvis /usr/local/bin/larvis
