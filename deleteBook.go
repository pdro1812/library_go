package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func deleteBook(w http.ResponseWriter, r *http.Request) {
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

	if b.ID == 0 {
		http.Error(w, "ID do livro é obrigatório para deletar", http.StatusBadRequest)
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

	// Deleta vínculos na tabela livro_autor
	_, err = tx.Exec("DELETE FROM livro_autor WHERE livro_id = $1", b.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Erro ao deletar vínculos livro_autor: %v", err), http.StatusInternalServerError)
		return
	}

	// Deleta o livro
	result, err := tx.Exec("DELETE FROM livro WHERE id = $1", b.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Erro ao deletar livro: %v", err), http.StatusInternalServerError)
		return
	}

	// Verifica se algum livro foi deletado
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Erro ao verificar deleção: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("Livro com id %d não encontrado", b.ID), http.StatusNotFound)
		return
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Erro ao commitar transação", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Livro com id %d e seus autores foram deletados com sucesso!", b.ID)))
}
