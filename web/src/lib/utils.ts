export interface SchoolYearRange {
		StartYear: number,
		EndYear: number
}

export interface DayTimeRange {
	DayOff: boolean,
	StartTimeHour: number,
	StartTimeMinute: number,
	EndTimeHour: number,
	EndTimeMinute: number
}

export interface Schedule {
	SchoolYear: SchoolYearRange,
	Monday: DayTimeRange,
	Tuesday: DayTimeRange,
	Wednesday: DayTimeRange,
	Thursday: DayTimeRange,
	Friday: DayTimeRange,
	Saturday: DayTimeRange,
	Sunday: DayTimeRange
}

export interface ApiRes {
	IdNumber: number,
	Schedules: Schedule[]
}

export interface DayDate {
	Year: number,
	Month: number,
	Day: number
}

export const monthsName: string[] = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];

export interface AttendancesDates {
	Years: number[],
	Months: number[],
	Days: number[],
}

export interface GetDatesBody {
	IdNumber: number,
	Date: DayDate
}

export interface AttendanceTime {
	Hour: number
	Minute: number
}

export interface Attendance {
	Date: DayDate
	State: string
	TimeIn: AttendanceTime
	TimeOut: AttendanceTime
	Suspended: string
}

export interface SuspendedDay {
	Date: DayDate
	Type: string
}

export async function getSchedules(idNumber: number): Promise<ApiRes> {
	let res = await fetch("/api/getallschedule/" + idNumber)
	let jsonRes: ApiRes = await res.json();
	return jsonRes;
}
