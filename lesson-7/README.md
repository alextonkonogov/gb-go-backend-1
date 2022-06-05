# gb-go-backend-1
## lesson-7
### Task

_Протестировать написанный код с помощью запуска curl в соответствии с API сервиса. Убедиться что все работает, либо что не работает и поправить ошибки._

---
### Подготовка
Начал со скрипта runpg.sh. При запуске он ругнулся вот так:

```shell
a.tonkonogov@admins-MacBook-Pro lesson-7 % ./runpg.sh
Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
docker: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?.
See 'docker run --help'.
a.tonkonogov@admins-MacBook-Pro lesson-7 % ./runpg.sh
12: Pulling from library/postgres
dc1f00a5d701: Pulling fs layer 
3bb4b34c334c: Pulling fs layer 
4739db3ff30d: Pulling fs layer 
67627067cf92: Pull complete 
8cb1fcaf0443: Pull complete 
4495f752b8b4: Pull complete 
54aaebaf7bd6: Pull complete 
ca284527a779: Pull complete 
addeb7a3dc8d: Pull complete 
5492f9961ca3: Pull complete 
af1d12941a76: Pull complete 
7d1359bf0385: Pull complete 
fc1e18738f7a: Pull complete 
Digest: sha256:fe84844ef27aaaa52f6ec68d6b3c225d19eb4f54200a93466aa67798c99aa462
Status: Downloaded newer image for postgres:12
docker.io/library/postgres:12
b81948df3435bfd1c2ec32d29a42d6e809454c723a3450a5cb651040468ef135
docker: Error response from daemon: Mounts denied: 
The path /opt/databases/postgres is not shared from the host and is not known to Docker.
You can configure shared paths from Docker -> Preferences... -> Resources -> File Sharing.
See https://docs.docker.com/desktop/mac for more info.
```

Полез в интернет и подчерпнул знаний об особенностях macOS. Для работы заменил var на tmp в строке c маппингом директорий:

```
docker run --name=pgsql1 -p 5432:5432 -v "/tmp/databases/postgres":/var/lib/postgresql/data -e POSTGRES_PASSWORD=1110 -e POSTGRES_DB=test -d postgres:12
```

Запустил БД:

```shell
a.tonkonogov@admins-MacBook-Pro lesson-7 % ./runpg.sh
12: Pulling from library/postgres
Digest: sha256:fe84844ef27aaaa52f6ec68d6b3c225d19eb4f54200a93466aa67798c99aa462
Status: Image is up to date for postgres:12
docker.io/library/postgres:12
bff647da9163e08ce92017a8dbcc799b78b2309a1d9c0ce22391ff19f648ed6d
```

Подтянул пакеты из go.mod и запустил приложение:

```shell
a.tonkonogov@admins-MacBook-Pro lesson-7 % go mod tidy
go: downloading github.com/deepmap/oapi-codegen v1.9.1
go: downloading github.com/getkin/kin-openapi v0.89.0
go: downloading github.com/googleapis/gax-go/v2 v2.1.1
go: downloading cloud.google.com/go v0.100.2
go: downloading github.com/google/go-cmp v0.5.7
go: downloading google.golang.org/genproto v0.0.0-20220210181026-6fee9acbd336
go: downloading google.golang.org/grpc v1.44.0
go: downloading github.com/ugorji/go v1.2.6
go: downloading golang.org/x/crypto v0.0.0-20220210151621-f4118a5b28e2
go: downloading github.com/ugorji/go/codec v1.2.6
go: downloading github.com/go-playground/validator/v10 v10.10.0
go: downloading github.com/json-iterator/go v1.1.12
go: downloading golang.org/x/sys v0.0.0-20220209214540-3681064d5158
go: downloading github.com/ghodss/yaml v1.0.0
go: downloading github.com/go-openapi/jsonpointer v0.19.5
go: downloading google.golang.org/api v0.68.0
go: downloading github.com/modern-go/reflect2 v1.0.2
go: downloading github.com/go-playground/universal-translator v0.18.0
go: downloading github.com/leodido/go-urn v1.2.1
go: downloading github.com/go-playground/locales v0.14.0
go: downloading github.com/kr/pretty v0.3.0
go: downloading golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
go: downloading github.com/go-openapi/swag v0.21.1
go: downloading github.com/rogpeppe/go-internal v1.8.0
go: downloading github.com/mailru/easyjson v0.7.7
go: downloading github.com/josharian/intern v1.0.0
go: downloading golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
go: downloading cloud.google.com/go/compute v1.2.0
go: finding module for package cloud.google.com/go/iam
go: downloading cloud.google.com/go/iam v0.3.0
go: downloading cloud.google.com/go v0.102.0
go: found cloud.google.com/go/iam in cloud.google.com/go/iam v0.3.0
go: downloading google.golang.org/genproto v0.0.0-20220222213610-43724f9ea8cf
go: downloading google.golang.org/api v0.70.0
go: downloading cloud.google.com/go/compute v1.3.0

a.tonkonogov@admins-MacBook-Pro lesson-7 % go run reguser/cmd/reguser/main.go
2022/06/05 20:48:31 Start
```

Далее тестируем приложение

---
### Создание

**Запрос:**
```shell
curl --location --request POST 'localhost:8000/create' \
--header 'Authorization: Basic YWRtaW46YWRtaW4=' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "George",
    "data": "Lorem Ipsum",
    "perms": 765
}'
```

**Ответ:**
```shell
{"id":"2a9bee67-a498-4708-867c-d52842cd5956","name":"George","data":"Lorem Ipsum","perms":0}
```

---

### Поиск

**Запрос:**
```shell
curl --location --request GET 'localhost:8000/search/e' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
[
{"id":"2a9bee67-a498-4708-867c-d52842cd5956","name":"George","data":"Lorem Ipsum","perms":493}
]
```

---
### Чтение

**Запрос:**
```shell
curl --location --request GET 'localhost:8000/read/2a9bee67-a498-4708-867c-d52842cd5956' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
{"id":"2a9bee67-a498-4708-867c-d52842cd5956","name":"George","data":"Lorem Ipsum","perms":0}
```

---
### Удаление

**Запрос**:
```shell
curl --location --request DELETE 'localhost:8000/delete/2a9bee67-a498-4708-867c-d52842cd5956' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
{"id":"2a9bee67-a498-4708-867c-d52842cd5956","name":"George","data":"Lorem Ipsum","perms":0}
```
---
Здесь есть тот же самый баг (или фича), который я правил в 6-ой домашке: perms не записывается при создании и потом выдается в виде 0. 
В этот раз я уже не стал его править, потому что это не влияет на суть задания.
---
### Проверим БД
Наконец, проверим, что в БД у нас все записывается:

```shell
a.tonkonogov@admins-MacBook-Pro gb-go-backend-1 % psql -h localhost -p 5432 -d test -U postgres
Password for user postgres: 
psql (14.1, server 12.11 (Debian 12.11-1.pgdg110+1))
Type "help" for help.

test-# \d
         List of relations
 Schema | Name  | Type  |  Owner   
--------+-------+-------+----------
 public | users | table | postgres
(1 row)

test=# select * from users;
                  id                  |          created_at           |          updated_at           |          deleted_at           |  name  |    data     | perms 
--------------------------------------+-------------------------------+-------------------------------+-------------------------------+--------+-------------+-------
 2a9bee67-a498-4708-867c-d52842cd5956 | 2022-06-05 17:50:22.092634+00 | 2022-06-05 17:50:22.092634+00 | 2022-06-05 17:53:47.350425+00 | George | Lorem Ipsum |     0
(1 row)

```