FROM redis:latest

## For security settings uncomment, make the dir, copy conf, and also start with the conf, to use it

#RUN mkdir -p /usr/local/etc/redis_client
#COPY redis_client.conf /usr/local/etc/redis_client/redis_client.conf

VOLUME /data

EXPOSE 6379

#CMD ["redis_client-server", "/usr/local/etc/redis_client/redis_client.conf"]
CMD ["redis-server"]
