FROM mongo:8.2.7




RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
RUN  curl -o /usr/local/bin/gomplate -sSL https://github.com/hairyhenderson/gomplate/releases/download/v5.0.0/gomplate_linux-amd64
RUN chmod 755 /usr/local/bin/gomplate


COPY /databases/mongo/init /docker-entrypoint-initdb.d/

COPY /databases/mongo/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/docker-entrypoint.sh


