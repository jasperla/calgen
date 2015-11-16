FROM golang:1.5-onbuild

MAINTAINER Jasper Lievisse Adriaanse <j@jasper.la>

RUN apt-get -qy update && apt-get -qy install wkhtmltopdf

EXPOSE 9090
