# syntax=docker/dockerfile:1

FROM cgr.dev/chainguard/go AS build

WORKDIR /usr/src

COPY go.* .

RUN go mod download

COPY . .

RUN go build -tags=timetzdata -o=/opt ./cmd/...

RUN mkdir -p /var/lib/sqlite

FROM cgr.dev/chainguard/glibc-dynamic

COPY --from=build /opt /usr/bin/
COPY --from=build /var/lib/sqlite /var/lib/sqlite

VOLUME [ "/var/lib/sqlite" ]

EXPOSE 8080

ENTRYPOINT [ "ph" ]

CMD ["serve"]