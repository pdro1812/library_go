package main

import (
	"encoding/json"
    "net/http"
	"log"
)

func listAutors(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }


	rows, err := DB.Query("SELECT * FROM autor")
	if err != nil {
		log.Fatalf("Erro ao executar SELECT: %v", err)
	}
	defer rows.Close()

	type Autor struct {
		ID   int    `json:"id"`
		Nome string `json:"nome"`
	}

	var autors []Autor

	for rows.Next() {
		var autor Autor
		if err := rows.Scan(&autor.ID, &autor.Nome); err != nil {
			http.Error(w, "Erro ao ler os dados do autor", http.StatusInternalServerError)
			log.Printf("Erro ao ler linha: %v", err)
			return
		}
		autors = append(autors, autor)
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(autors)

}
