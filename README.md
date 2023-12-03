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


# Просмотр документации

1. `godoc -http=:8080`
2. http://localhost:8080/pkg/github.com/tiraill/go_collect_metrics/?m=all
3. wget -r -np -E -p -k -nH -P doc "http://localhost:8080/pkg/github.com/tiraill/go_collect_metrics/?m=all"

# Запуск линтеров

1. `go build cmd/staticlint/mycheck.go`
2. `./mycheck ./...`

# Сборка приложений

### agent
`go build -ldflags "-X main.buildVersion=v0.1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(git log --pretty=format:'%h' -n1)" cmd/agent/agent.go`

### server
`go build -ldflags "-X main.buildVersion=v0.1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(git log --pretty=format:'%h' -n1)" cmd/server/server.go`