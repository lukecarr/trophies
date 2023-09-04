FROM oven/bun:0.8 AS frontend-builder

WORKDIR /build
COPY ./frontend ./

RUN bun install --frozen-lockfile
RUN bun run build

FROM golang:1.21-alpine AS backend-builder

RUN apk update && apk add --no-cache git ca-certificates tzdata build-base && update-ca-certificates

ENV USER=trophies
ENV UID=10001

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

WORKDIR $GOPATH/src/github.com/lukecarr/trophies
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .
COPY --from=frontend-builder /build/dist ./frontend/dist/

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ARG BUILD_VERSION
ARG BUILD_COMMIT
ARG BUILD_DATE
RUN go build -ldflags="-w -s -extldflags '-static' -X github.com/lukecarr/trophies/internal/info.Version=${BUILD_VERSION} -X github.com/lukecarr/trophies/internal/info.Commit=${BUILD_COMMIT} -X github.com/lukecarr/trophies/internal/info.Date=${BUILD_DATE}" -a -o /usr/bin/trophies main.go

FROM scratch

COPY --from=backend-builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-builder /etc/passwd /etc/passwd
COPY --from=backend-builder /etc/group /etc/group

COPY --from=backend-builder /usr/bin/trophies /usr/bin/trophies

ENV DSN=/data/trophies.db

EXPOSE 3000

USER trophies

ENTRYPOINT ["trophies"]
CMD ["serve"]
