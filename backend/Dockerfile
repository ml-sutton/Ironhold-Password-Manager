FROM golang:1.26.3-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ENV CGO_ENABLED=0 GOMAXPROCS=1
RUN go build -trimpath -ldflags="-s -w" -o main ./cmd/api/

FROM busybox:1.36 AS toolbox


FROM gcr.io/distroless/static-debian13:nonroot AS runtime
WORKDIR /api
COPY --from=builder /app/main .
EXPOSE 80
COPY --from=toolbox /bin/busybox /bin/busybox
USER nonroot:nonroot

HEALTHCHECK --interval=30s --timeout=5s --start-period=90s --retries=5 \
  CMD ["/bin/busybox", "sh", "-c", "pgrep main >/dev/null 2>&1 || exit 1"]

CMD ["./main"]