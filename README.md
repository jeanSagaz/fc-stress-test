## Dêe uma estreal! :star:
Se vc gostou do projeto 'go Cobra CLI', por favor dêe uma estrela

## Como executar:
Execute no prompt:  
docker run jeansagaz/go-cli:latest test --url=http://google.com --requests=1000 --concurrency=10  

Caso desejar rodar local execute no prompt de comando na pasta raiz:  
go run main.go test --url=http://google.com --requests=1000 --concurrency=10  
go run main.go test -u=http://google.com —r=1000 —c=10  

## Tecnologias implementadas:

go 1.20
 - Cobra CLI
 - Wait-Groups
 - Channels
 