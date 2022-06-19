Имеется два REST маршрута:
  1. /receive?GUID=xxxx - для пользователя с GUID-xxxx выдается пара токенов (Access, Refresh)
  
  ![image](https://user-images.githubusercontent.com/98939753/174469225-1c277283-0031-407b-a034-3dc249ab4679.png)
  
  Access token - тип JWT, алгоритм SHA512,
  Refresh token - тип JWT, в формате base64.
  
  Refresh token переводится в bcrypt хеш и сессия записывается в БД.
  ![image](https://user-images.githubusercontent.com/98939753/174469219-b6b5a3e2-64a5-45d7-bff9-2cd7aa8975c0.png)
  
  Д ля одного пользователя возможно открыть неограниченное количество сессий (можно ограничить).
  При повторной попытке получить токены, выдается новая пара токенов. 
    Причем, чтобы обновить Acess Token для второй сессии, необходимо использовать тот Refresh Token, с которым
  выдавался Access Token, т.к. они свзяаны с помощью уникального идентификатора. 
      
  2. /refresh?GUID=xxxx - выполняется обновления токенов, Refresh token передается в теле запроса.
  
  ![image](https://user-images.githubusercontent.com/98939753/174469400-b40e73b2-a538-41e3-abc5-b657517072ea.png)
  В теле ответа возвращается новая пара токенов. Далее в бд создается новая сессия, старая удаляется. 
  
  Если токен был изменен или происходит попытка обновить токены для пользователя, который их не получал, то выводится текст о наличии ошибки.
  
  ![image](https://user-images.githubusercontent.com/98939753/174469654-5c74d06b-f56d-46b3-8c50-56a613c26b30.png)
  
  или
  
  ![image](https://user-images.githubusercontent.com/98939753/174469664-4f250798-6cf3-4e44-8cec-baac7b0747b8.png)
