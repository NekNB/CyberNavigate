FROM mongo:8.2.7

RUN apt update && apt install -y curl 

RUN  curl -o /usr/local/bin/gomplate -sSL https://github.com/hairyhenderson/gomplate/releases/download/v5.0.0/gomplate_linux-amd64
RUN chmod 755 /usr/local/bin/gomplate


COPY /databases/mongo/init /docker-entrypoint-initdb.d/

COPY /databases/mongo/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh