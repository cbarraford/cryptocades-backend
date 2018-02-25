FROM golang:alpine

RUN \
	apk add --update \
		python \
		python-dev \
		py-pip \
		build-base \
	&& \
	pip install dumb-init && \
	apk del \
		python \
		python-dev \
		py-pip \
		build-base \
	&& \
	rm -rf /var/cache/apk/* && \
	:

RUN apk update && \
    apk add curl make git postgresql-client && \
    apk del curl && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["dumb-init"]
CMD ["/bin/sh"]
