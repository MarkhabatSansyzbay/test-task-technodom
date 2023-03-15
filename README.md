## Описание проблемы

Периодически на сайте меняется роутинг, например: ссылка на категорию товаров выглядела следующим образом - /sad-dacha-ogorod/tehnika-dlja-sada-i-doma/izmel-chiteli-sadovye, потом данной категории поменяли родительскую категорию, и ссылка стала выглядеть следующим образом - /sad-dacha-ogorod/tehnika-dlja-sada-i-doma/tehnika-dlja-sada-i-doma/izmel-chiteli-sadovy.

 Такие вещи происходят довольно часто, и так как в поисковых системах таких как Google, Yandex ссылка на страницу может быть проиндексирована все еще в старом виде, требуется редайректить пользователей с старых ссылок на новые.

Задание состоит из пунктов обязательных для реализации и опциональных.
Можно использовать любое решение на свой выбор: базы данных BoltDB, SQLite, или можно развернуть любую СУБД ( MySQL, PostgreSQL, MongoDB и тд ).
Иметь установленный VSCode, Goland или подобный кодовый редактор.

## Обязательные требования:

1. Требуется реализовать веб сервис - REST API, который будет предоставлять возможность управления условиями редайретов которые должны храниться в базе данных, а также сам метод перенаправления ссылок.

Должны быть реализованы следующие методы:

- Методы для управления данными
    • GET /admin/redirects - метод для получения списка записей из базы ( нужно учесть пагинацию, так как в базе хранится довольно много объектов )
    • GET /admin/redirects/{id} - метод для получения конкретного объекта из базы
    • POST /admin/redirects - метод для создания новой записи в базе
    • PATCH /admin/redirects/{id} - метод для записи текущей актуальной ссылки в список исторических, и изменения активной ссылки
    • DELETE /admin/redirects/{id} - метод для удаления записи из базы
	
- Пользовательские методы:
    • GET /redirects - метод для перенаправления пользовательского запроса на корректную ссылку. Метод должен принимать в query параметрах ссылку, по которой пользователь пытается запросить данные. Метод должен возвращать статус 200 если пользовательская ссылка является активной, или вернуть статус 301 с новой ссылкой, если требуется перенаправить запрос


2. Данные нужно хранить в базе данных. Можно использовать любое решение на свой выбор: in Memory базы для простоты: BoltDB, SQLite, или же можно развернуть любую СУБД ( MySQL, PostgreSQL, MongoDB и тд ).

    Файл с датасетом можете найти по ссылке: https://drive.google.com/file/d/1dTCcIDwHWnejJTMZWbSngTxtCOoAdQ9v/view?usp=sharing

    Объект выглядит следующим образом:
	{
		“active_link”: “какая-то активная ссылка”,
		“history_link”: какая-то старая сОписание проблемы:
    }  


3. Требуется имплементировать простой inMemory cache не используя сторонние библиотеки. Кэш должен быть использован для сохранения запрашиваемых ссылок и ссылок куда нужно перенаправлять пользователя в памяти, для того чтобы постоянно не обращаться к базе данных. Так же размер кэша должно быть ограниченным ( например, в кэше должно храниться не больше 1000 ключей ):

	Требуется имплементировать следующий интерфейс:
```
type Cache interface {
    // Метод для добавление новой записи в кэш или обновления уже существующей
    Add(key, value string)
    // Метод для получения значения из кэша по ключу. Если значения нет, вернуть пустую строку и false
    Get(key string) (value string, ok bool)
    // Метод для получения количества ключей в кэше
    Len() int
}
```

## Запуск 
```
make psqlrun
```
and 
```
make run
```
