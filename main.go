package main

import (
	"fmt"
	"net/http"
	
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	err := ConnectPostgres()
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}else{
		fmt.Println("Database connected successfully")
	}
	defer DB.Close()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/receiveAutors", receiveAutors)
	http.HandleFunc("/listAutors", listAutors)
	http.HandleFunc("/receiveBook", receiveBook)


	fmt.Println("Server running at http://localhost:8080/")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
