# base image
# FROM node:10.15.0
# FROM node:alpine
FROM node:lts-alpine as build-stage

WORKDIR /app
COPY ./frontend/package*.json ./
RUN apk --no-cache --virtual .gyp add \
        python \
        make \
        g++ \
    && npm install \
    && apk del .gyp 
# RUN npm install
COPY ./frontend .
RUN npm run build-prod

FROM golang:1.14.4-alpine3.12 as go-build
ENV CGO_ENABLED 1
RUN apk add --update gcc musl-dev
RUN mkdir -p /usr/src/app/backend_go
RUN mkdir /usr/src/app/db
WORKDIR /usr/src/app/backend_go
COPY ./backend_go ./
RUN go mod vendor
RUN go build -gcflags "all=-N -l" -o /server main.go 

# #sqlite stage
# FROM coleifer/sqlite as build-stage
# RUN apk add --no-cache --virtual .build-reqs build-base gcc make \
#       && pip install --no-cache-dir cython \
#       && pip install --no-cache-dir flask peewee sqlite-web \
#       && apk del .build-reqs
# # EXPOSE 8080
# VOLUME /data
# WORKDIR /data

# production stage
# FROM nginx:stable-alpine as production-stage
FROM alpine:3.12.0 as production-stage

# RUN apk add --no-cache python py-pip 
# RUN apk add --no-cache --virtual build-dependencies build-base gcc make \
#       && pip install --no-cache-dir cython \
#       && pip install --no-cache-dir flask peewee sqlite-web \
#       && apk del build-dependencies

# RUN apk add --no-cache supervisor && apk add --no-cache bash
# RUN apk add --update supervisor && rm  -rf /tmp/* /var/cache/apk/*


# RUN apk add --update npm
# WORKDIR /usr/src/app/backend
# COPY ./backend/package*.json ./
# RUN npm ci --only=production
# COPY ./backend .
# WORKDIR /usr/src/app/backend_go
WORKDIR /usr/share/nginx/html
COPY --from=build-stage /app/dist ./
COPY ./entrypoints/env_load.sh .env ./

COPY --from=go-build /server /server
COPY ./backend_go /usr/src/app/backend_go

COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
# COPY ./nginx/gzip.conf /etc/nginx/gzip.conf
# COPY .env .
# RUN apk add --no-cache bash

# RUN apk add --update --no-cache supervisor && apk add --no-cache bash && chmod +x env_load.sh
# RUN apk add --update --no-cache supervisor && apk add --no-cache nginx && apk add --no-cache bash && chmod +x env_load.sh
# RUN apk add --no-cache --virtual .pipbuild python3 py3-pip && pip install --no-cache-dir supervisor && apk del .pipbuild \ 
RUN apk add --no-cache py3-pip && pip install --no-cache-dir supervisor \ 
&& apk add --no-cache nginx && apk add --no-cache bash && chmod +x env_load.sh


# RUN apk add --update supervisor && apk add --update bash && rm -rf /tmp/* /var/cache/apk/* && chmod +x env_load.sh
# RUN chmod +x env_load.sh

COPY ./entrypoints/supervisord.conf /etc/supervisord.conf
ENTRYPOINT ["supervisord", "--nodaemon", "--configuration", "/etc/supervisord.conf"]
# CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]

# EXPOSE 80
# CMD ["nginx", "-g", "daemon off;"]


# FROM coleifer/sqlite
# RUN apk add --no-cache --virtual .build-reqs build-base gcc make \
#       && pip install --no-cache-dir cython \
#       && pip install --no-cache-dir flask peewee sqlite-web \
#       && apk del .build-reqs
# # EXPOSE 8080
# VOLUME /data
# WORKDIR /data
# CMD sqlite_web -H 0.0.0.0 -x $SQLITE_DATABASE

 
# # set working directory
# RUN mkdir /usr/src/app
# WORKDIR /usr/src/app
 
# # add `/usr/src/app/node_modules/.bin` to $PATH
# ENV PATH /usr/src/app/node_modules/.bin:$PATH

# # RUN ls /usr/local/bin
# # install and cache app dependencies
# COPY package.json /usr/src/app/package.json
# RUN npm install
# RUN npm install -g @vue/cli @vue/cli-service vue-template-compiler vuetify
# # RUN ls /usr/src/app

# # FROM coleifer/sqlite

# # VOLUME /data
# # RUN mkdir /usr/src/app
# # VOLUME /usr/src/app
# # WORKDIR /usr/src/app

# # COPY --from=0 /usr/ /usr/
# # ENV PATH /usr/src/app/.local/lib

# # RUN apk add --no-cache --virtual .build-reqs build-base gcc make \
# RUN apt-get update && apt-get install -y gcc make python3 python3-pip\
#     #   && pip install --no-cache-dir cython \
#     #   && pip install --no-cache-dir flask peewee sqlite-web \
#     && pip3 install --no-cache-dir sqlite-web
# # EXPOSE 8080

# ENV PATH /usr/local/lib/python3.7/site-packages:$PATH
# # ENV PATH /usr/src/app/node_modules/.bin:$PATH
# # RUN ls /usr/src/app/node_modules/.bin
# # RUN ls /usr/local/lib/python3.7/site-packages
# # RUN ls /usr/src/app/
# # RUN ls /usr/local/bin
# # CMD sqlite_web -H 0.0.0.0 -x data.sqlite

# # RUN pip show sqlite-web
# # COPY sqlite_web sqlite_web

# # Add Tini
# # ENV TINI_VERSION v0.18.0
# # ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
# # RUN chmod +x /tini


# ADD ./entrypoints/start-stage.sh /
# RUN chmod +x /start-stage.sh
# CMD ["/bin/sh", "/start-stage.sh"]


# ENTRYPOINT ["/tini", "--", "/start.sh"]
# CMD ["/start.sh"]

# start app
# CMD ["npm", "run", "dev"]