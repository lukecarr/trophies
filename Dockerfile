FROM node:20 AS node-builder

RUN npm i -g pnpm@8.7.0

WORKDIR /build
COPY ./frontend ./

RUN pnpm install --frozen-lockfile
RUN pnpm build

FROM golang:1.21-alpine AS go-builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

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
COPY --from=node-builder /build/dist ./frontend/dist/

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ARG BUILD_VERSION
ARG BUILD_COMMIT
ARG BUILD_DATE
RUN go build -ldflags="-w -s -extldflags '-static' -X github.com/lukecarr/trophies/internal/info.Version=${BUILD_VERSION} -X github.com/lukecarr/trophies/internal/info.Commit=${BUILD_COMMIT} -X github.com/lukecarr/trophies/internal/info.Date=${BUILD_DATE}" -a -o /usr/bin/trophies main.go

FROM scratch

COPY --from=go-builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /etc/passwd /etc/passwd
COPY --from=go-builder /etc/group /etc/group

COPY --from=go-builder /usr/bin/trophies /usr/bin/trophies

EXPOSE 3000

USER trophies

ENTRYPOINT ["trophies"]
CMD ["serve"]
