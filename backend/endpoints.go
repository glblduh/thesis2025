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

	addErr := addEmployee(strconv.Itoa(body.IdNumber), body.IsFaculty, employee{
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

	removeErr := removeEmployee(strconv.Itoa(body.IdNumber))
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

	updateErr := updateEmployeeSchedule(strconv.Itoa(body.IdNumber), body.Schedule)
	if updateErr != nil {
		errorRes(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiGetAllYearsSchedule(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	employeeSchedules, getScheduleErr := getEmployeeAllYearsSchedule(idNumber)
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
		Schedules: employeeSchedules,
	})
}

func apiGetSchedule(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	schoolYearRequest, httpVarUnescapeErr := url.QueryUnescape(httpVars["schoolYear"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	employeeSchedule, getScheduleErr := getEmployeeSchedule(idNumber, schoolYearRequest)
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
		IdNumber: idNumberInt,
		Schedule: employeeSchedule,
	})
}

func apiGetEmployee(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	employeeStruct, employeeGetError := getEmployee(idNumber)
	if employeeGetError != nil {
		errorRes(w, employeeGetError.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, employeeStruct)
}
