# О проекте

Индивидуальный дипломный проект курса «Go-разработчик»

# Шаблон

Шаблон взят из репозитория <https://github.com/yandex-praktikum/go-musthave-diploma-tpl>

# Спецификация

Спецификация проекта находится в файле [SPECIFICATION.md](https://github.com/k0st1a/gophermart/blob/master/SPECIFICATION.md)

# Обновление автотестов

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.