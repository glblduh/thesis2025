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

func splitSchoolYear(schoolYearString string) (schoolYearRange, error) {
	schoolYearStruct := schoolYearRange{}
	schoolYear := strings.Split(schoolYearString, "-")
	schoolYearStart, schoolYearConvertErr := strconv.Atoi(schoolYear[0])
	if schoolYearConvertErr != nil {
		return schoolYearStruct, schoolYearConvertErr
	}
	schoolYearEnd, schoolYearConvertErr := strconv.Atoi(schoolYear[1])
	if schoolYearConvertErr != nil {
		return schoolYearStruct, schoolYearConvertErr
	}
	schoolYearStruct.StartYear = schoolYearStart
	schoolYearStruct.EndYear = schoolYearEnd
	return schoolYearStruct, nil
}

func createSchoolYearStruct(schoolYear string) (schoolYearRange, error) {
	schoolYearStruct, splitSchoolYearErr := splitSchoolYear(schoolYear)
	if splitSchoolYearErr != nil {
		return schoolYearStruct, splitSchoolYearErr
	}
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

func getLatestSchoolYear(scheduleBucket *bbolt.Bucket) (schoolYearRange, string, error) {
	date := time.Now()
	currentYearString := strconv.Itoa(date.Year())
	lastestSchoolYearStruct := schoolYearRange{}
	lastestSchoolYear := ""

	lastSchoolYearSum := 0
	scheduleBucketCursor := scheduleBucket.Cursor()
	for schoolYear, _ := scheduleBucketCursor.First(); schoolYear != nil; schoolYear, _ = scheduleBucketCursor.Next() {
		schoolYearString := string(schoolYear)
		if !strings.Contains(schoolYearString, currentYearString) {
			continue
		}

		splittedSchoolYear, schoolYearSplitErr := splitSchoolYear(schoolYearString)
		if schoolYearSplitErr != nil {
			return lastestSchoolYearStruct, lastestSchoolYear, schoolYearSplitErr
		}

		schoolYearSum := splittedSchoolYear.StartYear + splittedSchoolYear.EndYear
		if schoolYearSum > lastSchoolYearSum {
			lastSchoolYearSum = schoolYearSum
			lastestSchoolYearStruct = splittedSchoolYear
			lastestSchoolYear = schoolYearString
		}
	}

	return lastestSchoolYearStruct, lastestSchoolYear, nil
}

func getDaySchedule(schoolYear schoolYearRange, date dayDate, scheduleBucket *bbolt.Bucket) (dayTimeRange, error) {
	daySchedule := dayTimeRange{}

	parsedDate := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC)
	schoolYearString := createSchoolYearString(schoolYear.StartYear, schoolYear.EndYear)

	scheduleBucketCursor := scheduleBucket.Cursor()
	for schoolYear, _ := scheduleBucketCursor.First(); schoolYear != nil; schoolYear, _ = scheduleBucketCursor.Next() {
		if schoolYearString != string(schoolYear) {
			continue
		}

		dayBucket := scheduleBucket.Bucket(schoolYear)
		dayByteSchedule := dayBucket.Get([]byte(parsedDate.Weekday().String()))

		unmarshalErr := json.Unmarshal(dayByteSchedule, &daySchedule)
		if unmarshalErr != nil {
			return daySchedule, unmarshalErr
		}
		break
	}

	return daySchedule, nil
}

func createAttendanceStruct(attendanceDayBucket *bbolt.Bucket, date dayDate, daySchedule dayTimeRange, suspended SuspensionType) (attendance, error) {
	attendanceStruct := attendance{}
	attendanceStruct.Date = date
	attendanceStruct.Suspended = suspended

	if daySchedule.DayOff {
		attendanceStruct.State = DAYOFF
		return attendanceStruct, nil
	}

	if attendanceDayBucket != nil && attendanceDayBucket.Get([]byte("LEAVE")) != nil {
		attendanceStruct.State = LEAVE
		return attendanceStruct, nil
	}

	if attendanceDayBucket == nil && !daySchedule.DayOff {
		attendanceStruct.State = ABSENT
		return attendanceStruct, nil
	}

	timeInByte := attendanceDayBucket.Get([]byte("TIMEIN"))
	timeOutByte := attendanceDayBucket.Get([]byte("TIMEOUT"))

	if timeInByte != nil {
		unmarshalErr := json.Unmarshal(timeInByte, &attendanceStruct.TimeIn)
		if unmarshalErr != nil {
			return attendanceStruct, unmarshalErr
		}
	}

	if timeInByte != nil && timeOutByte == nil {
		attendanceStruct.State = NOOUT
		return attendanceStruct, nil
	}

	if timeOutByte != nil {
		unmarshalErr := json.Unmarshal(timeOutByte, &attendanceStruct.TimeOut)
		if unmarshalErr != nil {
			return attendanceStruct, unmarshalErr
		}
	}

	attendanceStruct.State = AttendanceState(ATTENDED)
	return attendanceStruct, nil
}

func checkIfBucketEmpty(bucket *bbolt.Bucket) bool {
	bucketCursor := bucket.Cursor()
	first, _ := bucketCursor.First()
	return first == nil
}

func getDateSuspension(suspendedBucket *bbolt.Bucket, date dayDate) SuspensionType {
	suspensionType := SuspensionType(NOTSUSPENDED)

	yearBucket := suspendedBucket.Bucket([]byte(strconv.Itoa(date.Year)))

	var monthBucket *bbolt.Bucket
	if yearBucket != nil {
		monthBucket = yearBucket.Bucket([]byte(strconv.Itoa(date.Month)))
	}

	var dayBucket *bbolt.Bucket
	if monthBucket != nil {
		dayBucket = monthBucket.Bucket([]byte(strconv.Itoa(date.Day)))
	}

	if dayBucket != nil {
		suspensionType = SuspensionType((dayBucket.Get([]byte("TYPE"))))
	}

	return suspensionType
}
