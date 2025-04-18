FROM golang:1.24.2-alpine AS todo_builder

WORKDIR /app

COPY  . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./todo_app ./cmd/main.go

FROM alpine:3.21

# Создаю директорию для db так как  Go создаст scheduler.db динамически, а дирректория нужна
RUN mkdir -p /pkg/db

COPY --from=todo_builder ./app/todo_app ./todo_app
COPY --from=todo_builder ./app/web ./web
COPY --from=todo_builder ./app/.env ./.env

EXPOSE 7540

CMD [ "/todo_app" ]