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

export async function getSchedules(idNumber: number): Promise<ApiRes> {
	let res = await fetch("/api/getallschedule/" + idNumber)
	let jsonRes = await res.json();
	return jsonRes;
}
