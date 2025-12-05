package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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
		_, leaveBucketErr := bucket.CreateBucketIfNotExists([]byte("Leave"))
		if leaveBucketErr != nil {
			return leaveBucketErr
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

	schoolYear := strconv.Itoa(schedule.SchoolYear.StartYear) + "-" + strconv.Itoa(schedule.SchoolYear.EndYear)

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
			schoolYear := strings.Split(string(key), "-")
			schoolYearStart, schoolYearConvertErr := strconv.Atoi(schoolYear[0])
			if schoolYearConvertErr != nil {
				return schoolYearConvertErr
			}
			schoolYearEnd, schoolYearConvertErr := strconv.Atoi(schoolYear[1])
			if schoolYearConvertErr != nil {
				return schoolYearConvertErr
			}
			schoolYearStruct := schoolYearRange{
				StartYear: schoolYearStart,
				EndYear:   schoolYearEnd,
			}

			currentYearSchedule.SchoolYear = schoolYearStruct

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

func getEmployeeSchedule(idNumber string, schoolYear string) (employeeSchedule, error) {
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
		schoolYearBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte(schoolYear))
		if schoolYearBucket == nil {
			return errors.New("school year not found")
		}

		schoolYear := strings.Split(schoolYear, "-")
		schoolYearStart, schoolYearConvertErr := strconv.Atoi(schoolYear[0])
		if schoolYearConvertErr != nil {
			return schoolYearConvertErr
		}
		schoolYearEnd, schoolYearConvertErr := strconv.Atoi(schoolYear[1])
		if schoolYearConvertErr != nil {
			return schoolYearConvertErr
		}
		employeeSchedule.SchoolYear = schoolYearRange{
			StartYear: schoolYearStart,
			EndYear:   schoolYearEnd,
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

func addAttendance(idNumber string, attendanceTime attendance) error {
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
		currentYear, currentMonth, currentDay := time.Now().Date()
		currentYearString := strconv.Itoa(currentYear)
		currentMonthString := strconv.Itoa(int(currentMonth))
		currentDayString := strconv.Itoa(currentDay)

		attendanceBucket := tx.Bucket([]byte(employeeStruct.EmployeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Attendance"))

		yearBucket, yearBucketErr := attendanceBucket.CreateBucketIfNotExists([]byte(currentYearString))
		if yearBucketErr != nil {
			return yearBucketErr
		}

		monthBucket, monthBucketErr := yearBucket.CreateBucketIfNotExists([]byte(currentMonthString))
		if monthBucketErr != nil {
			return monthBucketErr
		}

		dayBucket, dayBucketErr := monthBucket.CreateBucketIfNotExists([]byte(currentDayString))
		if dayBucketErr != nil {
			return dayBucketErr
		}

		if !attendanceTime.TimeIn.DontChange {
			timeInByte, timeInMarshalErr := json.Marshal(attendanceTime.TimeIn)
			if timeInMarshalErr != nil {
				return timeInMarshalErr
			}
			timeInPutErr := dayBucket.Put([]byte("TimeIn"), timeInByte)
			if timeInPutErr != nil {
				return timeInPutErr
			}
		}

		if !attendanceTime.TimeOut.DontChange {
			timeOutByte, timeOutMarshalErr := json.Marshal(attendanceTime.TimeOut)
			if timeOutMarshalErr != nil {
				return timeOutMarshalErr
			}
			timeOutPutErr := dayBucket.Put([]byte("TimeOut"), timeOutByte)
			if timeOutPutErr != nil {
				return timeOutPutErr
			}
		}

		return nil
	})
}
