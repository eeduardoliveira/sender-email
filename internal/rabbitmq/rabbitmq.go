package rabbitmq

import (
    "encoding/json"
    "log"
    "github.com/streadway/amqp"
    "gaminifica-senders/internal/email"
)

// ConnectRabbitMQ estabelece uma conexão com o RabbitMQ.
func ConnectRabbitMQ() (*amqp.Connection, error) {
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
    if err != nil {
        return nil, err
    }
    return conn, nil
}

// CreateChannel cria um novo canal no RabbitMQ.
func CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    return ch, nil
}

// DeclareQueue declara uma fila no RabbitMQ.
func DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
    q, err := ch.QueueDeclare(
        queueName, // Nome da fila
        false,     // Durable
        false,     // Auto-delete
        false,     // Exclusive
        false,     // No-wait
        nil,       // Args
    )
    if err != nil {
        return amqp.Queue{}, err
    }
    return q, nil
}

// ConsumeQueue consome mensagens de uma fila e as processa usando o handler fornecido.
func ConsumeQueue(queueName string, handler func(email.EmailRequest) error) error {
    conn, err := ConnectRabbitMQ()
    if err != nil {
        return err
    }
    defer conn.Close()

    ch, err := CreateChannel(conn)
    if err != nil {
        return err
    }
    defer ch.Close()

    q, err := DeclareQueue(ch, queueName)
    if err != nil {
        return err
    }

    msgs, err := ch.Consume(
        q.Name, // Nome da fila
        "",     // Nome do consumidor
        true,   // Auto-ack
        false,  // Exclusive
        false,  // No-local
        false,  // No-wait
        nil,    // Args
    )
    if err != nil {
        return err
    }

    // Processa as mensagens recebidas
    for msg := range msgs {
        var req email.EmailRequest
        if err := json.Unmarshal(msg.Body, &req); err != nil {
            log.Printf("Erro ao desserializar mensagem: %v", err)
            continue
        }

        if err := handler(req); err != nil {
            log.Printf("Erro ao processar a mensagem: %v", err)
            msg.Nack(false, false) 
        } else {
            msg.Ack(false)
        }    
    }
    return nil
}
