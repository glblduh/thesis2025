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

	schoolYearStruct, createSchoolYearStructErr := createSchoolYearStruct(schoolYearRequest)
	if createSchoolYearStructErr != nil {
		errorRes(w, createSchoolYearStructErr.Error(), http.StatusInternalServerError)
		return
	}

	employeeSchedule, getScheduleErr := getEmployeeSchedule(idNumber, schoolYearStruct)
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

func apiAddAttendance(w http.ResponseWriter, r *http.Request) {
	body := apiAddAttendanceBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	addAttendanceErr := addAttendance(strconv.Itoa(body.IdNumber), body.IsLeave, body.LeaveReason, body.AttendanceTime)
	if addAttendanceErr != nil {
		errorRes(w, addAttendanceErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiGetAttendance(w http.ResponseWriter, r *http.Request) {
	body := apiGetAttendanceBody{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	employeeAttendance, getAttendanceErr := getAttendance(strconv.Itoa(body.IdNumber), body.SchoolYear, body.Date)
	if getAttendanceErr != nil {
		errorRes(w, getAttendanceErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, apiGetAttendanceRes{
		IdNumber: body.IdNumber,
		State:    employeeAttendance.State,
		Reason:   employeeAttendance.Reason,
		TimeIn:   employeeAttendance.TimeIn,
		TimeOut:  employeeAttendance.TimeOut,
	})
}

func apiGetMonthAttendances(w http.ResponseWriter, r *http.Request) {
	body := apiGetAttendanceBody{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	employeeMonthAttendances, getMonthAttendancesErr := getMonthAttendances(strconv.Itoa(body.IdNumber), body.SchoolYear, body.Date)
	if getMonthAttendancesErr != nil {
		errorRes(w, getMonthAttendancesErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, apiGetMonthAttendancesRes{
		IdNumber:    body.IdNumber,
		Attendances: employeeMonthAttendances,
	})
}

func apiGetAllEmployees(w http.ResponseWriter, r *http.Request) {
	allEmployeesStruct, getAllEmployeesErr := getAllEmployees()
	if getAllEmployeesErr != nil {
		errorRes(w, getAllEmployeesErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, allEmployeesStruct)
}

func apiRemoveSchedule(w http.ResponseWriter, r *http.Request) {
	body := apiRemoveScheduleBodyRes{};
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	removeScheduleErr := removeSchedule(strconv.Itoa(body.IdNumber), body.SchoolYear)
	if removeScheduleErr != nil {
		errorRes(w, removeScheduleErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}
