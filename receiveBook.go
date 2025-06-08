package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	
)

type Book struct {
	ID      int    `json:"id"`
	Nome    string `json:"nome"`
	Edition int    `json:"Edition"`
	Year    int    `json:"Year"`
	Autores []int  `json:"autores"`
}

func receiveBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var b Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Erro ao receber JSON", http.StatusBadRequest)
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		http.Error(w, "Erro ao iniciar transação", http.StatusInternalServerError)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			http.Error(w, "Erro inesperado", http.StatusInternalServerError)
		}
	}()

	_, err = tx.Exec("INSERT INTO livro (id, nome, edicao, publication_year) VALUES ($1, $2, $3, $4)", b.ID, b.Nome, b.Edition, b.Year)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Erro ao inserir livro: %v", err), http.StatusInternalServerError)
		return
	}

	for _, autorID := range b.Autores {
		_, err = tx.Exec("INSERT INTO livro_autor (livro_id, autor_id) VALUES ($1, $2)", b.ID, autorID)
		if err != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("Erro ao inserir livro_autor (autor_id=%d): %v", autorID, err), http.StatusInternalServerError)
			return
		}
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Erro ao commitar transação", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Livro e autores inseridos com sucesso!"))
}