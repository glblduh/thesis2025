package main

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func apiAddEmployee(w http.ResponseWriter, r *http.Request) {
	body := apiAddEmployeeBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	addErr := addEmployee(strconv.Itoa(body.IdNumber), body.isFaculty, employee{
		FirstName:  body.FirstName,
		MiddleName: body.MiddleName,
		LastName:   body.LastName,
	})

	if addErr != nil {
		errorRes(w, addErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiRemoveEmployee(w http.ResponseWriter, r *http.Request) {
	body := apiRemoveEmployeeBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	removeErr := removeEmployee(strconv.Itoa(body.IdNumber), body.isFaculty)
	if removeErr != nil {
		errorRes(w, removeErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiUpdateSchedule(w http.ResponseWriter, r *http.Request) {
	body := apiUpdateScheduleBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	updateErr := updateEmployeeSchedule(strconv.Itoa(body.IdNumber), body.isFaculty, body.Schedule)
	if updateErr != nil {
		errorRes(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiGetAllYearsSchedule(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	employeeType, httpVarUnescapeErr := url.QueryUnescape(httpVars["employeeType"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	isFaculty, isFacultyCheckErr := checkIfFaculty(employeeType)
	if isFacultyCheckErr != nil {
		errorRes(w, isFacultyCheckErr.Error(), http.StatusBadRequest)
		return
	}

	employeeSchedules, getScheduleErr := getEmployeeAllYearsSchedule(idNumber, isFaculty)
	if getScheduleErr != nil {
		errorRes(w, getScheduleErr.Error(), http.StatusInternalServerError)
		return
	}

	idNumberInt, idNumberConvertErr := strconv.Atoi(idNumber)
	if idNumberConvertErr != nil {
		errorRes(w, idNumberConvertErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, apiGetAllYearsScheduleRes{
		IdNumber:  idNumberInt,
		isFaculty: isFaculty,
		Schedules: employeeSchedules,
	})
}

func apiGetSchedule(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	employeeType, httpVarUnescapeErr := url.QueryUnescape(httpVars["employeeType"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	schoolYearRequest, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	isFaculty, isFacultyCheckErr := checkIfFaculty(employeeType)
	if isFacultyCheckErr != nil {
		errorRes(w, isFacultyCheckErr.Error(), http.StatusBadRequest)
		return
	}

	employeeSchedule, getScheduleErr := getEmployeeSchedule(idNumber, isFaculty, schoolYearRequest)
	if getScheduleErr != nil {
		errorRes(w, getScheduleErr.Error(), http.StatusInternalServerError)
		return
	}

	idNumberInt, idNumberConvertErr := strconv.Atoi(idNumber)
	if idNumberConvertErr != nil {
		errorRes(w, idNumberConvertErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, apiGetScheduleRes{
		IdNumber:  idNumberInt,
		isFaculty: isFaculty,
		Schedule:  employeeSchedule,
	})
}
