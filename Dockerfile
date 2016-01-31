FROM golang:1.5-onbuild

MAINTAINER Jasper Lievisse Adriaanse <j@jasper.la>

RUN apt-get -qy update && apt-get -qy install wkhtmltopdf xvfb xauth --no-install-recommends && \
	apt-get clean && rm -fr /tmp/* /var/tmp/*

EXPOSE 9090
