# Протокол для общения с сервером

## Типы методов
- 0x00 - getMe
- 0x01 - sendSensorsData
- 0x02 - setLEDColor
- 0x03 - sendLocationData
- 0x04 - sendBalance

## getMe (0x00)
Клиент -> Сервер -> Клиент
### Запрос к серверу
| Поле        | Количество байт |     Название      | Тип данных | Комментарий  |
|-------------|:---------------:|:-----------------:|:----------:|:------------:|
| METHOD TYPE |        1        |    Тип метода     |    byte    | 0x00 - getMe |
| DEVICE ID   |        2        |   ID устройства   |    int     |              |
| TOKEN       |       16        | Токен авторизации |   []byte   |              |
### Ответ от сервера (если успешно)
| Поле        | Количество байт |       Название        | Тип данных |  Комментарий   |
|-------------|:---------------:|:---------------------:|:----------:|:--------------:|
| METHOD TYPE |        1        |      Тип метода       |    byte    |  0x00 - getMe  |
| RESULT      |        1        | Результат авторизации |    byte    | 0x01 - успешно |
| DEVICE ID   |        2        |     ID устройства     |    int     |                |
| NAME        |       16        |  Название устройства  |   []char   |                |
### Ответ от сервера (если неуспешно)
| Поле        | Количество байт |       Название        | Тип данных |   Комментарий    |
|-------------|:---------------:|:---------------------:|:----------:|:----------------:|
| METHOD TYPE |        1        |      Тип метода       |    byte    |   0x00 - getMe   |
| RESULT      |        1        | Результат авторизации |    byte    | 0x00 - неуспешно |
| ERROR CODE  |        2        |      Код ошибки       |    int     |                  |
| DESCRIPTION |       32        |    Описание ошибки    |   []char   |                  |
## sendSensorsData (0x01)
Клиент -> Сервер
### Запрос к серверу
| Поле        | Количество байт |     Название      | Тип данных |      Комментарий       |
|-------------|:---------------:|:-----------------:|:----------:|:----------------------:|
| METHOD TYPE |        1        |    Тип метода     |    byte    | 0x01 - sendSensorsData |
| DEVICE ID   |        2        |   ID устройства   |    int     |                        |
| TOKEN       |       16        | Токен авторизации |   []byte   |                        |
| TEMPERATURE |        4        |    Температура    |   float    |                        |
| HUMIDITY    |        4        |     Влажность     |   float    |                        |
## setLEDColor (0x02)
Сервер -> Клиент
### Запрос к клиенту
| Поле        | Количество байт |          Название           | Тип данных |    Комментарий     |
|-------------|:---------------:|:---------------------------:|:----------:|:------------------:|
| METHOD TYPE |        1        |         Тип метода          |    byte    | 0x02 - setLEDColor |
| DEVICE ID   |        2        |        ID устройства        |    int     |                    |
| RED         |        1        | Яркость красного светодиода |    byte    |                    |
| GREEN       |        1        | Яркость зеленого светодиода |    byte    |                    |
| BLUE        |        1        |  Яркость синего светодиода  |    byte    |                    |
## sendLocationData (0x03)
Клиент -> Сервер
### Запрос к серверу
| Поле        | Количество байт |     Название      | Тип данных |       Комментарий       |
|-------------|:---------------:|:-----------------:|:----------:|:-----------------------:|
| METHOD TYPE |        1        |    Тип метода     |    byte    | 0x03 - sendLocationData |
| DEVICE ID   |        2        |   ID устройства   |    int     |                         |
| TOKEN       |       16        | Токен авторизации |   []byte   |                         |
| LATITUDE    |        4        |      Широта       |   float    |                         |
| LONGITUDE   |        4        |      Долгота      |   float    |                         |
| ALTITUDE    |        4        |      Высота       |   int32    |                         |
| ACCURACY    |        4        |     Точность      |   int32    |                         |
## sendBalance (0x04)
Клиент -> Сервер
### Запрос к серверу
| Поле        | Количество байт |     Название      | Тип данных |    Комментарий     |
|-------------|:---------------:|:-----------------:|:----------:|:------------------:|
| METHOD TYPE |        1        |    Тип метода     |    byte    | 0x04 - sendBalance |
| DEVICE ID   |        2        |   ID устройства   |    int     |                    |
| TOKEN       |       16        | Токен авторизации |   []byte   |                    |
| BALANCE     |        4        |      Баланс       |   float    |                    |

## Коды ошибок
401 - Unauthorized