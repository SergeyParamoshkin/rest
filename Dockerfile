FROM alpine:3.10 as app

RUN apk --no-cache upgrade && apk --no-cache add ca-certificates
ADD testing /usr/local/bin/testing
WORKDIR /usr/local/bin/

CMD ["testing"]