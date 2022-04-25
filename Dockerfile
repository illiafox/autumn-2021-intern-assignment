# build stage
FROM golang:1.18-alpine AS build-env
RUN apk --no-cache add build-base git curl
ADD . /server
RUN cd /server/cmd/app && go build -o bin

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /server/cmd/app/bin /app/bin
COPY --from=build-env /server/cmd/app/config.toml /app/config.toml
ENTRYPOINT ./bin $ARGS