# Проект по сокращению ссылок на golang

Этот учебный проект [Ссылка на пример работы](https://tiny-url-nu.vercel.app/) представляет собой веб приложение, написанное на языке Go предназначенное для сокращения ссылок и создания qr кодов. При сокращении ссылка записывается в базу данных в виде :
`сгенерированный ключ` - сокращенная ссылка, 
`значение` исходная ссылка
вид таблицы `urls`:

| `id` int4 | `key` text | `value` text | 'count' int4 default 0 |
| --------- | ---------- | ------------ |---|
|          |            |              | |

после при переходе на сокращенную ссылку веб приложение ищет ее значение в базе данных и переадресует пользователя

## Использованные библиотеки:
- [github.com/jackc/pgx/v4](https://github.com/jackc/pgx/v4)
- [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2)
- [github.com/gofiber/adaptor/v2](https://github.com/gofiber/adaptor/v2)
- [github.com/ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- [github.com/skip2/go-qrcode](https://github.com/skip2/go-qrcode)

Проект также развертывается на [Vercel](https://vercel.com/) с использованием файла `vercel.json`, который настроен следующим образом:

```json
{
    "rewrites": [
      { "source": "(.*)", "destination": "api/index.go" }
    ]
}
```

## Секретные данные
Для корректной работы с базой данных и интеграции Vercel с GitHub Secrets, необходимо добавить следующие переменные в GitHub Secrets:
```plaintext
VERCEL_TOKEN
DB_PORT
DB_HOST
DB_NAME
DB_USER
DB_PASSWORD
```
(при локальной работе нужно заменить функцию в файле /config/config.go а данные записать в файл .env)

## Структура проекта
```plaintext
Project:.
├───api
│   └───index.go   // Главный файл API
├───config
│   └───config.go   // Конфиг базы данных
├───db
│   └───operations   // Файлы с операциями с базой данных
└───public   // html css файлы для проекта
```


