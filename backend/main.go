package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	Info.Println("Initializing database")
	initializeDB()

	httpRouter := mux.NewRouter().StrictSlash(true)

	httpRouter.HandleFunc("/api/addemployee", apiAddEmployee).Methods("POST")
	httpRouter.HandleFunc("/api/updateschedule", apiUpdateSchedule).Methods("POST")
	httpRouter.HandleFunc("/api/addattendance", apiAddAttendance).Methods("POST")

	httpRouter.HandleFunc("/api/removeemployee", apiRemoveEmployee).Methods("DELETE")

	httpRouter.HandleFunc("/api/getallschedule/{idNumber}", apiGetAllYearsSchedule).Methods("GET")
	httpRouter.HandleFunc("/api/getschedule/{idNumber}/{schoolYear}", apiGetSchedule).Methods("GET")
	httpRouter.HandleFunc("/api/getemployee/{idNumber}", apiGetEmployee).Methods("GET")
	httpRouter.HandleFunc("/api/getattendance", apiGetAttendance).Methods("GET")

	httpCors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler(httpRouter)

	Error.Fatalln(http.ListenAndServe(":8080", httpCors))
}
