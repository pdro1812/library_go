package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
)


func listBook(w http.ResponseWriter, r *http.Request) {
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

	// Montar query dinamicamente
	query := "SELECT id, nome, edicao, publication_year FROM livro WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if b.ID != 0 {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, b.ID)
		argIndex++
	}
	if b.Nome != "" {
		query += fmt.Sprintf(" AND nome ILIKE $%d", argIndex)
		args = append(args, "%"+b.Nome+"%")
		argIndex++
	}
	if b.Edition != 0 {
		query += fmt.Sprintf(" AND edicao = $%d", argIndex)
		args = append(args, b.Edition)
		argIndex++
	}
	if b.Year != 0 {
		query += fmt.Sprintf(" AND publication_year = $%d", argIndex)
		args = append(args, b.Year)
		argIndex++
	}

	// Executa query
	rows, err := DB.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao consultar livros: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var livros []Book

	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Nome, &book.Edition, &book.Year)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao ler livro: %v", err), http.StatusInternalServerError)
			return
		}

		// Busca autores desse livro
		autorRows, err := DB.Query("SELECT autor_id FROM livro_autor WHERE livro_id = $1", book.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao consultar autores do livro %d: %v", book.ID, err), http.StatusInternalServerError)
			return
		}

		var autores []int
		for autorRows.Next() {
			var autorID int
			err = autorRows.Scan(&autorID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Erro ao ler autor_id: %v", err), http.StatusInternalServerError)
				return
			}
			autores = append(autores, autorID)
		}
		autorRows.Close()

		book.Autores = autores
		livros = append(livros, book)
	}

	// Retorna a lista de livros como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(livros)
}
