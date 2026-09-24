FROM docker.io/library/golang:1.27.0-alpine AS builder

ARG VERSION=dev
ARG COMMIT_HASH
ARG DATETIME

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X github.com/fr33dman/go-template/internal/build.Version=${VERSION} -X github.com/fr33dman/go-template/internal/build.CommitHash=${COMMIT_HASH} -X github.com/fr33dman/go-template/internal/build.Datetime=${DATETIME}" \
    -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X github.com/fr33dman/go-template/internal/build.Version=${VERSION} -X github.com/fr33dman/go-template/internal/build.CommitHash=${COMMIT_HASH} -X github.com/fr33dman/go-template/internal/build.Datetime=${DATETIME}" \
    -o /out/migrate ./cmd/migrate

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/server /app/server
COPY --from=builder /out/migrate /app/migrate

USER 65534

CMD ["/app/server"]
