FROM --platform=linux/x86_64 mysql:8.4-oraclelinux9

RUN --mount=type=cache,target=/root/.cache/microdnf:rw\
  microdnf install -y wget unzip openssl
RUN wget https://downloads.mysql.com/docs/world-db.zip \
  && unzip world-db.zip
RUN mv world-db/world.sql /docker-entrypoint-initdb.d/world.sql
