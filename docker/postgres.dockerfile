FROM postgres:18.2


RUN apt update && apt install -y curl 

RUN  curl -o /usr/local/bin/gomplate -sSL https://github.com/hairyhenderson/gomplate/releases/download/v5.0.0/gomplate_linux-amd64
RUN chmod 755 /usr/local/bin/gomplate




COPY /databases/postgres/templates /docker-entrypoint-initdb.d/

COPY /databases/postgres/scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN sed -i 's/\r$//'  /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
