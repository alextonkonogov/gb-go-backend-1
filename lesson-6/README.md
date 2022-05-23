# gb-go-backend-1
## lesson-6
### Task

_Протестировать написанный код с помощью запуска curl в соответствии с API сервиса. Убедиться что все работает, либо что не работает и поправить ошибки._

Запустил, все работает. Был баг (или фича), что perms не записывался при создании и потом в виде 0 выдавался - я поправил. 
Все остальное, вроде, как и задумано.
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
{"id":"64cb6bf6-0fd6-4900-9839-9b8bfdc6bcd9","name":"George","data":"Lorem Ipsum","perms":765}
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
{"id":"2d321745-f443-405f-bd50-561bb9dbeeba","name":"Alex","data":"Lorem Ipsum","perms":0}
,
{"id":"64cb6bf6-0fd6-4900-9839-9b8bfdc6bcd9","name":"George","data":"Lorem Ipsum","perms":765}
,
{"id":"5669e804-174c-4e54-a68a-82bfb2d72e14","name":"Michael","data":"Lorem Ipsum ","perms":0}
]

```

---
### Чтение

**Запрос:**
```shell
curl --location --request GET 'localhost:8000/read/64cb6bf6-0fd6-4900-9839-9b8bfdc6bcd9' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
{
    "id": "64cb6bf6-0fd6-4900-9839-9b8bfdc6bcd9",
    "name": "George",
    "data": "Lorem Ipsum",
    "perms": 765
}
```

---
### Удаление

**Запрос**:
```shell
curl --location --request DELETE 'localhost:8000/delete/64cb6bf6-0fd6-4900-9839-9b8bfdc6bcd9' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
{
    "id": "64cb6bf6-0fd6-4900-9839-9b8bfdc6bcd9",
    "name": "George",
    "data": "Lorem Ipsum",
    "perms": 765
}
```
