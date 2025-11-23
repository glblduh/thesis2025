package main

import (
	"encoding/json"
	"io"
	"net/http"

	"go.etcd.io/bbolt"
)

func encodeRes(w http.ResponseWriter, v any) error {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		errorRes(w, "JSON Encoder error", http.StatusInternalServerError)
	}
	return err
}

func errorRes(w http.ResponseWriter, errorResponse string, code int) {
	w.WriteHeader(code)
	err := encodeRes(w, &jsonErrorRes{
		Error: errorResponse,
	})
	if err != nil {
		w.Write([]byte(errorResponse))
	}
}

func decodeBody(w http.ResponseWriter, body io.Reader, v any) error {
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		errorRes(w, "JSON Encoder error", http.StatusInternalServerError)
	}
	return err
}

func writeDaySchedule(schoolYearBucket *bbolt.Bucket, dayName string, daySchedule dayTimeRange) error {
	if daySchedule.dontChange {
		return nil
	}

	dayScheduleByte, scheduleMarshalErr := json.Marshal(daySchedule)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}

	daySchedulePutErr := schoolYearBucket.Put([]byte(dayName), dayScheduleByte)
	if daySchedulePutErr != nil {
		return scheduleMarshalErr
	}

	return nil
}

func employeeInfoDBToStruct(employeeBucket *bbolt.Bucket, idNumber int, isFaculty bool) employee {
	return employee{
		IdNumber:   idNumber,
		isFaculty:  isFaculty,
		FirstName:  string(employeeBucket.Get([]byte("FirstName"))),
		MiddleName: string(employeeBucket.Get([]byte("MiddleName"))),
		LastName:   string(employeeBucket.Get([]byte("Name"))),
	}
}
