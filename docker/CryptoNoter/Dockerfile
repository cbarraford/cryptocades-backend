FROM yellowblackandwhite/alpine-npm

RUN set -x \
    && apk update \
    && apk add \
        git

# git pull CryptoNoter source code
RUN cd /srv \
    && git clone https://github.com/cbarraford/CryptoNoter.git -o CryptoNoter

WORKDIR /srv/CryptoNoter

RUN git checkout redis

RUN npm install

CMD ["npm", "start"]
