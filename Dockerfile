FROM golang:alpine3.18 AS builder
RUN apk add --no-cache --update git build-base

WORKDIR /app
COPY . .
RUN go build \
    -a \
    -trimpath \
    -o annona_core \
    -ldflags "-s -w -buildid=" \
    "./cmd/annona_core" && \
    ls -lah

FROM alpine:3.18 AS runner
RUN apk --no-cache add ca-certificates tzdata
#ENV LANG C.UTF-8
#ENV LANGUAGE en_US:en
#ENV LC_ALL C.UTF-8
ENV TZ UTC
WORKDIR /app

COPY --from=builder /app/annona_core .
VOLUME /app/conf
VOLUME /app/autocert
VOLUME /app/log
EXPOSE 8080

#ENTRYPOINT ["./annona_core" ,"-c","/app/conf/config.yaml"]
ENTRYPOINT ["./annona_core"]