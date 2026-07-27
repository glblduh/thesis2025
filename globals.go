package main

import (
	"log"
	"os"
	"time"
	"errors"
)

var (
	Info  = log.New(os.Stderr, "["+time.Now().Format("2006/01/02 15:04:05")+"] [INFO] ", log.Lmsgprefix)
	Warn  = log.New(os.Stderr, "["+time.Now().Format("2006/01/02 15:04:05")+"] [WARN] ", log.Lmsgprefix)
	Error = log.New(os.Stderr, "["+time.Now().Format("2006/01/02 15:04:05")+"] [ERROR] ", log.Lmsgprefix)
)

var (
	ErrYearNoSchoolYear = errors.New("current year does not belong to any registered school year")
	ErrDayOff = errors.New("day off")
	ErrDayComplete = errors.New("attendance for today is already completed")

	ErrIdNumberNotFound = errors.New("id number not found")
	ErrInvalidEmployeeType = errors.New("not a valid employee type")

	ErrSchoolYearNotFound = errors.New("school year not found")
	ErrYearNotFound = errors.New("year not found")
	ErrMonthNotFound = errors.New("month not found")
	ErrDayNotFound = errors.New("day not found")

	ErrInvalidDate = errors.New("invalid date")
	ErrInvalidMonth = errors.New("invalid month")
	ErrInvalidDay = errors.New("invalid day")

	ErrSuspendedBucketNotFound = errors.New("suspended bucket not found")
	ErrYearBucketNotFound = errors.New("year bucket not found")
	ErrMonthBucketNotFound = errors.New("month bucket not found")

	ErrAuthUserNotFound = errors.New("user not found")
	ErrAuthKeyIncorrect = errors.New("auth key is incorrect")
)

type AttendanceState string
type AttendState     string
type SuspensionType  string
type UserType        string

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

	FACULTY UserType = "FACULTY"
	STAFF   UserType = "STAFF"
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
		EmployeeInfo employee
		Attendances []attendance
	}

	allAttendances struct {
		Faculty []monthAttendances
		Staff []monthAttendances
	}

	attend struct {
		State AttendState
		SchoolYear string
		Date dayDate
		Time attendanceTime
	}

	userAuth struct {
		Username string
		Type UserType
		Key string
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
		SchoolYear string
		Date dayDate
		Time     attendanceTime
	}

	apiAddSuspendedBodyRes struct {
		Date dayDate
		Type SuspensionType
	}

	apiRemoveSuspendedBodyRes struct {
		Date dayDate
	}

	apiGetAllMonthAttendancesRes struct {
		SchoolYear string
		Date dayDate
		AttendancesEmpty bool
		Attendances allAttendances
	}

	apiGetAllSchoolYearsBody struct {
		Faculty bool
		Staff bool
	}

	apiUserAuthBody struct {
		Username string
		Password string
	}
)
