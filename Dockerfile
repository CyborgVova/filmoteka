FROM golang:alpine
WORKDIR /filmoteka
COPY . .
ENV GOOSE_DBSTRING="host=postgres dbname=docker user=docker password=docker sslmode=disable"
ENV GOOSE_DRIVER="postgres"
EXPOSE 8080
RUN chmod +x app.sh
CMD ["sh", "./app.sh"]
