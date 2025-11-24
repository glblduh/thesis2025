package main

import "github.com/gorilla/mux"

func main() {
	Info.Println("Initializing database")
	initializeDB()

	httpRouter := mux.NewRouter().StrictSlash(true)

	httpRouter.HandleFunc("/api/addemployee", apiAddEmployee).Methods("POST")
	httpRouter.HandleFunc("/api/updateschedule", apiUpdateSchedule).Methods("POST")

	httpRouter.HandleFunc("/api/removeemployee", apiRemoveEmployee).Methods("DELETE")

	httpRouter.HandleFunc("/api/getallschedule/{idNumber}", apiGetAllYearsSchedule).Methods("GET")
	httpRouter.HandleFunc("/api/getschedule/{idNumber}/{schoolYear}", apiGetSchedule).Methods("GET")
	httpRouter.HandleFunc("/api/getemployee/{idNumber}", apiGetEmployee).Methods("GET")
}
