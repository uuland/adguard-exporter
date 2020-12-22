FROM --platform=$BUILDPLATFORM golang:1.14-alpine AS build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /tmp/adguard_exporter

RUN apk update && apk --no-cache add git alpine-sdk upx
COPY . .
RUN GO111MODULE=on go mod vendor
RUN CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags '-s -w' -o adguard_exporter ./
RUN upx -f --brute adguard_exporter

FROM scratch
LABEL name="adguard-exporter"

WORKDIR /root
COPY --from=build /tmp/adguard_exporter/adguard_exporter adguard_exporter

CMD ["./adguard_exporter"]