FROM mysql:5.7

COPY . /docker-entrypoint-initdb.d
RUN chmod -R 0755 /docker-entrypoint-initdb.d
