FROM alpine:3.19

WORKDIR /app
COPY ./build .
EXPOSE 8080

# community app config env for config.yaml
ENV DB_HOST 127.0.0.1:3306
ENV DB_DATABASE community
ENV DB_USER root
ENV DB_PASS 123456
ENV KODO_ACCESS_KEY ""
ENV KODO_BUCKET ""
ENV KODO_SECRET_KEY ""

ENTRYPOINT ./community
