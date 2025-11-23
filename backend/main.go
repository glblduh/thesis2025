package main

import "github.com/gorilla/mux"

func main() {
	Info.Println("Initializing database")
	initializeDB()

	httpRouter := mux.NewRouter().StrictSlash(true)

	httpRouter.HandleFunc("/api/addemployee", apiAddEmployee).Methods("POST")
	httpRouter.HandleFunc("/api/updateschedule", apiUpdateSchedule).Methods("POST")

	httpRouter.HandleFunc("/api/removeemployee", apiRemoveEmployee).Methods("DELETE")

	httpRouter.HandleFunc("/api/getallschedule/{employeeType}/{idNumber}", apiGetAllYearsSchedule).Methods("GET")
	httpRouter.HandleFunc("/api/getschedule/{employeeType}/{idNumber}/{schoolYear}", apiGetSchedule).Methods("GET")
}
