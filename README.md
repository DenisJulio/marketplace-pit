# Marketplace PIT

## Dependências necessárias

Certifique-se de que as seguintes ferramentas estão instaladas no seu ambiente:

* [golang](https://go.dev/doc/install) 
* [docker](https://docs.docker.com/get-started/get-docker/)
* [templ](https://templ.guide/quick-start/installation/)

## Intruções para executar a aplicação

Antes de começar, certifique-se de instalar as dependências do projeto usando:

```sh
go mod tidy
```

1. Iniciar o container com o banco de dados:

```sh
docker run -d \
  --name db \
  --shm-size=128mb \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_USER=denis \
  -e POSTGRES_DB=db \
  -v $(pwd)/sql:/docker-entrypoint-initdb.d \
  -p 5432:5432 \
  postgres:15
```

2. Compilar os templates:

```sh
templ generate
```
3. Compilar a aplicação:

```sh
go build -o ./tmp/main .
```
4. Executar a aplicação:

```sh
./tmp/main
```
## Removendo o container do banco de dados

Para parar e remover o container do banco de dados após o uso:

```sh
docker stop db && docker rm db
```
