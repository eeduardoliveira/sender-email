package middleware

import (
    "encoding/json"
    "net/http"
    "fmt"
)

var authServiceURL = "http://auth-service:8080/validate-token" // URL do microsserviço de login

// ValidateTokenResponse representa a resposta da validação do token
type ValidateTokenResponse struct {
    Valid bool `json:"valid"`
}

// JWTMiddleware valida o token com o serviço de autenticação
func JWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Autorização necessária", http.StatusUnauthorized)
            return
        }

        // Fazendo uma requisição para o serviço de login para validar o token
        token := authHeader[len("Bearer "):]
        valid, err := validateTokenWithAuthService(token)
        if err != nil || !valid {
            http.Error(w, "Token inválido ou não autorizado", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// validateTokenWithAuthService valida o token com o microsserviço de login
func validateTokenWithAuthService(token string) (bool, error) {
    req, err := http.NewRequest("POST", authServiceURL, nil)
    if err != nil {
        return false, err
    }

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return false, nil
    }

    var validateResp ValidateTokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
        return false, err
    }

    return validateResp.Valid, nil
}
