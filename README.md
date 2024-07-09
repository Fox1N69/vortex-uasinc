# Текстовое задание для Effective Mobile

В данном проекте я реализовал api для task-tracker c использованием сторонего api для получения доп информации о пользователях.
#### Документация к api  - localhost:4000/swagger/index.html#/

## Запуск приложения
Перед запуском приложения поменяйте конфигурационные данные в файлe который находится по пути [./config/config.json](./config/config.json)


Запуск миграции зависимостей

```console
make dep
```

Запуск проекта локально

```console
make run-tracker
```

Запуск тестов

```console
make test
```

Сборка проекта

```console
make build-tacker
```

Запуск с hot reload

```console
air
```


## Используемые библиотеки

| Library    | Usage             |
| ---------- | ----------------- |
| gin        | Base framework    |
| database/sql | SQL library       |
| postgres   | Database          |
| logrus     | Logger library    |
| viper      | Config library    |

## Комментарий
В проекте подразумевалось использованние Redis для кэширования частых запросов, но из-за отсуцтвия уверености, что у проверяющего будут желани заморачиватся с подключение, было принято решение оставить эту **идею**.