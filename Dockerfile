# --- builder ---
FROM golang:1.22 as builder

ARG BUILD_VERSION
ARG BUILD_LAST_COMMIT

RUN mkdir /app
RUN mkdir /app/rueidis-test

WORKDIR /app/rueidis-test

COPY go.mod /app/rueidis-test/go.mod
COPY go.sum /app/rueidis-test/go.sum

RUN go env
RUN go env -w CGO_ENABLED=0
RUN go mod download

COPY . /app/rueidis-test

RUN go build \
    -a \
    -o "release/main" \
    "./main.go"

# --- runner ---
FROM debian as runner

RUN apt update && apt upgrade -y && apt install -y ca-certificates curl && update-ca-certificates

COPY --from=builder /app/rueidis-test/release/main /app/rueidis-test/release/main

RUN mkdir -p /usr/local/bin
RUN ln -s /app/rueidis-test/release/main /usr/local/bin/rueidis-test

CMD [ "/usr/local/bin/rueidis-test" ]
