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

type AttendanceState string
type AttendState     string
type SuspensionType  string

const (
	DAYOFF     AttendanceState = "DAYOFF"
	LEAVE      AttendanceState = "LEAVE"
	ATTENDED   AttendanceState = "ATTENDED"
	NOOUT      AttendanceState = "NOOUT"
	ABSENT     AttendanceState = "ABSENT"

	TIMEIN     AttendState     = "TIMEIN"
	TIMEOUT    AttendState     = "TIMEOUT"

	NOTSUSPENDED SuspensionType = "NOTSUSPENDED"
	SUSPENSION   SuspensionType = "SUSPENSION"
	HOLIDAY      SuspensionType = "HOLIDAY"
)

type (
	schoolYearRange struct {
		StartYear int
		EndYear   int
	}

	dayTimeRange struct {
		DayOff          bool
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
		Hour   int
		Minute int
		Unix   int
	}

	attendance struct {
		Date      dayDate
		State     AttendanceState
		TimeIn    attendanceTime
		TimeOut   attendanceTime
		Suspended SuspensionType
	}

	dayDate struct {
		Year  int
		Month int
		Day   int
	}

	allEmployees struct {
		Faculty []employee
		Staff   []employee
	}

	attendanceDates struct {
		Years  []int
		Months []int
		Days   []int
	}

	suspendedDay struct {
		Date dayDate
		Type SuspensionType
	}

	monthAttendances struct {
		IdNumber int
		Attendances []attendance
	}

	allAttendances struct {
		Faculty []monthAttendances
		Staff []monthAttendances
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

	apiUpdateAttendanceBodyRes struct {
		IdNumber       int
		Attendance     attendance
	}

	apiGetAttendanceBody struct {
		IdNumber   int
		Date       dayDate
	}

	apiGetAttendanceRes struct {
		IdNumber int
		State    AttendanceState
		Reason   string
		TimeIn   attendanceTime
		TimeOut  attendanceTime
	}

	apiGetMonthAttendancesRes struct {
		Attendances []attendance
	}

	apiRemoveScheduleBodyRes struct {
		IdNumber   int
		SchoolYear string
	}

	apiRemoveAttendanceBodyRes struct {
		IdNumber   int
		Date       dayDate
	}

	apiAttendRes struct {
		IdNumber int
		State    AttendState
		Time     attendanceTime
	}

	apiAddSuspendedBodyRes struct {
		Date dayDate
		Type SuspensionType
	}

	apiRemoveSuspendedBodyRes struct {
		Date dayDate
	}

	apiGetAllMonthAttendancesBody struct {
		SchoolYear string
		Date dayDate
	}

	apiGetAllMonthAttendancesRes struct {
		SchoolYear string
		Date dayDate
		Employees allEmployees
		Attendances allAttendances
	}
)
