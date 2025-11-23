package main

import (
	"net/http"
	"strconv"
)

func apiAddEmployee(w http.ResponseWriter, r *http.Request) {
	body := apiAddEmployeeBody{}
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

	encodeRes(w, &apiAddEmployeeRes{
		IdNumber:  body.IdNumber,
		isFaculty: body.isFaculty,
	})
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
	body := apiUpdateScheduleBody{}
	if decodeBody(w, r.Body, &body) != nil {
		return
	}

	updateErr := updateEmployeeSchedule(strconv.Itoa(body.IdNumber), body.isFaculty, body.Schedule)
	if updateErr != nil {
		errorRes(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}
}
