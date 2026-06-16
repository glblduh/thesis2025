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

func apiUpdateAttendance(w http.ResponseWriter, r *http.Request) {
	body := apiUpdateAttendanceBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	addAttendanceErr := updateAttendance(strconv.Itoa(body.IdNumber), body.Attendance)
	if addAttendanceErr != nil {
		errorRes(w, addAttendanceErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiGetAttendance(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	idNumberInt, convertErr := strconv.Atoi(idNumber)
	if convertErr != nil {
		errorRes(w, convertErr.Error(), http.StatusInternalServerError)
		return
	}

	schoolYear, httpVarUnescapeErr := url.QueryUnescape(httpVars["schoolYear"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	yearString, httpVarUnescapeErr := url.QueryUnescape(httpVars["year"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	yearInt, convertErr := strconv.Atoi(yearString)
	if convertErr != nil {
		errorRes(w, convertErr.Error(), http.StatusInternalServerError)
		return
	}

	monthString, httpVarUnescapeErr := url.QueryUnescape(httpVars["month"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	monthInt, convertErr := strconv.Atoi(monthString)
	if convertErr != nil {
		errorRes(w, convertErr.Error(), http.StatusInternalServerError)
		return
	}

	dayString, httpVarUnescapeErr := url.QueryUnescape(httpVars["day"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	dayInt, convertErr := strconv.Atoi(dayString)
	if convertErr != nil {
		errorRes(w, convertErr.Error(), http.StatusInternalServerError)
		return
	}

	employeeAttendance, getAttendanceErr := getAttendance(idNumber, schoolYear, dayDate{
		Year: yearInt,
		Month: monthInt,
		Day: dayInt,
	})
	if getAttendanceErr != nil {
		errorRes(w, getAttendanceErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, apiGetAttendanceRes{
		IdNumber: idNumberInt,
		State:    employeeAttendance.State,
		TimeIn:   employeeAttendance.TimeIn,
		TimeOut:  employeeAttendance.TimeOut,
	})
}

func apiGetMonthAttendances(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	schoolYear, httpVarUnescapeErr := url.QueryUnescape(httpVars["schoolYear"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	yearString, httpVarUnescapeErr := url.QueryUnescape(httpVars["year"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	yearInt, convertErr := strconv.Atoi(yearString)
	if convertErr != nil {
		errorRes(w, convertErr.Error(), http.StatusInternalServerError)
		return
	}

	monthString, httpVarUnescapeErr := url.QueryUnescape(httpVars["month"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	monthInt, convertErr := strconv.Atoi(monthString)
	if convertErr != nil {
		errorRes(w, convertErr.Error(), http.StatusInternalServerError)
		return
	}

	employeeMonthAttendances, getMonthAttendancesErr := getMonthAttendances(idNumber, schoolYear, dayDate{
		Year: yearInt,
		Month: monthInt,
	})
	if getMonthAttendancesErr != nil {
		errorRes(w, getMonthAttendancesErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, apiGetMonthAttendancesRes{
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
	body := apiRemoveScheduleBodyRes{}
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

func apiGetAttendancesDates(w http.ResponseWriter, r*http.Request) {
	httpVars := mux.Vars(r)

	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}

	date := dayDate{Year: 0, Month: 0, Day: 0}
	yearUnescaped, isOk := httpVars["year"]
	if isOk {
		yearString, httpVarUnescapeErr := url.QueryUnescape(yearUnescaped)
		if httpVarUnescapeErr != nil {
			errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
			return
		}
		yearInt, convertErr := strconv.Atoi(yearString)
		if convertErr != nil {
			errorRes(w, convertErr.Error(), http.StatusInternalServerError)
			return
		}
		date.Year = yearInt
	}

	dates, getDatesErr := getAttendancesDates(idNumber, date)
	if getDatesErr != nil {
		errorRes(w, getDatesErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, dates)
}

func apiRemoveAttendance(w http.ResponseWriter, r *http.Request) {
	body := apiRemoveAttendanceBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	removeErr := removeAttendance(strconv.Itoa(body.IdNumber), body.Date)
	if removeErr != nil {
		errorRes(w, removeErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiAttend(w http.ResponseWriter, r *http.Request) {
	httpVars := mux.Vars(r)

	idNumber, httpVarUnescapeErr := url.QueryUnescape(httpVars["idNumber"])
	if httpVarUnescapeErr != nil {
		errorRes(w, httpVarUnescapeErr.Error(), http.StatusInternalServerError)
		return
	}
	idNumberInt, convertErr := strconv.Atoi(idNumber)
	if convertErr != nil {
		errorRes(w, convertErr.Error(), http.StatusInternalServerError)
		return
	}

	attendState, attendTime, checkAttendErr := checkAndAttend(idNumber)
	if checkAttendErr != nil {
		errorRes(w, checkAttendErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, apiAttendRes{
		IdNumber: idNumberInt,
		State: attendState,
		Time: attendTime,
	})
}

func apiUpdateSuspended(w http.ResponseWriter, r *http.Request) {
	body := apiAddSuspendedBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	updateSuspendedErr := updateSuspended(body.Date, body.Type)
	if updateSuspendedErr != nil {
		errorRes(w, updateSuspendedErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiRemoveSuspended(w http.ResponseWriter, r *http.Request) {
	body := apiRemoveSuspendedBodyRes{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	removeSuspendedErr := removeSuspended(body.Date)
	if removeSuspendedErr != nil {
		errorRes(w, removeSuspendedErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, body)
}

func apiGetAllSuspended(w http.ResponseWriter, r *http.Request) {
	allSuspended, getAllSuspendedErr := getAllSuspended()
	if getAllSuspendedErr != nil {
		errorRes(w, getAllSuspendedErr.Error(), http.StatusInternalServerError)
		return
	}

	encodeRes(w, allSuspended)
}
