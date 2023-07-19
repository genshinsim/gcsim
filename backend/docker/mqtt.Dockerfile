FROM eclipse-mosquitto:latest

COPY mosquitto/docker-entrypoint.sh /
RUN chmod +x docker-entrypoint.sh 

ENTRYPOINT ["sh", "docker-entrypoint.sh"]

CMD ["/usr/sbin/mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]