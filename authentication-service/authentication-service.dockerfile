FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY app.env /app

COPY ./db/migrations  /app/db/migrations

COPY authApp /app

CMD [ "/app/authApp"]
