package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//go:embed web/dist
var svelteFiles embed.FS

func main() {
	Info.Println("Initializing database")
	initializeDB()

	httpRouter := mux.NewRouter().StrictSlash(true)

	apiRouter := httpRouter.PathPrefix("/api").Subrouter()
	apiPostRouter := apiRouter.Methods(http.MethodPost).Subrouter()
	apiDeleteRouter := apiRouter.Methods(http.MethodDelete).Subrouter()
	apiGetRouter := apiRouter.Methods(http.MethodGet).Subrouter()

	apiPostRouter.HandleFunc("/addemployee", apiAddEmployee)
	apiPostRouter.HandleFunc("/updateschedule", apiUpdateSchedule)
	apiPostRouter.HandleFunc("/addattendance", apiAddAttendance)

	apiDeleteRouter.HandleFunc("/removeemployee", apiRemoveEmployee)

	apiGetRouter.HandleFunc("/getallschedule/{idNumber}", apiGetAllYearsSchedule)
	apiGetRouter.HandleFunc("/getschedule/{idNumber}/{schoolYear}", apiGetSchedule)
	apiGetRouter.HandleFunc("/getemployee/{idNumber}", apiGetEmployee)
	apiGetRouter.HandleFunc("/getattendance", apiGetAttendance)
	apiGetRouter.HandleFunc("/getmonthattendances", apiGetMonthAttendances)
	apiGetRouter.HandleFunc("/getallemployees", apiGetAllEmployees)

	svelteFS, fsErr := fs.Sub(svelteFiles, "web/dist")
	if fsErr != nil {
		Error.Fatalln("Cannot get embedded files")
	}

	httpRouter.PathPrefix("/").Handler(http.FileServer(http.FS(svelteFS)))

	httpCors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler(httpRouter)

	Error.Fatalln(http.ListenAndServe(":8080", httpCors))
}
