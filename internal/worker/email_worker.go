package worker

import (
    "log"
    "gaminifica-senders/internal/email"
    "gaminifica-senders/internal/rabbitmq"
)

func ProcessEmailRequest(req email.EmailRequest) error {
    err := email.SendEmail(req.To, req.Subject, req.Body)
    if err != nil {
        log.Printf("Erro ao enviar e-mail: %v", err)
        return err
    }
    log.Printf("E-mail enviado para %s", req.To)
    return nil
}

func StartEmailWorker() {
    err := rabbitmq.ConsumeQueue("email_queue", ProcessEmailRequest)
    if err != nil {
        log.Fatalf("Erro ao consumir fila: %v", err)
    }
}
