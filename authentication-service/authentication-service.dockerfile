FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY app.env /app

COPY authApp /app

CMD [ "./authApp"]
