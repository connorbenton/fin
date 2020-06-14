# base image
FROM node:10.15.0
 
RUN mkdir /usr/src/app
# install and cache app dependencies
RUN mkdir /usr/src/app/frontend
WORKDIR /usr/src/app/frontend
COPY ./frontend/package*.json ./
RUN npm install
RUN npm install -g @vue/cli @vue/cli-service vue-template-compiler vuetify

WORKDIR /usr/src/app
RUN apt-get update && apt-get install -y gcc make python3 python3-pip\
    && pip3 install --no-cache-dir sqlite-web

ENV PATH /usr/local/lib/python3.7/site-packages:$PATH

COPY ./entrypoints/start-dev.sh /
RUN chmod +x /start-dev.sh
CMD ["/bin/bash", "-c", "source /start-dev.sh"]