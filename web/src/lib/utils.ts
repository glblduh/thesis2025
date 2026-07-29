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

export interface EmployeeInfo {
	IdNumber: number
	IsFaculty: boolean
	EmployeeType: string
	FirstName: string
	MiddleName: string
	LastName: string
}

export interface InfoAttendance {
	EmployeeInfo: EmployeeInfo
	Attendances: Attendance[]
}

export interface AllAttendances {
	Faculty: InfoAttendance[]
	Staff: InfoAttendance[]
}

export interface MonthAttendances {
	SchoolYear: string
	Date: DayDate
	AttendancesEmpty: boolean
	Attendances: AllAttendances
}

export async function getSchedules(idNumber: number): Promise<ApiRes> {
	let res = await fetch("/api/getallschedule/" + idNumber)
	let jsonRes: ApiRes = await res.json();
	return jsonRes;
}

export function badgeColor(state: string): string {
	let color = "primary";
	switch (state) {
		case "DAYOFF":
			color = "secondary";
			break;
		case "LEAVE":
			color = "info";
			break;
		case "ATTENDED":
			color = "success";
			break;
		case "NOOUT":
			color = "warning";
			break;
		case "ABSENT":
			color = "danger";
			break;
	}
	return color;
}

export function logout() {
	sessionStorage.removeItem("username");
	sessionStorage.removeItem("usertype");
	sessionStorage.removeItem("userkey");
	window.location.reload();
}

export function getUsername(): string {
	return sessionStorage.getItem("username") || "";
}

export function getUserType(): string {
	return sessionStorage.getItem("usertype") || "";
}

function getUserKey(): string {
	return sessionStorage.getItem("userkey") || "";
}

function getCredentials(): string {
	return getUsername() + ":" + getUserKey();
}

export function isAuthenticated(): boolean {
	return getUsername() != "" && getUserType() != "" && getUserKey() != "";
}

export async function modFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
	const modHeaders: HeadersInit = new Headers(init?.headers || {});
	modHeaders.set("Authorization", getCredentials());

	return fetch(input, { ...init, headers: modHeaders });
}
