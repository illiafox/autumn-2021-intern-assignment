# build stage
FROM golang:1.18-alpine AS build-env
RUN apk --no-cache add build-base git
ADD . /server
RUN cd /server/cmd/app && go build -o bin

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /server/cmd/app/ /app/
ENTRYPOINT ./bin