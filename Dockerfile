FROM golang:1.15 AS api
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go build

FROM node:15 AS ux
WORKDIR /app
COPY ux .
RUN yarn install
RUN yarn build

FROM ubuntu:20.04
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive  apt-get install -y python3 python3-pip ffmpeg \
 && pip3 install youtube-dl \
 && apt-get remove -y python3-pip \
 && apt-get autoremove -y \
 && rm -rf /var/lib/apt/lists/*

COPY --from=api /app/mixbooth /usr/bin/mixbooth
COPY --from=ux  /app/dist     /htdocs
COPY            ingest        /usr/bin

EXPOSE 5000

ENV HTDOCS_ROOT=/htdocs
CMD ["mixbooth"]
