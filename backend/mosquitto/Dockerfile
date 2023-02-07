FROM eclipse-mosquitto:latest

COPY backend/mosquitto/docker-entrypoint.sh /

ENTRYPOINT ["sh", "docker-entrypoint.sh"]

CMD ["/usr/sbin/mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]