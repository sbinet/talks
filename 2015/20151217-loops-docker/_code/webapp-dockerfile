## image for the web-app
FROM binet/web-base

MAINTAINER binet@cern.ch

## add the whole git-repo
ADD . /go/src/github.com/sbinet/loops-20151217-tp

RUN go install github.com/sbinet/loops-20151217-tp/web-app

CMD web-app
