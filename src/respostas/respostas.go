package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON retorna uma respos em JSON para a requisição
func JSON(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(dados); err != nil {
		log.Fatal(err)
	}
}

// Erro retorna um erro em formato JSON
func Erro(w http.ResponseWriter, statusCode int, erro error) {
	JSON(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})
}
