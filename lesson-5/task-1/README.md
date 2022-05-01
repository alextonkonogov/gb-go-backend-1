# gb-go-backend-1
## lesson-5
### Task-1
_Выбрать роутер для курсового проекта. Приложить ссылки на код или Pull Request. Добавить в
README проекта обоснование выбора роутера._

Изначально на первом или втором занятии при создании первой версии проекта мною был выбран [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter). 
Я пару лет назад читал книжку [Building web apps with Go. Jeremy Saenz](https://codegangsta.gitbooks.io/building-web-apps-with-go/content/index.html) и там он рекомендовался именно как очень простой и легкий. 
Плюс, в отличие от стандартного, он умеет следующие важные для меня возможности:
* Указывать метод HTTP запроса при описании роута (`r.POST("/short", a.ShortURL)`)
* Использовать именованные параметры в составе URL (`r.GET("/s/:ID/:code", a.LongToShort)`)

Из недостатков отметил, что он просит хендлер с третьим параметром. Мне же хочется, чтобы хендлеры принимали стандартные `w http.ResponseWriter и r *http.Request`

``` go 

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
```

Поэтому, еще раз пересмотрев лекцию, я остановился на роутере [go-chi](https://github.com/go-chi/chi):
* Он имеет все те же возможности, что и [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) + массу дополнительных
* Он идиоматичный и легкий
* Он чаще используется в настоящих проектах и более популярен (хотя у [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) больше звезд :) )
* До сих пор обновляется и поддерживается (последняя версия была от 19 ноября прошлого года)
* Он помог мне решить вопрос с правильным получением IP адреса клиента
* Это один из роутеров, которые рекомендовал преподаватель

Версия проекта с роутером go-chi - https://github.com/alextonkonogov/gb-go-url-shortener/pull/2