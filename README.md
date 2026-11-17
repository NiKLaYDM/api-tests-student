# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
TEST_FAVORITES_LOGIN=your_second_login
TEST_FAVORITES_PASSWORD=your_second_password
```

`TEST_FAVORITES_LOGIN` и `TEST_FAVORITES_PASSWORD` задают данные второго тестового аккаунта. Первый аккаунт создаёт объявления, второй добавляет их в избранное. Для проверки лимита избранное второго аккаунта должно быть пустым.

8. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -count=1 -v ./...
```

Для проверки повторяемости выполните команду три раза подряд.

9. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -count=1 -v ./tests/scenarios/myAdvertisement/
```

## Покрытие кода

Последний прогон на TestBoard покрыл 94,2% операторов Go-кода клиентов, менеджеров и тестовых хелперов (277 из 294).

```bash
go test -count=1 -coverpkg=./internal/...,./tests/... -coverprofile=/tmp/testboard-coverage.out ./...
go tool cover -func=/tmp/testboard-coverage.out
go tool cover -html=/tmp/testboard-coverage.out
```

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```