package email

// EmailRequest representa os dados necessários para enviar um e-mail.
type EmailRequest struct {
    To      string `json:"to" example:"recipient@example.com"`
    Subject string `json:"subject" example:"Hello World"`
    Body    string `json:"body" example:"Este é um exemplo"`
}
