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
		dontChange      bool
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
		IdNumber   int
		FirstName  string
		MiddleName string
		LastName   string
	}
)

type (
	jsonErrorRes struct {
		Error string
	}

	apiAddEmployeeBody struct {
		IdNumber   int
		isFaculty  bool
		FirstName  string
		MiddleName string
		LastName   string
	}

	apiAddEmployeeRes struct {
		IdNumber  int
		isFaculty bool
	}

	apiRemoveEmployeeBodyRes struct {
		IdNumber  int
		isFaculty bool
	}

	apiUpdateScheduleBody struct {
		IdNumber  int
		isFaculty bool
		Schedule  employeeSchedule
	}
)
