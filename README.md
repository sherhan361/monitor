# go-musthave-devops-tpl

Шаблон репозитория для практического трека «Go в DevOps».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# Обновление шаблона

Чтобы получать обновления автотестов и других частей шаблона, выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-devops-tpl.git
```

Для обновления кода автотестов выполните команду:

(Для Unix систем)

```
git fetch template && git checkout template/main .github
```

(Для Windows PowerShell)

```
(git fetch template) -and (git checkout template/main .github)
```

Затем добавьте полученные изменения в свой репозиторий.



Запуск базы
``` shell script
 docker run --name monitor-db -e POSTGRES_PASSWORD=example -e POSTGRES_DB=monitor -p 5432:5432 postgres
```
