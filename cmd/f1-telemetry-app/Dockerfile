FROM golang:1.18-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /build/f1-2020-telemetry

FROM scratch

COPY --from=builder ["/build/f1-2020-telemetry", "/"]

EXPOSE 8080
EXPOSE 20777/udp

# Command to run when starting the container.
ENTRYPOINT ["/f1-2020-telemetry"]