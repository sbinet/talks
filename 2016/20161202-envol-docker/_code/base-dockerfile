## base environment w/ deps for the web-app
FROM golang
MAINTAINER binet@cern.ch

## dependencies already contained in `golang`
## but we need 'pkg-config'
RUN apt-get update -y && apt-get install -y pkg-config

## prepare for web server
EXPOSE 8080
