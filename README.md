# Balance microservice

Спасибо команде Авито за предоставленный проект.

После разработки на Java это первый такой крупный проект на Go.
Очень надеюсь, что я правильно выбрал архитектуру файлов.

### Проект

- Выбранная БД: Postgresql
- Версия Go - 1.19.2 darwin/amd64

В проекте выполнены, помимо основного, дополнительные задания, а также реализован сценарий
разрезервирования средств и возвращения на баланс пользователя.

Во время выполнения задания возник вопрос: откуда брать услуги и заказы. 
Я решил реализовать услуги, как будто они уже существуют, а айди заказов приходят
со стороннего микросервиса. 

При экспорте месячного отчета создается файл в папке files/csv/ с именем "год-месяц".
При повторном запросе этого месяца файл перезаписывается. Помимо сохранения файла в папке
предусмотрено скачивание файла при запросе на эндпоинт.


### Конфигурация приложения 

Чтобы конфигурировать приложение, нужно изменить данные в файле application.yml, 
который находится в директории configs.

#### Доступные конфигурации: 
1. База данных
   - host - хост БД
   - dbname - имя БД
   - sslmode (enable\disable) - включение или выключение sslmode
   - user - имя пользователя БД
   - password - пароль БД
   - port - порт БД
   
2. Конфигурация приложения
   - port - порт приложения


### Запуск приложения

Для запуска приложения я использую docker-compose. 
Создается два контейнера: БД и само приложение. 
Предусмотрена проверка на работоспособность БД:
пока БД не будет "живой", приложение не запуститься.

### Api

1. Health 
   - get health status: GET /health
2. Users
   - deposit on balance: POST /users/deposit
   - get balance of user: GET /users/{userId}
3. Accounting
   - get accounting records: GET /accounting
   - export to csv: GET /accounting/csv
4. Service 
   - Get service by id: GET /services/{serviceId}
5. Bill
   - reserve funds: POST /bills
   - approve reservation: PATCH /bills/{billId}/approve
   - reject reservation: PATCH /bills/{billId}/reject
6. Transaction
   - get transactions of user: GET /transactions/{userId}