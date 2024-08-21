package main

import (
   "encoding/json"
    "log"
    "net/http"
    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
    "github.com/swaggo/http-swagger"
    _ "gaminifica-senders/docs"
    "gaminifica-senders/internal/email"
    "gaminifica-senders/internal/rabbitmq"
    "gaminifica-senders/internal/worker"
    "github.com/streadway/amqp"
)

// @title  Sender API
// @version 1.0
// @description Esta é uma API para enviar e-mails.
// @host localhost:8080
// @BasePath /api
func main() {
        // Carregar o arquivo .env
        err := godotenv.Load()
        if err != nil {
            log.Fatalf("Erro ao processar o arquivo .env: %v", err)
        }

    r := mux.NewRouter()

    apiRouter := r.PathPrefix("/api").Subrouter()

    // Rotas da API
    apiRouter.HandleFunc("/send-email", sendEmailHandler).Methods("POST")
    // Rotas da API
    apiRouter.HandleFunc("/active-worker", activeWorker).Methods("GET")

    // Rota para a documentação Swagger
    r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)


    log.Println("Servidor iniciado na porta 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Não foi possível iniciar o servidor: %v", err)
    }
}

// @Summary Envia um e-mail
// @Description Envia um e-mail para o destinatário especificado
// @Accept json
// @Produce json
// @Param email body email.EmailRequest true "Dados do e-mail"
// @Success 200 {string} string "E-mail enviado com sucesso"
// @Failure 400 {string} string "Payload inválido"
// @Failure 500 {string} string "Erro ao enviar e-mail"
// @Router /send-email [post]
func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
    var req email.EmailRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Payload inválido", http.StatusBadRequest)
        return
    }

    if err := sendEmailToQueue(req); err != nil {
        http.Error(w, "Erro ao publicar a mensagem na fila", http.StatusInternalServerError)
        return
    }
    log.Println("Iniciando o worker de e-mails...")
    worker.StartEmailWorker()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("E-mail enviado para processamento"))
}

// @Summary Inicia o Worker
// @Description -
// @Success 200 {string} string "Worker Executado com sucesso"
// @Failure 500 {string} string "Erro ao executar worker"
// @Router /active-worker [get]
func activeWorker(w http.ResponseWriter, r *http.Request) {
    log.Println("Iniciando o worker de e-mails...")
    worker.StartEmailWorker(); 

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Fila Executada"))

}



// Função para conectar ao RabbitMQ e criar um canal
func connectToRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
    conn, err := rabbitmq.ConnectRabbitMQ()
    if err != nil {
        return nil, nil, err
    }

    ch, err := rabbitmq.CreateChannel(conn)
    if err != nil {
        conn.Close()
        return nil, nil, err
    }

    return conn, ch, nil
}

// Função para enviar um e-mail para a fila do RabbitMQ
func sendEmailToQueue(req email.EmailRequest) error {
    conn, ch, err := connectToRabbitMQ()
    if err != nil {
        return err
    }
    defer conn.Close()
    defer ch.Close()

    q, err := rabbitmq.DeclareQueue(ch, "email_queue")
    if err != nil {
        return err
    }

    body, err := json.Marshal(req)
    if err != nil {
        return err
    }

    return ch.Publish(
        "",        // Exchange
        q.Name,    // Nome da fila
        false,     // Mandatory
        false,     // Immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}


