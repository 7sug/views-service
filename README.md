# views-service
Сервис для накрутки просмотров в телеграмм канал.

Поддерживаемый протокол: socks5

Эндпоинт: POST /views, где в теле передается ссылка на пост в канале (формат string)
Ответ: Кол-во прокси, кол-во успешных запросов

Кол-во просмотров: от 10 до 50 (зависит от качества спаршеных прокси)

