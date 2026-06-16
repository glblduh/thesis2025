package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"go.etcd.io/bbolt"
)

func openDB() (*bbolt.DB, error) {
	return bbolt.Open(".records.db", 0660, &bbolt.Options{
		Timeout: time.Second,
	})
}

func initializeDB() {
	db, dbErr := openDB()
	if dbErr != nil {
		Error.Fatalln("Opening DB error")
	}
	defer db.Close()

	db.Update(func(tx *bbolt.Tx) error {
		_, facultyCreateErr := tx.CreateBucketIfNotExists([]byte("Faculty"))
		if facultyCreateErr != nil {
			Error.Fatalln("Faculty bucket creation error")
		}
		_, staffCreateErr := tx.CreateBucketIfNotExists([]byte("Staff"))
		if staffCreateErr != nil {
			Error.Fatalln("Staff bucket creation error")
		}
		_, suspendedCreateErr := tx.CreateBucketIfNotExists([]byte("Suspended"))
		if suspendedCreateErr != nil {
			Error.Fatalln("Suspended bucket creation error")
		}
		return nil
	})
}

func addEmployee(idNumber string, isFaculty bool, employeeStruct employee) error {
	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	employeeType := "Staff"
	if isFaculty {
		employeeType = "Faculty"
	}

	return db.Update(func(tx *bbolt.Tx) error {
		bucket, bucketErr := tx.Bucket([]byte(employeeType)).CreateBucketIfNotExists([]byte(idNumber))
		if bucketErr != nil {
			return bucketErr
		}
		firstNamePutErr := bucket.Put([]byte("FirstName"), []byte(employeeStruct.FirstName))
		if firstNamePutErr != nil {
			return firstNamePutErr
		}
		middleNamePutErr := bucket.Put([]byte("MiddleName"), []byte(employeeStruct.MiddleName))
		if middleNamePutErr != nil {
			return middleNamePutErr
		}
		lastNamePutErr := bucket.Put([]byte("LastName"), []byte(employeeStruct.LastName))
		if lastNamePutErr != nil {
			return lastNamePutErr
		}
		employeeTypePutErr := bucket.Put([]byte("EmployeeType"), []byte(employeeType))
		if employeeTypePutErr != nil {
			return employeeTypePutErr
		}
		_, scheduleBucketErr := bucket.CreateBucketIfNotExists([]byte("Schedule"))
		if scheduleBucketErr != nil {
			return scheduleBucketErr
		}
		_, attendanceBucketErr := bucket.CreateBucketIfNotExists([]byte("Attendance"))
		if attendanceBucketErr != nil {
			return attendanceBucketErr
		}
		return nil
	})
}

func removeEmployee(idNumber string) error {
	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(employeeStruct.EmployeeType)).DeleteBucket([]byte(idNumber))
	})
}

func updateEmployeeSchedule(idNumber string, schedule employeeSchedule) error {
	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	schoolYear := createSchoolYearString(schedule.SchoolYear.StartYear, schedule.SchoolYear.EndYear)

	return db.Update(func(tx *bbolt.Tx) error {
		scheduleBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Schedule"))
		schoolYearBucket, schoolYearBucketErr := scheduleBucket.CreateBucketIfNotExists([]byte(schoolYear))
		if schoolYearBucketErr != nil {
			return schoolYearBucketErr
		}

		var scheduleAddErr error
		scheduleAddErr = writeDaySchedule(schoolYearBucket, "Monday", schedule.Monday)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		scheduleAddErr = writeDaySchedule(schoolYearBucket, "Tuesday", schedule.Tuesday)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		scheduleAddErr = writeDaySchedule(schoolYearBucket, "Wednesday", schedule.Wednesday)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		scheduleAddErr = writeDaySchedule(schoolYearBucket, "Thursday", schedule.Thursday)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		scheduleAddErr = writeDaySchedule(schoolYearBucket, "Friday", schedule.Friday)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		scheduleAddErr = writeDaySchedule(schoolYearBucket, "Saturday", schedule.Saturday)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		scheduleAddErr = writeDaySchedule(schoolYearBucket, "Sunday", schedule.Sunday)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		return nil
	})
}

func getEmployee(idNumber string) (employee, error) {
	employeeStruct := employee{}

	db, dbErr := openDB()
	if dbErr != nil {
		return employeeStruct, dbErr
	}
	defer db.Close()

	employeeCheckErr := db.View(func(tx *bbolt.Tx) error {
		staffBucket := tx.Bucket([]byte("Staff"))
		facultyBucket := tx.Bucket([]byte("Faculty"))

		staffCheck := staffBucket.Bucket([]byte(idNumber))
		facultyCheck := facultyBucket.Bucket([]byte(idNumber))

		if staffCheck == nil && facultyCheck == nil {
			return errors.New("id number not found")
		}

		idNumberInt, idNumberConvertErr := strconv.Atoi(idNumber)
		if idNumberConvertErr != nil {
			return idNumberConvertErr
		}

		if staffCheck != nil {
			employeeStruct = employeeInfoDBToStruct(staffBucket.Bucket([]byte(idNumber)), idNumberInt, false)
			return nil
		}

		employeeStruct = employeeInfoDBToStruct(facultyBucket.Bucket([]byte(idNumber)), idNumberInt, true)
		return nil
	})
	if employeeCheckErr != nil {
		return employeeStruct, employeeCheckErr
	}

	return employeeStruct, nil
}

func getEmployeeAllYearsSchedule(idNumber string) ([]employeeSchedule, error) {
	allYearsSchedule := []employeeSchedule{}

	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return allYearsSchedule, verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return nil, dbErr
	}
	defer db.Close()

	scheduleIterateErr := db.View(func(tx *bbolt.Tx) error {
		employeeScheduleBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Schedule"))

		cursor := employeeScheduleBucket.Cursor()

		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			currentYearSchedule := employeeSchedule{}
			splittedSchoolYearStruct, schoolYearSplitErr := splitSchoolYear(string(key))
			if schoolYearSplitErr != nil {
				return schoolYearSplitErr
			}
			currentYearSchedule.SchoolYear = splittedSchoolYearStruct

			currentKeyCursor := employeeScheduleBucket.Bucket(key).Cursor()
			for dayKey, dayValue := currentKeyCursor.First(); dayKey != nil; dayKey, dayValue = currentKeyCursor.Next() {
				scheduleConvertGetErr := dayScheduleDBToStruct(dayKey, dayValue, &currentYearSchedule)
				if scheduleConvertGetErr != nil {
					return scheduleConvertGetErr
				}
			}
			allYearsSchedule = append(allYearsSchedule, currentYearSchedule)
		}
		return nil
	})
	if scheduleIterateErr != nil {
		return nil, scheduleIterateErr
	}

	return allYearsSchedule, nil
}

func getEmployeeSchedule(idNumber string, schoolYear schoolYearRange) (employeeSchedule, error) {
	employeeSchedule := employeeSchedule{}

	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return employeeSchedule, verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return employeeSchedule, dbErr
	}
	defer db.Close()

	scheduleIterateErr := db.View(func(tx *bbolt.Tx) error {
		schoolYearString := createSchoolYearString(schoolYear.StartYear, schoolYear.EndYear)

		schoolYearBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte(schoolYearString))
		if schoolYearBucket == nil {
			return errors.New("school year not found")
		}

		employeeSchedule.SchoolYear = schoolYearRange{
			StartYear: schoolYear.StartYear,
			EndYear:   schoolYear.EndYear,
		}

		schoolYearCursor := schoolYearBucket.Cursor()
		for dayKey, dayValue := schoolYearCursor.First(); dayKey != nil; dayKey, dayValue = schoolYearCursor.Next() {
			scheduleConvertGetErr := dayScheduleDBToStruct(dayKey, dayValue, &employeeSchedule)
			if scheduleConvertGetErr != nil {
				return scheduleConvertGetErr
			}
		}
		return nil
	})
	if scheduleIterateErr != nil {
		return employeeSchedule, scheduleIterateErr
	}

	return employeeSchedule, nil
}

func updateAttendance(idNumber string, attendanceStruct attendance) error {
	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	month := attendanceStruct.Date.Month
	if month < 1 || month > 12 {
		return errors.New("invalid month")
	}

	day := attendanceStruct.Date.Day
	if day < 1 || day > 31 {
		return errors.New("invalid day")
	}

	year := attendanceStruct.Date.Year
	parseDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if parseDate.Year() != year || parseDate.Month() != time.Month(month) || parseDate.Day() != day {
		return errors.New("invalid date")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		attendanceBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Attendance"))

		yearBucket, yearBucketErr := attendanceBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(year)))
		if yearBucketErr != nil {
			return yearBucketErr
		}

		monthBucket, monthBucketErr := yearBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(month)))
		if monthBucketErr != nil {
			return monthBucketErr
		}

		dayBucket, dayBucketErr := monthBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(day)))
		if dayBucketErr != nil {
			return dayBucketErr
		}

		if attendanceStruct.State == LEAVE {
			statePutErr := dayBucket.Put([]byte("LEAVE"), []byte(""))
			if statePutErr != nil {
				return statePutErr
			}
			return nil
		}

		removeLeaveErr := dayBucket.Delete([]byte("LEAVE"))
		if removeLeaveErr != nil {
			return removeLeaveErr
		}

		timeInByte, timeInMarshalErr := json.Marshal(attendanceStruct.TimeIn)
		if timeInMarshalErr != nil {
			return timeInMarshalErr
		}
		timeInPutErr := dayBucket.Put([]byte("TIMEIN"), timeInByte)
		if timeInPutErr != nil {
			return timeInPutErr
		}

		timeOutByte, timeOutMarshalErr := json.Marshal(attendanceStruct.TimeOut)
		if timeOutMarshalErr != nil {
			return timeOutMarshalErr
		}
		timeOutPutErr := dayBucket.Put([]byte("TIMEOUT"), timeOutByte)
		if timeOutPutErr != nil {
			return timeOutPutErr
		}

		return nil
	})
}

func getAttendance(idNumber string, schoolYearString string, date dayDate) (attendance, error) {
	attendanceStruct := attendance{}

	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return attendanceStruct, verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return attendanceStruct, dbErr
	}
	defer db.Close()

	dbViewErr := db.View(func(tx *bbolt.Tx) error {
		attendanceBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Attendance"))
		suspendedBucket := tx.Bucket([]byte("Suspended"))

		schoolYearScheduleBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Schedule")).Bucket([]byte(schoolYearString))
		if schoolYearScheduleBucket == nil {
			return errors.New("school year not found")
		}

		yearBucket := attendanceBucket.Bucket([]byte(strconv.Itoa(date.Year)))
		if yearBucket == nil {
			return errors.New("year not found")
		}

		monthBucket := yearBucket.Bucket([]byte(strconv.Itoa(date.Month)))
		if monthBucket == nil {
			return errors.New("month not found")
		}

		dayBucket := monthBucket.Bucket([]byte(strconv.Itoa(date.Day)))
		if dayBucket == nil {
			return errors.New("day not found")
		}

		parsedDay := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC)
		dayScheduleByte := schoolYearScheduleBucket.Get([]byte(parsedDay.Weekday().String()))
		dayScheduleStruct := dayTimeRange{}

		unmarshalErr := json.Unmarshal(dayScheduleByte, &dayScheduleStruct)
		if unmarshalErr != nil {
			return unmarshalErr
		}

		var createAttendanceErr error
		attendanceStruct, createAttendanceErr = createAttendanceStruct(dayBucket, date, dayScheduleStruct, getDateSuspension(suspendedBucket, date))
		if createAttendanceErr != nil {
			return createAttendanceErr
		}

		return nil
	})
	if dbViewErr != nil {
		return attendanceStruct, dbViewErr
	}

	return attendanceStruct, nil
}

func getMonthAttendances(idNumber string, schoolYearString string, date dayDate) ([]attendance, error) {
	monthAttendances := []attendance{}

	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return monthAttendances, verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return monthAttendances, dbErr
	}
	defer db.Close()

	dbViewErr := db.View(func(tx *bbolt.Tx) error {
		attendanceBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Attendance"))
		suspendedBucket := tx.Bucket([]byte("Suspended"))

		schoolYearScheduleBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Schedule")).Bucket([]byte(schoolYearString))
		if schoolYearScheduleBucket == nil {
			return errors.New("school year not found")
		}

		yearBucket := attendanceBucket.Bucket([]byte(strconv.Itoa(date.Year)))
		if yearBucket == nil {
			return errors.New("year not found")
		}

		monthBucket := yearBucket.Bucket([]byte(strconv.Itoa(date.Month)))
		if monthBucket == nil {
			return errors.New("month not found")
		}

		parsedDate := time.Date(date.Year, time.Month(date.Month)+1, 0, 0, 0, 0, 0, time.UTC)
		for i := 1; i <= parsedDate.Day(); i++ {
			dayBucket := monthBucket.Bucket([]byte(strconv.Itoa(i)))
			parsedDay := time.Date(date.Year, time.Month(date.Month), i, 0, 0, 0, 0, time.UTC)
			dayScheduleByte := schoolYearScheduleBucket.Get([]byte(parsedDay.Weekday().String()))
			dayScheduleStruct := dayTimeRange{}

			unmarshalErr := json.Unmarshal(dayScheduleByte, &dayScheduleStruct)
			if unmarshalErr != nil {
				return unmarshalErr
			}

			iterationDate := dayDate{
				Year: date.Year,
				Month: date.Month,
				Day: i,
			}

			dayAttendance, createAttendanceErr := createAttendanceStruct(dayBucket, iterationDate, dayScheduleStruct, getDateSuspension(suspendedBucket, iterationDate))
			if createAttendanceErr != nil {
				return createAttendanceErr
			}
			monthAttendances = append(monthAttendances, dayAttendance)
		}

		return nil
	})
	if dbViewErr != nil {
		return monthAttendances, dbViewErr
	}

	return monthAttendances, nil
}

func getAllEmployees() (allEmployees, error) {
	allEmployeesStruct := allEmployees{}

	db, dbErr := openDB()
	if dbErr != nil {
		return allEmployeesStruct, dbErr
	}
	defer db.Close()

	dbViewErr := db.View(func(tx *bbolt.Tx) error {
		staffBucket := tx.Bucket([]byte("Staff"))
		facultyBucket := tx.Bucket([]byte("Faculty"))

		staffBucketCursor := staffBucket.Cursor()
		for staffKey, _ := staffBucketCursor.First(); staffKey != nil; staffKey, _ = staffBucketCursor.Next() {
			currentStaffBucket := staffBucket.Bucket(staffKey)
			currentStaffIdNumber, convertErr := strconv.Atoi(string(staffKey))
			if convertErr != nil {
				return convertErr
			}

			allEmployeesStruct.Staff = append(allEmployeesStruct.Staff, employeeInfoDBToStruct(currentStaffBucket, currentStaffIdNumber, false))
		}

		facultyBucketCursor := facultyBucket.Cursor()
		for facultyKey, _ := facultyBucketCursor.First(); facultyKey != nil; facultyKey, _ = facultyBucketCursor.Next() {
			currentFacultyBucket := facultyBucket.Bucket(facultyKey)
			currentFacultyIdNumber, convertErr := strconv.Atoi(string(facultyKey))
			if convertErr != nil {
				return convertErr
			}

			allEmployeesStruct.Faculty = append(allEmployeesStruct.Faculty, employeeInfoDBToStruct(currentFacultyBucket, currentFacultyIdNumber, true))
		}

		return nil
	})
	if dbViewErr != nil {
		return allEmployeesStruct, dbViewErr
	}

	return allEmployeesStruct, nil
}

func removeSchedule(idNumber string, schoolYear string) error {
	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Schedule")).DeleteBucket([]byte(schoolYear))
	})
}

func getAttendancesDates(idNumber string, date dayDate) (attendanceDates, error) {
	dates := attendanceDates{}

	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return dates, verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return dates, dbErr
	}
	defer db.Close()

	dbViewErr := db.View(func(tx *bbolt.Tx) error {
		attendanceBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Attendance"))

		if date.Year != 0 {
			yearString := strconv.Itoa(date.Year)

			yearBucket := attendanceBucket.Bucket([]byte(yearString))
			if yearBucket == nil {
				return errors.New("year not found")
			}

			yearBucketCursor := yearBucket.Cursor()
			for month, _ := yearBucketCursor.First(); month != nil; month, _ = yearBucketCursor.Next() {
				monthString := string(month)
				monthInt, convertErr := strconv.Atoi(monthString)
				if convertErr != nil {
					return convertErr
				}
				dates.Months = append(dates.Months, monthInt)
			}
			return nil
		}

		if date.Year != 0 && date.Month != 0 {
			yearString := strconv.Itoa(date.Year)
			monthString := strconv.Itoa(date.Month)

			yearBucket := attendanceBucket.Bucket([]byte(yearString))
			if yearBucket == nil {
				return errors.New("year not found")
			}

			monthBucket := yearBucket.Bucket([]byte(monthString))
			if monthBucket != nil {
				return errors.New("month not found")
			}

			monthBucketCursor := monthBucket.Cursor()
			for day, _ := monthBucketCursor.First(); day != nil; day, _ = monthBucketCursor.Next() {
				dayString := string(day)
				dayInt, convertErr := strconv.Atoi(dayString)
				if convertErr != nil {
					return convertErr
				}
				dates.Days = append(dates.Days, dayInt)
			}
			return nil
		}

		attendanceBucketCursor := attendanceBucket.Cursor()
		for year, _ := attendanceBucketCursor.First(); year != nil; year, _ = attendanceBucketCursor.Next() {
			yearString := string(year)
			yearInt, convertErr := strconv.Atoi(yearString)
			if convertErr != nil {
				return convertErr
			}
			dates.Years = append(dates.Years, yearInt)
		}
		return nil
	})
	if dbViewErr != nil {
		return dates, dbViewErr
	}

	return dates, nil
}

func removeAttendance(idNumber string, date dayDate) error {
	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	month := date.Month
	if month < 1 || month > 12 {
		return errors.New("invalid month")
	}

	day := date.Day
	if day < 1 || day > 31 {
		return errors.New("invalid day")
	}

	year := date.Year
	parseDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if parseDate.Year() != year || parseDate.Month() != time.Month(month) || parseDate.Day() != day {
		return errors.New("invalid date")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		attendanceBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Attendance"))

		yearBucket := attendanceBucket.Bucket([]byte(strconv.Itoa(date.Year)))
		if yearBucket == nil {
			return errors.New("year does not exist")
		}

		monthBucket := yearBucket.Bucket([]byte(strconv.Itoa(date.Month)))
		if monthBucket == nil {
			return errors.New("month does not exist")
		}

		removeDayErr := monthBucket.DeleteBucket([]byte(strconv.Itoa(date.Day)))
		if removeDayErr != nil {
			return removeDayErr
		}

		if checkIfBucketEmpty(monthBucket) {
			removeMonthErr := yearBucket.DeleteBucket([]byte(strconv.Itoa(date.Month)))
			if removeMonthErr != nil {
				return removeMonthErr
			}
		}

		if checkIfBucketEmpty(yearBucket) {
			removeYearErr := attendanceBucket.DeleteBucket([]byte(strconv.Itoa(date.Year)))
			if removeYearErr != nil {
				return removeYearErr
			}
		}

		return nil
	})
}

func checkAndAttend(idNumber string) (AttendState, attendanceTime, error) {
	attendState := AttendState("");
	attendTime := attendanceTime{}

	employeeStruct, verifyErr := getEmployee(idNumber)
	if verifyErr != nil {
		return attendState, attendTime, verifyErr
	}

	db, dbErr := openDB()
	if dbErr != nil {
		return attendState, attendTime, dbErr
	}
	defer db.Close()

	dbUpdate := db.Update(func(tx *bbolt.Tx) error {
		attendanceBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Attendance"))
		scheduleBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Schedule"))

		_, currentSchoolYearString, getSchoolYearErr := getLatestSchoolYear(scheduleBucket)
		if getSchoolYearErr != nil {
			return getSchoolYearErr
		}
		if currentSchoolYearString == "" {
			return errors.New("current year does not belong to any registered school year")
		}

		currentTime := time.Now()
		currentYear := currentTime.Year()
		currentMonth := int(currentTime.Month())
		currentDay := currentTime.Day()

		currentDayScheduleStruct := dayTimeRange{}
		currentDayScheduleByte := scheduleBucket.Bucket([]byte(currentSchoolYearString)).Get([]byte(currentTime.Weekday().String()))
		unmarshalErr := json.Unmarshal(currentDayScheduleByte, &currentDayScheduleStruct)
		if unmarshalErr != nil {
			return unmarshalErr
		}

		if currentDayScheduleStruct.DayOff {
			return errors.New("day off")
		}

		yearBucket, yearBucketErr := attendanceBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(currentYear)))
		if yearBucketErr != nil {
			return yearBucketErr
		}

		monthBucket, monthBucketErr := yearBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(currentMonth)))
		if monthBucketErr != nil {
			return monthBucketErr
		}

		dayBucket, dayBucketErr := monthBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(currentDay)))
		if dayBucketErr != nil {
			return dayBucketErr
		}

		if dayBucket.Get([]byte("TIMEIN")) != nil && dayBucket.Get([]byte("TIMEOUT")) != nil {
			return errors.New("attendance for today is already completed")
		}

		currentTimeStruct := attendanceTime{
			Hour: currentTime.Hour(),
			Minute: currentTime.Minute(),
			Unix: int(currentTime.Unix()),
		}
		currentTimeByte, marshalErr := json.Marshal(currentTimeStruct)
		if marshalErr != nil {
			return marshalErr
		}

		if dayBucket.Get([]byte("TIMEIN")) != nil {
			timeOutPutErr := dayBucket.Put([]byte("TIMEOUT"), currentTimeByte)
			if timeOutPutErr != nil {
				return timeOutPutErr
			}
			attendState = TIMEOUT
			attendTime = currentTimeStruct
			return nil
		}

		timeInPutErr := dayBucket.Put([]byte("TIMEIN"), currentTimeByte)
		if timeInPutErr != nil {
			return timeInPutErr
		}

		attendState = TIMEIN
		attendTime = currentTimeStruct
		return nil
	})
	if dbUpdate != nil {
		return attendState, attendTime, dbUpdate
	}

	return attendState, attendTime, nil
}

func updateSuspended(date dayDate, suspensionType SuspensionType) error {
	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		suspendedBucket := tx.Bucket([]byte("Suspended"))
		if suspendedBucket == nil {
			return errors.New("suspended bucket not found")
		}

		yearBucket, yearBucketErr := suspendedBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(date.Year)))
		if yearBucketErr != nil {
			return yearBucketErr
		}

		monthBucket, monthBucketErr := yearBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(date.Month)))
		if monthBucketErr != nil {
			return monthBucketErr
		}

		dayBucket, dayBucketErr := monthBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(date.Day)))
		if dayBucketErr != nil {
			return dayBucketErr
		}

		dateByte, dateMarshalErr := json.Marshal(date)
		if dateMarshalErr != nil {
			return dateMarshalErr
		}
		dayBucket.Put([]byte("DATE"), dateByte)
		dayBucket.Put([]byte("TYPE"), []byte(suspensionType))

		return nil
	})
}

func removeSuspended(date dayDate) error {
	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		suspendedBucket := tx.Bucket([]byte("Suspended"))
		if suspendedBucket == nil {
			return errors.New("suspended bucket not found")
		}

		yearBucket := suspendedBucket.Bucket([]byte(strconv.Itoa(date.Year)))
		if yearBucket == nil {
			return errors.New("year bucket not found")
		}

		monthBucket := yearBucket.Bucket([]byte(strconv.Itoa(date.Month)))
		if monthBucket == nil {
			return errors.New("month bucket not found")
		}

		removeDayErr := monthBucket.DeleteBucket([]byte(strconv.Itoa(date.Day)))
		if removeDayErr != nil {
			return removeDayErr
		}

		if checkIfBucketEmpty(monthBucket) {
			removeMonthErr := yearBucket.DeleteBucket([]byte(strconv.Itoa(date.Month)))
			if removeMonthErr != nil {
				return removeMonthErr
			}
		}

		if checkIfBucketEmpty(yearBucket) {
			removeYearErr := suspendedBucket.DeleteBucket([]byte(strconv.Itoa(date.Year)))
			if removeYearErr != nil {
				return removeYearErr
			}
		}

		return nil
	})
}

func getAllSuspended() ([]suspendedDay, error) {
	allSuspensions := []suspendedDay{}

	db, dbErr := openDB()
	if dbErr != nil {
		return allSuspensions, dbErr
	}
	defer db.Close()

	return allSuspensions, db.View(func(tx *bbolt.Tx) error {
		suspendedBucket := tx.Bucket([]byte("Suspended"))
		if suspendedBucket == nil {
			return errors.New("suspended bucket not found")
		}

		suspendedBucketCursor := suspendedBucket.Cursor()
		for year, _ := suspendedBucketCursor.First(); year != nil; year, _ = suspendedBucketCursor.Next() {
			yearBucket := suspendedBucket.Bucket(year)
			yearBucketCursor := yearBucket.Cursor()

			for month, _ := yearBucketCursor.First(); month != nil; month, _ = yearBucketCursor.Next() {
				monthBucket := yearBucket.Bucket(month)
				monthBucketCursor := monthBucket.Cursor()

				for day, _ := monthBucketCursor.First(); day != nil; day, _ = monthBucketCursor.Next() {
					dayBucket := monthBucket.Bucket(day)

					date := dayDate{}
					dateUnmarshalErr := json.Unmarshal(dayBucket.Get([]byte("DATE")), &date)
					if dateUnmarshalErr != nil {
						return dateUnmarshalErr
					}

					allSuspensions = append(allSuspensions, suspendedDay{
						Date: date,
						Type: SuspensionType(string(dayBucket.Get([]byte("TYPE")))),
					})
				}
			}
		}

		return nil
	})
}
