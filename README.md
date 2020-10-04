# bank-transactions-go
Bank Transactions é uma aplicação escrita em Go que simula transações bancárias básicas, como: compra, saque e pagamento.

Tais serviços são expostos em uma API REST.

A aplicação possui uma arquitetura baseada nos conceitos de **Clean Architecture**, descrita originalmente por Robert C. Martin, tornando o código de fácil leitura, independente de agentes externos, altamente testável e de fácil manutenção. Clique [aqui](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) para ver mais detalhes no **blog do Uncle Bob**.

## Dependências
Para facilitar a execução bem como evitar problemas em diferentes ambientes, a aplicação é containerizada usando **Docker**. E, para facilitar a manipulação dos containers, usa-se **Docker-Compose**. Logo, estas duas ferramentas são pré-requisitos para executar este projeto com mais facilidade.

Por questões de segurança, nenhuma credencial de servidores está exposta nesse repositório, porém, a aplicação depende que estas credenciais estejam definidas em variáveis de ambiente. Estas credenciais serão automaticamente lidas pelo container Docker do arquivo **.env** na raíz do projeto. Tais credenciais estão disponíveis no  arquivo **.env.example**, com valores pré-definidos para o ambiente de desenvolvimento local. O comando **make init** criará o arquivo .env com base no exemplo e a aplicação estará pronta para iniciar.

## Como Iniciar
Após executar **make init** para definir as variáveis de ambiente, deve-se executar o comando **make up**, que fará o download de todas as dependências da aplicação e iniciará os containeres necessários para executar todos os casos de uso.  

## API REST
A API HTTP está exposta através da porta 8080.

Quando a solicitação não pode ser atendida, será retornado um *HTTP Status Code* condizente com a situação, e o payload conterá mais detalhes do(s) erro(s). Exemplo de payload de resposta com erro:
```
{
    "errors": [
        {
            "field": "document.number",
            "description": "number is a required field"
        }
    ]
}
```

### Criar Conta

Cada cliente possui uma conta disponibilizada pelo banco, e para criar a mesma, deve-se informar um CPF válido, formatado ou não.

Endpoint: 
```
POST /accounts
```
Headers:
```
Content-type: application/json
```
Request Payload:
```
{
    "document": {
        "number": "00000000191"
    }
}
```
Response:
```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sun, 04 Oct 2020 13:44:59 GMT
Content-Length: 81

{
  "id": 1,
  "document": {
    "number": "00000000191"
  },
  "created_at": "2020-10-04T13:44:59Z"
}
```

### Buscar Conta

Para buscar uma conta deve-se informar o ID da conta.

Endpoint: 
```
GET /accounts/{:id}
```
Response:
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 04 Oct 2020 13:56:20 GMT
Content-Length: 81

{
  "id": 1,
  "document": {
    "number": "00000000191"
  },
  "created_at": "2020-10-04T13:44:59Z"
}
```

### Registrar Transação

Para registrar uma transação deve-se informar o ID de uma conta válida, o ID da operação (ver tabela abaixo) e o valor da transação.

Operações:

|ID|Descrição|
| ------------- |:-------------:|
|1|Compra à vista|
|2|Compra parcelada|
|3|Saque|
|4|Pagamento|

Caso a operação informada seja de compra (1, 2) ou saque (3), a transação será registrada com valor negativo, enquanto transações de pagamento (4) serão registradas com valor positivo.

Endpoint: 
```
POST /transactions
```
Headers:
```
Content-type: application/json
```
Request Payload:
```
{
    "account_id": 1,
    "operation_id": 4,
    "amount": 100.00
}
```
Response:
```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sun, 04 Oct 2020 14:12:31 GMT
Content-Length: 132

{
    "id": 1,
    "account": {
        "id": 40,
        "document": {}
    },
    "operation": {
        "id": 4,
        "type": "PAGAMENTO"
    },
    "amount": 100,
    "created_at": "2020-10-04T11:35:58Z"
}
```
