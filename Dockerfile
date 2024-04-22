FROM golang:1.22.1-alpine
WORKDIR /filmoteka
COPY . .
ENV GOOSE_DBSTRING="postgres://docker:docker@postgres/docker?sslmode=disable"
ENV GOOSE_DRIVER="postgres"
ENV GOOSE_MIGRATION_DIR="db/migrations"
EXPOSE 8080
RUN chmod +x app.sh
CMD ["./app.sh"]

