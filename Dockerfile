FROM golang:alpine AS builder

WORKDIR /build

COPY ./ ./

RUN go mod download
RUN go build -o surge .

WORKDIR /dist
RUN cp /build/surge .

FROM alpine AS runtime
COPY --from=builder /dist/surge .
COPY schema ./schema