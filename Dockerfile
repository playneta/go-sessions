FROM golang:1.13-alpine as builder

WORKDIR /app
RUN apk add git
ADD . .

RUN set -x \
    && export VERSION=$(git rev-parse --verify HEAD --short) \
    && export LDFLAGS="-w -s -X main.Version=${VERSION}" \
    && export CGO_ENABLED=0 \
    && go build -v -ldflags "${LDFLAGS}" -o /chat .

FROM alpine:3.9

RUN apk add  ca-certificates
COPY --from=builder /chat /chat

CMD ["./chat serve"]
