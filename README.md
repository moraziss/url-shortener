# URL Shortener

Простой сервис на Go для коротких ссылок (URL shortener), который перенаправляет пользователя с короткого пути на полный URL. Пути и соответствующие ссылки можно задавать как через Go-карту (map), так и через YAML-конфигурацию.

## Возможности

- Редирект с короткого пути на полный URL (`http.Redirect`)
- Конфигурация маршрутов через `map[string]string` в коде
- Конфигурация маршрутов через YAML-файл (`gopkg.in/yaml.v3`)
- Fallback-хендлер: если путь не найден среди коротких ссылок, запрос передаётся дальше по цепочке

## Структура проекта

```
url-shortener/
├── go.mod              # module github.com/moraziss/url-shortener
├── main/
│   └── main.go          # точка входа, настройка маршрутов
└── urlshort/
    └── handler.go       # MapHandler, YAMLHandler и парсинг YAML
```

## Установка

```bash
git clone https://github.com/moraziss/url-shortener.git
cd url-shortener
go mod tidy
go run main/main.go
```

## Как это работает

### MapHandler

Принимает `map[string]string` (путь → URL) и `http.Handler` в качестве fallback. Возвращает `http.HandlerFunc`, который:

1. Проверяет, есть ли текущий путь запроса (`r.URL.Path`) в карте.
2. Если есть — делает редирект (`http.StatusFound`) на соответствующий URL.
3. Если нет — передаёт запрос в fallback-хендлер.

### YAMLHandler

Принимает YAML-конфигурацию в виде `[]byte` и fallback-хендлер. Внутри:

1. Парсит YAML в список пар `path` / `url`.
2. Преобразует этот список в `map[string]string`.
3. Строит `MapHandler` на основе полученной карты.

Пример YAML-конфигурации:

```yaml
- path: /urlshort
  url: https://github.com/moraziss/url-shortener
- path: /urlshort-final
  url: https://github.com/moraziss/url-shortener/tree/final
```

## Пример запуска

```
Starting the server on port 8080
```

- `http://localhost:8080/urlshort` → редирект на `https://github.com/moraziss/url-shortener`
- `http://localhost:8080/urlshort-final` → редирект на ветку `final` репозитория
- `http://localhost:8080/urlshort-godoc` → редирект на godoc-страницу проекта (задаётся через карту `pathsToUrls` в `main.go`)
- `http://localhost:8080/yaml-godoc` → редирект на godoc `gopkg.in/yaml.v2`
- Любой другой путь (например `/`) → обрабатывается дефолтным хендлером `hello`, который выводит `Hello, World!`

Цепочка фолбэков: `yamlHandler → mapHandler → defaultMux (hello)`.
