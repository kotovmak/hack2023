# hack2023

https://leaders2023.innoagency.ru

```plantuml
# Получение ответа в чате
@startuml user_story_1
actor User order 10
participant "Чат-бот" as Chatbot order 20
participant "Консультация" as Consultation order 30
User -> Chatbot : Вопрос пользователя
activate Chatbot
Chatbot -> Chatbot : Обработка ключевых слов
Chatbot -> User : Вывод ответа, ответ удовлетворил пользователя
deactivate Chatbot
@enduml
```

# Запись на консультацию
```plantuml
@startuml user_story_2
actor User order 10
participant "Чат-бот" as Chatbot order 20
participant "Консультация" as Consultation order 30
actor "Представитель КНО" as Officer order 40
alt Если запись через чат-бота
User -> Chatbot : Вопрос пользователя
activate Chatbot
Chatbot -> Chatbot : Обработка ключевых слов
Chatbot -> User : Вывод ответа и кнопка записи на консультацию
deactivate Chatbot
end
User -> Consultation : Заполнение формы записи на консультацию
activate Consultation
Consultation -> Officer : Запрос на подтверждение консультации
activate Officer
Officer -> Consultation : Подтверждение даты и времени консультации
Consultation -> Consultation : Ожидание даты и времени консультации
Consultation -> Officer : Напоминание, видеозвонок
Consultation -> User : Напоминание, видеозвонок
deactivate Consultation
@enduml
```