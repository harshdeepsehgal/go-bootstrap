FROM golang:1.26-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /src
COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -o /api ./cmd/server

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /api /api
EXPOSE 8083
ENTRYPOINT ["/api"]
