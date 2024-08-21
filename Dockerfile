# Etapa 1: Construir o binário Go
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copiar todos os arquivos para o container de build, incluindo o .env
COPY . .

# Instalar dependências e construir o binário
RUN go mod tidy
RUN go build -o main

# Etapa 2: Container final para execução
FROM alpine:latest

WORKDIR /app

# Copiar o binário e o arquivo .env do estágio de build para o container final
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Definir o comando de execução
CMD [ "/app/main" ]
