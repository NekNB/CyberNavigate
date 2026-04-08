FROM postgres:18.2


RUN apt update && apt install -y gettext


COPY postgres/templates /docker-entrypoint-initdb.d/

COPY postgres/scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

# ENTRYPOINT [ "" ]
# CMD sleep infinity