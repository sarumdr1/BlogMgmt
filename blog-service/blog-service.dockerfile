FROM alpine:latest

RUN mkdir /app

COPY blogServiceApp /app

CMD [ "/app/blogServiceApp"]