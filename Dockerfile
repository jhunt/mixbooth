FROM node:15 AS ux
WORKDIR /app
COPY ux /app
RUN yarn install
RUN yarn build

FROM perl:5.30
RUN apt-get update \
 && apt-get install -y python3-pip \
 && pip3 install youtube-dl

RUN cpanm Carton Starman

WORKDIR /app
COPY cpanfile .
RUN carton install

COPY . .
RUN carton exec -- perl -I/app/lib -c /app/bin/app.psgi

COPY --from=ux /app/dist /app/dist
EXPOSE 5000
ENV PATH=/usr/bin:/bin:/usr/sbin:/sbin:/app/bin:/usr/local/bin
CMD ["carton", "exec", "starman", "--port", "5000", "/app/bin/app.psgi"]

RUN apt-get install -y ffmpeg
