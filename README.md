<p align="center">URLService</p>
<p align="center">
  <a href="https://travis-ci.org/github/x421/URLService"><img src="https://travis-ci.org/x421/URLService.svg" alt="Build Status"></a>
</p>
<hr>
## Сервис по сокращению ссылок.
База данных: MySQL
CI/CD: travis/heroku
Покрытие тестами составляет 82.1%
Приложение доступно по ссылке <a href="https://floating-plateau-23456.herokuapp.com/">https://floating-plateau-23456.herokuapp.com/</a>

## Описание контроллеров:
<hr>
  ### /setShort : метод Post, тип принимаемых данных JSON
  Пример входных данных
  ```json
  {
      "URI": "http://google.com",
      "Short": ""
  }
  ```
  где URI - адрес перенаправления(не более 200 символов), Short - желаемый адрес на сервере, при передаче пустого значения генерируется автоматически(не более 25 латинских букв или цифр и знака -). 

  ### /[ссылка либо ничего] : Любой метод, без входных параметров
  Отображает главную страницу либо перенаправляет на извлекаемую из базы ссылку, если такая имеется
