FROM alpine:3.8

ENV GOSLEEP_VERSION 1.0

RUN set -ex; \
	\
	apk add --no-cache --virtual .fetch-deps \
		dpkg \
		gnupg \
	; \
	\
	dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')"; \
	wget -O /usr/local/bin/gosleep "https://github.com/tianon/gosleep/releases/download/$GOSLEEP_VERSION/gosleep-$dpkgArch"; \
	wget -O /usr/local/bin/gosleep.asc "https://github.com/tianon/gosleep/releases/download/$GOSLEEP_VERSION/gosleep-$dpkgArch.asc"; \
	\
	export GNUPGHOME="$(mktemp -d)"; \
# gpg: key BF357DD4: public key "Tianon Gravi <tianon@tianon.xyz>" imported
	gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4; \
	gpg --batch --verify /usr/local/bin/gosleep.asc /usr/local/bin/gosleep; \
	rm -r "$GNUPGHOME" /usr/local/bin/gosleep.asc; \
	\
	chmod +x /usr/local/bin/gosleep; \
	\
	apk del .fetch-deps; \
	\
	gosleep --help; \
	time gosleep --for 1s

CMD ["gosleep"]
