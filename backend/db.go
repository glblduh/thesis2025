package main

import (
	"encoding/json"
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
		_, scheduleBucketErr := bucket.CreateBucketIfNotExists([]byte("Schedule"))
		if scheduleBucketErr != nil {
			return scheduleBucketErr
		}
		_, attendanceBucketErr := bucket.CreateBucketIfNotExists([]byte("Attendance"))
		if attendanceBucketErr != nil {
			return attendanceBucketErr
		}
		_, absenceBucketErr := bucket.CreateBucketIfNotExists([]byte("Absence"))
		if absenceBucketErr != nil {
			return absenceBucketErr
		}
		_, leaveBucketErr := bucket.CreateBucketIfNotExists([]byte("Leave"))
		if leaveBucketErr != nil {
			return leaveBucketErr
		}
		return nil
	})
}

func removeEmployee(idNumber string, isFaculty bool) error {
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
		return tx.Bucket([]byte(employeeType)).DeleteBucket([]byte(idNumber))
	})
}

func updateEmployeeSchedule(idNumber string, isFaculty bool, schedule employeeSchedule) error {
	db, dbErr := openDB()
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	employeeType := "Staff"
	if isFaculty {
		employeeType = "Faculty"
	}

	schoolYear := strconv.Itoa(schedule.SchoolYear.StartYear) + strconv.Itoa(schedule.SchoolYear.EndYear)

	mondaySchedule, scheduleMarshalErr := json.Marshal(schedule.Monday)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}
	tuesdaySchedule, scheduleMarshalErr := json.Marshal(schedule.Tuesday)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}
	wednesdaySchedule, scheduleMarshalErr := json.Marshal(schedule.Wednesday)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}
	thursdaySchedule, scheduleMarshalErr := json.Marshal(schedule.Thursday)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}
	fridaySchedule, scheduleMarshalErr := json.Marshal(schedule.Friday)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}
	saturdaySchedule, scheduleMarshalErr := json.Marshal(schedule.Saturday)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}
	sundaySchedule, scheduleMarshalErr := json.Marshal(schedule.Sunday)
	if scheduleMarshalErr != nil {
		return scheduleMarshalErr
	}

	return db.Update(func(tx *bbolt.Tx) error {
		scheduleBucket := tx.Bucket([]byte(employeeType)).Bucket([]byte(idNumber)).Bucket([]byte("Schedule"))
		schoolYearBucket, schoolYearBucketErr := scheduleBucket.CreateBucketIfNotExists([]byte(schoolYear))
		if schoolYearBucketErr != nil {
			return schoolYearBucketErr
		}

		var scheduleAddErr error
		scheduleAddErr = schoolYearBucket.Put([]byte("Monday"), mondaySchedule)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}
		scheduleAddErr = schoolYearBucket.Put([]byte("Tuesday"), tuesdaySchedule)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}
		scheduleAddErr = schoolYearBucket.Put([]byte("Wednesday"), wednesdaySchedule)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}
		scheduleAddErr = schoolYearBucket.Put([]byte("Thursday"), thursdaySchedule)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}
		scheduleAddErr = schoolYearBucket.Put([]byte("Friday"), fridaySchedule)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}
		scheduleAddErr = schoolYearBucket.Put([]byte("Saturday"), saturdaySchedule)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}
		scheduleAddErr = schoolYearBucket.Put([]byte("Sunday"), sundaySchedule)
		if scheduleAddErr != nil {
			return scheduleAddErr
		}

		return nil
	})
}
