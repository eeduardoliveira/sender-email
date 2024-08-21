# Gaminifica Sender

Este projeto é um exemplo de aplicação Go integrada com RabbitMQ. Ele usa Docker para facilitar o ambiente de desenvolvimento e execução.

## Requisitos

Certifique-se de ter os seguintes softwares instalados:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Configuração

1. Clone o repositório:

    ```bash
    git clone <URL_DO_SEU_REPOSITORIO>
    cd <NOME_DO_PROJETO>
    ```

2. Crie um arquivo `.env` na raiz do projeto (caso ele já não exista). Adicione as variáveis de ambiente necessárias para o seu projeto:

    ```dotenv
    # Exemplo de variáveis no arquivo .env
    EMAIL_FROM=
    SMTP_PORT=
    SMTP_USER=
    SMTP_PASSWORD=

    ```

3. Build e inicie os containers com Docker Compose:

    ```bash
    docker-compose up --build
    ```

   Esse comando vai:

   - Buildar o projeto Go.
   - Subir o RabbitMQ com a interface de gerenciamento.
   - Subir a aplicação Go.

4. Acesse a interface de gerenciamento do RabbitMQ:

    - URL: `http://localhost:15672`
    - Usuário: `guest`
    - Senha: `guest`

5. A aplicação Go estará rodando em:

    - URL: `http://localhost:8080`

## Estrutura do Projeto

```plaintext
.
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── .env
└── README.md

    - Dockerfile: Contém as etapas para buildar a aplicação Go e configurar o ambiente.
    - docker-compose.yml: Orquestra a aplicação Go e o RabbitMQ.
    - go.mod e go.sum: Gerenciamento de dependências do Go.
    - main.go: Código principal da aplicação Go.
    - .env: Configurações de variáveis de ambiente.
    - README.md: Instruções do projeto.

## Considerações

    - O arquivo .env deve estar presente na raiz do projeto para que as variáveis de ambiente sejam carregadas corretamente.
    - Certifique-se de que as portas usadas (8080, 5672, 15672) não estejam em uso por outros processos.