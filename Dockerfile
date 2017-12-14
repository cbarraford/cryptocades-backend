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

RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    apk update && \
    apk add curl make git "libpq@edge<9.7" "postgresql-client@edge<9.7" "postgresql@edge<9.7" "postgresql-contrib@edge<9.7" && \
    apk del curl && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["dumb-init"]
CMD ["/bin/sh"]
