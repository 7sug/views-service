# views-service
Сервис для накрутки просмотров в телеграмм канал.

Поддерживаемый протокол: socks5

Эндпоинты: POST /views, где в теле передается ссылка на пост в канале (формат string); Ответ: Кол-во прокси, кол-во успешных запросов
           GET /ping ; Ответ: im alive - OK, другое - не ОК
           GET /test-parse ; Ответ: список успешно спаршеных прокси

Кол-во просмотров: от 10 до 50 (зависит от качества спаршеных прокси)

