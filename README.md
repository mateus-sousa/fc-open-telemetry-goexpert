Passo a passo para execução da aplicação:

* Para subir as dependencias do projeto execute o comando:
```
    docker-compose up -d --build
```
A aplicação estará de pé respondendo nas seguintes portas:
* Web Service A: 8080
* Web Service B: 8081
* Zipkin: 9411

1. Para testar o Web Server, execute o arquivo get_weather.http que esta na pasta requests.

2. Para verificar o funcionamento do tracing, no seu navegador favorito, abra: http://localhost:9411/