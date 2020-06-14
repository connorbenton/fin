FROM node:lts-alpine as build-stage

WORKDIR /app
COPY ./frontend/package*.json ./
RUN apk --no-cache --virtual .gyp add \
        python \
        make \
        g++ \
    && npm install \
    && apk del .gyp 
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

# production stage
FROM alpine:3.12.0 as production-stage

WORKDIR /usr/share/nginx/html
COPY --from=build-stage /app/dist ./
COPY ./entrypoints/env_load.sh ./

COPY --from=go-build /server /server
COPY ./backend_go /usr/src/app/backend_go

COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
COPY ./s6 /
 
RUN apk add --no-cache nginx && apk add --no-cache bash && chmod +x env_load.sh

ENTRYPOINT ["/init"]