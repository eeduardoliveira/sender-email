package email

import (
    "log"
    "os"
    "strconv"
    "gopkg.in/gomail.v2"
    "github.com/joho/godotenv"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Erro ao processar o arquivo .env: %v", err)
    }
}

func SendEmail(to string, subject string, body string) error {
    from := os.Getenv("EMAIL_FROM")
    smtpServer := os.Getenv("SMTP_SERVER")
    smtpPort := os.Getenv("SMTP_PORT")
    smtpUser := os.Getenv("SMTP_USER")
    smtpPassword := os.Getenv("SMTP_PASSWORD")

    m := gomail.NewMessage()
    m.SetHeader("From", from)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    port, err := strconv.Atoi(smtpPort)
    if err != nil {
        return err
    }

    d := gomail.NewDialer(smtpServer, port, smtpUser, smtpPassword)

    if err := d.DialAndSend(m); err != nil {
        return err
    }
    log.Printf("Email enviado com sucesso para %s", to)
    return nil
}
