package main

import "github.com/gorilla/mux"

func main() {
	Info.Println("Initializing database")
	initializeDB()

	httpRouter := mux.NewRouter().StrictSlash(true)

	httpRouter.HandleFunc("/api/addemployee", apiAddEmployee).Methods("POST")
	httpRouter.HandleFunc("/api/updateschedule", apiAddEmployee).Methods("POST")

	httpRouter.HandleFunc("/api/removeemployee", apiAddEmployee).Methods("DELETE")
}
