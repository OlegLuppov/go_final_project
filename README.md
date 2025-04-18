# Итоговый проект Scheduler Tasks

### Описание
Планировщик задач позволяет выполнять следующие действия:
- Создавать задачу
- Редактировать задачу
- Помечать как выполненную
- Удалять задачу

### Список выполненных заданий со звёздочкой:
- Реализована возможность определять порт из переменной окружения `TODO_PORT` при запуске сервера.
- Реализована возможность определять путь к базе данных из переменной окружения `TODO_DBFILE`.
- Реализован расширенный функционал определения `следующей даты` для переноса задачи по неделям и месяцам.
- Реализована возможность поиска задач через `строку поиска`.
- Реализована `Аутентификация` по токену.
- Написан `Dockerfile` для создания образа и дальнейшего запуска контейнера.

### Проверка работы приложения локально
1. Сделайте `fork` репозитория
2. Склонируте удаленный репозиторий к себе на локальную машину `git clone yours_link`
3. Запустите сервис командой в терминале `go run ./cmd/main.go`
4. Если есть переменная окружения `TODO_DBFILE`, тогда путь до базы данных будет браться из ее значения, если нет тогда будет использован дефолтный путь. Так же если файл с базой данных отсутствует, тогда он создасться динамически по указанному пути.
5. Запустите команду `go test ./tests` для тестирования сервиса, если тесты прошли успешно, можно переходить к следующим шагам.
6. Если у вас локально нет переменной окружения `TODO_PORT`, тогда приложение будет доступно под дефолтным портом 7540, нужно будет перейти по ссылке в браузере: `http://localhost:7540`, если такая переменная есть тогда порт подтянется из нее и нужно будет подставить порт `http://localhost:yours_port`
7. Если есть переменная окружения `TODO_PASSWORD`, тогда при первом входе на страницу откроется форма авторизации при помощи пароля, далее при успешной аутентификации будет доступен список задач, пока токен установлен в Cookie, если почистить Cookie или закрыть браузер, авторизацию нужно будет проходить заного. Если переменной нет то список задач будет доступен без авторизации.
8. При успешной авторизации попробуйте `Добавить`, `Отредактировать`, `Пометить выполненной`, `Удалить` задачу в браузере.

### Проверка сборки образа и запуска контейнера
1. Для сборки образа выполните комманду в терминале `docker build -t yours_nikname/todo-service:v1.0.0 .`
2. Для запуска контейнера выполните следующую комманду:
    - Если используется дефолтный порт: `docker run -p 7540:7540 image_name`
    - Если задана переменная окружения `TODO_PORT`: `docker run -p yours_port:yours_port image_name`
3. Перейдите по ссылке `http://localhost:yours_port` если задан `TODO_PORT` или `http://localhost:7540` если не задан.
4. Попробуйте `Добавить`, `Отредактировать`, `Пометить выполненной`, `Удалить` задачу в браузере.

### STACK

Back-end:

- [Go](https://go.dev/)
- [Sqlite](https://sqlite.org/)

Front-end:
- [HTML](https://developer.mozilla.org/ru/docs/Learn_web_development/Getting_started/Your_first_website/Creating_the_content)
- [CSS](https://developer.mozilla.org/ru/docs/Web/CSS)
- [JavaScript](https://developer.mozilla.org/ru/docs/Web/JavaScript)