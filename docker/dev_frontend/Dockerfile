# base image
FROM node:10.15.0
# FROM node:alpine
 
RUN mkdir /usr/src/app
# RUN mkdir /usr/src/app/backend
# WORKDIR /usr/src/app/backend
# COPY ./backend/package*.json ./
# RUN npm install

# # RUN mkdir /usr/local
# WORKDIR /usr/local
# RUN wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
# RUN tar -xvf go1.14.2.linux-amd64.tar.gz 
# ENV GOROOT /usr/local/go
# ENV PATH="/usr/local/go/bin:${PATH}"
# RUN mkdir /usr/local/gosrc
# ENV GOPATH /usr/local/gosrc
# RUN go get -u github.com/derekparker/delve/cmd/dlv
# ENV PATH="/usr/local/gosrc/bin:${PATH}"
# # RUN ls /usr/local/go/bin
# # RUN ls /usr/local/gosrc/bin

# # RUN go get github.com/codegangsta/gin

# RUN apt-get -y update && apt-get -y install inotify-tools
# RUN mkdir /usr/src/app/backend_go
# WORKDIR /usr/src/app/backend_go
# COPY ./backend_go ./
# # RUN go build -o /server main.go 
# RUN go build -gcflags "all=-N -l" -o /server main.go 
# # RUN go build -gcflags="-N -l" -o /server main.go 
# RUN chmod +x /usr/src/app/backend_go/start_script.sh

# set working directory
# RUN mkdir /usr/src/app
 
# add `/usr/src/app/node_modules/.bin` to $PATH
# ENV PATH /usr/src/app/node_modules/.bin:$PATH

# RUN ls /usr/local/bin
# install and cache app dependencies
RUN mkdir /usr/src/app/frontend
WORKDIR /usr/src/app/frontend
COPY ./frontend/package*.json ./
RUN npm install
RUN npm install -g @vue/cli @vue/cli-service vue-template-compiler vuetify
# RUN ls /usr/src/app


# FROM coleifer/sqlite

# VOLUME /data
# RUN mkdir /usr/src/app
# VOLUME /usr/src/app
# WORKDIR /usr/src/app

# COPY --from=0 /usr/ /usr/
# ENV PATH /usr/src/app/.local/lib

# RUN apk add --no-cache --virtual .build-reqs build-base gcc make \
WORKDIR /usr/src/app
RUN apt-get update && apt-get install -y gcc make python3 python3-pip\
    #   && pip install --no-cache-dir cython \
    #   && pip install --no-cache-dir flask peewee sqlite-web \
    && pip3 install --no-cache-dir sqlite-web
# EXPOSE 8080

ENV PATH /usr/local/lib/python3.7/site-packages:$PATH
# ENV PATH /usr/src/app/node_modules/.bin:$PATH
# RUN ls /usr/src/app/node_modules/.bin
# RUN ls /usr/local/lib/python3.7/site-packages
# RUN ls /usr/src/app/
# RUN ls /usr/local/bin
# CMD sqlite_web -H 0.0.0.0 -x data.sqlite

# RUN pip show sqlite-web
# COPY sqlite_web sqlite_web

# Add Tini
# ENV TINI_VERSION v0.18.0
# ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
# RUN chmod +x /tini

COPY ./entrypoints/start-dev.sh /
RUN chmod +x /start-dev.sh
# RUN ls /
# RUN ls ./
# RUN ls /usr/src/app
# RUN ls /usr/src/app/server
CMD ["/bin/bash", "-c", "source /start-dev.sh"]
# ENTRYPOINT ["/tini", "--", "/start.sh"]
# CMD ["/start.sh"]

# start app
# CMD ["npm", "run", "dev"]