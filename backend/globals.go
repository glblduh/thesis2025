package main

import (
	"log"
	"os"
	"time"
)

var (
	Info  = log.New(os.Stderr, "["+time.Now().Format("2006/01/02 15:04:05")+"] [INFO] ", log.Lmsgprefix)
	Warn  = log.New(os.Stderr, "["+time.Now().Format("2006/01/02 15:04:05")+"] [WARN] ", log.Lmsgprefix)
	Error = log.New(os.Stderr, "["+time.Now().Format("2006/01/02 15:04:05")+"] [ERROR] ", log.Lmsgprefix)
)

type (
	schoolYearRange struct {
		StartYear int
		EndYear   int
	}

	dayTimeRange struct {
		DontChange      bool
		StartTimeHour   int
		StartTimeMinute int
		EndTimeHour     int
		EndTimeMinute   int
	}

	employeeSchedule struct {
		SchoolYear schoolYearRange
		Monday     dayTimeRange
		Tuesday    dayTimeRange
		Wednesday  dayTimeRange
		Thursday   dayTimeRange
		Friday     dayTimeRange
		Saturday   dayTimeRange
		Sunday     dayTimeRange
	}

	employee struct {
		IdNumber     int
		IsFaculty    bool
		EmployeeType string
		FirstName    string
		MiddleName   string
		LastName     string
	}

	attendanceTime struct {
		DontChange bool
		Hour       int
		Minute     int
		Unix       int
	}

	attendance struct {
		State   string
		Reason  string
		TimeIn  attendanceTime
		TimeOut attendanceTime
	}

	dayDate struct {
		Year  int
		Month int
		Day   int
	}
)

type (
	jsonErrorRes struct {
		Error string
	}

	apiAddEmployeeBodyRes struct {
		IdNumber   int
		IsFaculty  bool
		FirstName  string
		MiddleName string
		LastName   string
	}

	apiRemoveEmployeeBodyRes struct {
		IdNumber int
	}

	apiUpdateScheduleBodyRes struct {
		IdNumber int
		Schedule employeeSchedule
	}

	apiGetAllYearsScheduleRes struct {
		IdNumber  int
		Schedules []employeeSchedule
	}

	apiGetScheduleRes struct {
		IdNumber int
		Schedule employeeSchedule
	}
)
