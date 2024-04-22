# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app
COPY go.mod config.conf ./
RUN go mod download
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /ctc
CMD [ "/ctc" ]
