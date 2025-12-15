package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	if !daySchedule.Change {
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
		IdNumber:     idNumber,
		IsFaculty:    isFaculty,
		EmployeeType: string(employeeBucket.Get([]byte("EmployeeType"))),
		FirstName:    string(employeeBucket.Get([]byte("FirstName"))),
		MiddleName:   string(employeeBucket.Get([]byte("MiddleName"))),
		LastName:     string(employeeBucket.Get([]byte("LastName"))),
	}
}

func dayScheduleDBToStruct(dayKey []byte, dayValue []byte, currentYearSchedule *employeeSchedule) error {
	currentDaySchedule := dayTimeRange{}

	dayValueUnmarshalErr := json.Unmarshal(dayValue, &currentDaySchedule)
	if dayValueUnmarshalErr != nil {
		return dayValueUnmarshalErr
	}

	switch string(dayKey) {
	case "Monday":
		currentYearSchedule.Monday = currentDaySchedule
	case "Tuesday":
		currentYearSchedule.Tuesday = currentDaySchedule
	case "Wednesday":
		currentYearSchedule.Wednesday = currentDaySchedule
	case "Thursday":
		currentYearSchedule.Thursday = currentDaySchedule
	case "Friday":
		currentYearSchedule.Friday = currentDaySchedule
	case "Saturday":
		currentYearSchedule.Saturday = currentDaySchedule
	case "Sunday":
		currentYearSchedule.Sunday = currentDaySchedule
	}

	return nil
}

func checkIfFaculty(employeeType string) (bool, error) {
	var isFaculty bool
	switch strings.ToLower(employeeType) {
	case "staff":
		isFaculty = false
	case "faculty":
		isFaculty = true
	default:
		return false, errors.New("not a valid employee type")
	}
	return isFaculty, nil
}

func createSchoolYearString(startYear int, endYear int) string {
	return strconv.Itoa(startYear) + "-" + strconv.Itoa(endYear)
}

func splitSchoolYear(schoolYearString string) (int, int, error) {
	schoolYear := strings.Split(schoolYearString, "-")
	schoolYearStart, schoolYearConvertErr := strconv.Atoi(schoolYear[0])
	if schoolYearConvertErr != nil {
		return 0, 0, schoolYearConvertErr
	}
	schoolYearEnd, schoolYearConvertErr := strconv.Atoi(schoolYear[1])
	if schoolYearConvertErr != nil {
		return 0, 0, schoolYearConvertErr
	}
	return schoolYearStart, schoolYearEnd, nil
}

func createSchoolYearStruct(schoolYear string) (schoolYearRange, error) {
	schoolYearStruct := schoolYearRange{}
	startYear, endYear, splitSchoolYearErr := splitSchoolYear(schoolYear)
	if splitSchoolYearErr != nil {
		return schoolYearStruct, splitSchoolYearErr
	}
	schoolYearStruct.StartYear = startYear
	schoolYearStruct.EndYear = endYear
	return schoolYearStruct, nil
}

func checkIfWorkingDay(scheduleBucket *bbolt.Bucket, schoolYear schoolYearRange, date dayDate) (bool, error) {
	schoolYearString := createSchoolYearString(schoolYear.StartYear, schoolYear.EndYear)
	schoolYearBucket := scheduleBucket.Bucket([]byte(schoolYearString))
	if schoolYearBucket == nil {
		return false, errors.New("school year not found")
	}
	dayName := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC).Weekday().String()
	dayBucket := scheduleBucket.Bucket([]byte(dayName))
	if dayBucket == nil {
		return false, nil
	}
	return true, nil
}

func getAttendanceState(dayBucket *bbolt.Bucket, scheduleBucket *bbolt.Bucket, schoolYear schoolYearRange, date dayDate) (attendance, error) {
	attendanceStruct := attendance{}

	isWorkingDay, checkWorkingDayErr := checkIfWorkingDay(scheduleBucket, schoolYear, date)
	if checkWorkingDayErr != nil {
		return attendanceStruct, checkWorkingDayErr
	}

	leaveReason := dayBucket.Get([]byte("LEAVE"))
	if leaveReason != nil {
		attendanceStruct.State = "LEAVE"
		attendanceStruct.Reason = string(leaveReason)
		return attendanceStruct, nil
	}

	if !isWorkingDay {
		attendanceStruct.State = "DAYOFF"
		return attendanceStruct, nil
	}

	if isWorkingDay && dayBucket == nil {
		attendanceStruct.State = "ABSENT"
		return attendanceStruct, nil
	}

	timeInUnmarshalErr := json.Unmarshal(dayBucket.Get([]byte("TIMEIN")), &attendanceStruct.TimeIn)
	if timeInUnmarshalErr != nil {
		return attendanceStruct, timeInUnmarshalErr
	}

	timeOutUnmarshalErr := json.Unmarshal(dayBucket.Get([]byte("TIMEOUT")), &attendanceStruct.TimeOut)
	if timeOutUnmarshalErr != nil {
		return attendanceStruct, timeOutUnmarshalErr
	}

	attendanceStruct.State = "ATTENDED"

	return attendanceStruct, nil
}
