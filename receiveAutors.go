package main

import (
    "encoding/csv"
    "fmt"
    "net/http"
	"log"
)

func receiveAutors(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Erro ao processar o formulário", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo CSV", http.StatusInternalServerError)
		return
	}

	for _, record := range records[1:] {
	_, err = DB.Exec("INSERT INTO autor (nome) VALUES ($1)", record[0])
	if err != nil {
		log.Fatalf("Erro ao inserir autor: %v", err)
	}
	fmt.Println("Usuário inserido com sucesso!")


		fmt.Println(record)
	}	

    w.WriteHeader(http.StatusOK)
	w.Write([]byte("Arquivo recebido e processado com sucesso!"))
}
