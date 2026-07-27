package main

import (
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func startHTTP() {
	httpRouter := mux.NewRouter().StrictSlash(true)

	authRouter := httpRouter.PathPrefix("/auth").Subrouter()
	authPostRouter := authRouter.Methods(http.MethodPost).Subrouter()
	authPostRouter.HandleFunc("/validate", apiUserAuth)

	apiRouter := httpRouter.PathPrefix("/api").Subrouter()
	apiRouter.Use(apiAuthMiddleware)
	apiPostRouter := apiRouter.Methods(http.MethodPost).Subrouter()
	apiDeleteRouter := apiRouter.Methods(http.MethodDelete).Subrouter()
	apiGetRouter := apiRouter.Methods(http.MethodGet).Subrouter()

	apiPostRouter.HandleFunc("/addemployee", apiAddEmployee)
	apiPostRouter.HandleFunc("/updateschedule", apiUpdateSchedule)
	apiPostRouter.HandleFunc("/updateattendance", apiUpdateAttendance)
	apiPostRouter.HandleFunc("/updateattendance", apiUpdateAttendance)
	apiPostRouter.HandleFunc("/attend/{idNumber}", apiAttend)
	apiPostRouter.HandleFunc("/updatesuspended", apiUpdateSuspended)

	apiDeleteRouter.HandleFunc("/removeemployee", apiRemoveEmployee)
	apiDeleteRouter.HandleFunc("/removeschedule", apiRemoveSchedule)
	apiDeleteRouter.HandleFunc("/removeattendance", apiRemoveAttendance)
	apiDeleteRouter.HandleFunc("/removesuspended", apiRemoveSuspended)

	apiGetRouter.HandleFunc("/getallschedule/{idNumber}", apiGetAllYearsSchedule)
	apiGetRouter.HandleFunc("/getschedule/{idNumber}/{schoolYear}", apiGetSchedule)
	apiGetRouter.HandleFunc("/getemployee/{idNumber}", apiGetEmployee)
	apiGetRouter.HandleFunc("/getattendance/{idNumber}/{schoolYear}/{year}/{month}/{day}", apiGetAttendance)
	apiGetRouter.HandleFunc("/getmonthattendances/{idNumber}/{schoolYear}/{year}/{month}", apiGetMonthAttendances)
	apiGetRouter.HandleFunc("/getallemployees", apiGetAllEmployees)
	apiGetRouter.HandleFunc("/getallattendancesyears/{idNumber}", apiGetAttendancesDates)
	apiGetRouter.HandleFunc("/getallattendancesmonths/{idNumber}/{year}", apiGetAttendancesDates)
	apiGetRouter.HandleFunc("/getallsuspended", apiGetAllSuspended)
	apiGetRouter.HandleFunc("/getallmonthattendances/{schoolYear}/{year}/{month}", apiGetAllMonthAttendances)
	apiGetRouter.HandleFunc("/getallschoolyears", apiGetAllSchoolYears)

	svelteFS, fsErr := fs.Sub(svelteFiles, "web/dist")
	if fsErr != nil {
		Error.Fatalln("Cannot get embedded files")
	}

	httpRouter.PathPrefix("/").Handler(http.FileServer(http.FS(svelteFS)))

	httpCors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler(httpRouter)

	httpPort := ":8080"
	Error.Fatalln(http.ListenAndServe(httpPort, httpCors))
	Info.Printf("HTTP server listening on port %s", httpPort)
}
