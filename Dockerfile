FROM alpine:3.10 as app

RUN apk --no-cache upgrade && apk --no-cache add ca-certificates
ADD rest /usr/local/bin/rest
WORKDIR /usr/local/bin/

CMD ["rest"]