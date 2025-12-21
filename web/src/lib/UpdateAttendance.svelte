<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Modal, ModalBody, FormGroup, Input, ModalFooter, InputGroup } from "@sveltestrap/sveltestrap";
	import type { Attendance, Schedule } from "./utils";
	import { monthsName } from "./utils";

	const defaultAttendance: Attendance = {
		Date:{Year: 0, Month: 0, Day: 0},
		State: "",
		TimeIn: {Hour: 0, Minute: 0},
		TimeOut: {Hour: 0, Minute: 0}
	};

	let { isModalOpen, modalToggle } = $props();
	let selectedEmployee: number = $state(0);
	let attendanceEdit: boolean = $state(false);
	let selectedSchedules: Schedule[] = $state([]);
	let selectedSchoolYear: string = $state("");
	let years: number[] = $state([]);
	let attendance = $state(defaultAttendance) as Attendance;
	let isLeave = $state(false);

	export function init(idNumber: number, isEdit: boolean, schedules: Schedule[], schoolYear: string, selectedAttendance: Attendance | undefined) {
		selectedEmployee = idNumber;
		attendanceEdit = isEdit;
		selectedSchedules = schedules;
		if (isEdit) {
			selectedSchoolYear = schoolYear;
			getYears();
			attendance = selectedAttendance as Attendance;
			isLeave = (selectedAttendance?.State == "LEAVE");
		}
	}

	function clearVars() {
		selectedEmployee = 0;
		attendanceEdit = false;
		selectedSchedules.length = 0;
		years.length = 0;
		selectedSchoolYear = "";
		attendance = defaultAttendance;
		isLeave = false;
		modalToggle();
	}

	function setHeader(): string {
		let header: string = "Add Attendance";
		if (attendanceEdit) {
			header = "Edit Attendance";
		}
		return header;
	}

	function getYears() {
		let splittedSchoolYear = selectedSchoolYear?.split("-");
		years = [(Number((splittedSchoolYear as string[])[0])), (Number((splittedSchoolYear as string[])[1]))];
	}

	interface ApiBody {
		IdNumber: number,
		Attendance: Attendance
	}

	async function updateAttendance() {
		if (isLeave) {
			attendance.State = "LEAVE";
		} else {
			attendance.State = "";
		}
		let body: ApiBody = {
			IdNumber: selectedEmployee,
			Attendance: attendance
		};
		await fetch("/api/updateattendance", {method: "POST", body: JSON.stringify(body)});
		clearVars();
	}
</script>

<Modal isOpen={isModalOpen} toggle={clearVars} header={setHeader()}>
	<ModalBody>
		<FormGroup floating label="School Year">
			<Input type="select" disabled={attendanceEdit} bind:value={selectedSchoolYear} on:change={getYears}>
				{#each selectedSchedules as schedule}
					<option>{schedule.SchoolYear.StartYear + "-" + schedule.SchoolYear.EndYear}</option>
				{/each}
			</Input>
		</FormGroup>
		<InputGroup>
			<FormGroup floating label="Year">
				<Input type="select" disabled={selectedSchoolYear=="" || attendanceEdit} bind:value={attendance.Date.Year}>
					{#each years as year}
						<option>{year}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Month">
				<Input type="select" disabled={attendance.Date.Year==0 || attendanceEdit} bind:value={attendance.Date.Month}>
					{#each monthsName as month, index}
						<option value={index+1}>{month}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Day">
				<Input type="select" disabled={attendance.Date.Month==0 || attendanceEdit} bind:value={attendance.Date.Day}>
					{#each Array(new Date(attendance.Date.Year, attendance.Date.Month, 0).getDate()) as _, day}
						<option>{day+1}</option>
					{/each}
				</Input>
			</FormGroup>
		</InputGroup>
		<InputGroup>
			<FormGroup floating label="Time In Hour">
				<Input disabled={attendance.Date.Day==0 || isLeave} type="select" bind:value={attendance.TimeIn.Hour}>
					{#each Array.from(Array(24).keys()) as hour }
						<option>{hour}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Time In Minute">
				<Input disabled={attendance.Date.Day==0 || isLeave} type="select" bind:value={attendance.TimeIn.Minute}>
					{#each Array.from(Array(60).keys()) as minute }
						<option>{minute}</option>
					{/each}
				</Input>
			</FormGroup>
		</InputGroup>
		<InputGroup>
			<FormGroup floating label="Time Out Hour">
				<Input disabled={attendance.Date.Day==0 || isLeave} type="select" bind:value={attendance.TimeOut.Hour}>
					{#each Array.from(Array(24).keys()) as hour }
						<option>{hour}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Time Out Minute">
				<Input disabled={attendance.Date.Day==0 || isLeave} type="select" bind:value={attendance.TimeOut.Minute}>
					{#each Array.from(Array(60).keys()) as minute }
						<option>{minute}</option>
					{/each}
				</Input>
			</FormGroup>
		</InputGroup>
		<Input disabled={attendance.Date.Day==0} type="switch" label="Leave?" bind:checked={isLeave} />
	</ModalBody>
	<ModalFooter>
		<Button color="secondary" on:click={clearVars}>Back</Button>
		<Button color="success" on:click={updateAttendance}>Submit</Button>
	</ModalFooter>
</Modal>
