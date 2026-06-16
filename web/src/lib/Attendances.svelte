<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, FormGroup, Input, ModalFooter, InputGroup, Badge, ButtonGroup } from "@sveltestrap/sveltestrap";
	import type { AttendancesDates, DayDate, Attendance, ApiRes } from "./utils";
    import { getSchedules, monthsName } from "./utils";
    import UpdateAttendance from "./UpdateAttendance.svelte";
    import RemoveAttendance from "./RemoveAttendance.svelte";

	let { isModalOpen, modalToggle } = $props();
	let selectedEmployee: number = 0;
	let selectedDate = $state({}) as DayDate;
	let dates = $state({}) as AttendancesDates;
	let attendances = $state([]) as Attendance[];
	let schedules = $state({}) as ApiRes;
	let selectedSchoolYear: string | undefined = $state();

	function clearVars() {
		selectedEmployee = 0;
		selectedDate = {} as DayDate;
		dates = {} as AttendancesDates;
		attendances = {} as Attendance[];
		schedules = {} as ApiRes;
		selectedSchoolYear = undefined;
		modalToggle();
	}

	async function getSchoolYears() {
		getSchedules(selectedEmployee).then((resJson: ApiRes) => {
			schedules = resJson;
		})
	}

	function getYears() {
		let splittedSchoolYear = selectedSchoolYear?.split("-");
		dates.Years = [(Number((splittedSchoolYear as string[])[0])), (Number((splittedSchoolYear as string[])[1]))];
	}

	async function getMonths() {
		attendances = {} as Attendance[];
		fetch("/api/getallattendancesmonths/" + selectedEmployee + "/" + selectedDate.Year).then((res) => {
			res.json().then((resJson: AttendancesDates) => {
				dates.Months = resJson.Months;
			})
		})
	}

	async function getAttendances() {
		fetch("/api/getmonthattendances/" + selectedEmployee + "/" + selectedSchoolYear + "/" + selectedDate.Year + "/" + selectedDate.Month).then((res) => {
			res.json().then((resJson) => {
				attendances = resJson.Attendances
			})
		})
	}

	function badgeColor(state: string): string {
		let color = "primary";
		switch(state) {
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

	export async function init(idNumber: number) {
		selectedEmployee = idNumber;
		await getSchoolYears();
	}

	let updateAttendanceModal: UpdateAttendance;
	let updateAttendanceModalState = $state(false);
	function updateAttendanceModalToggle() {
		clearVars();
		updateAttendanceModalState = !updateAttendanceModalState;
	}

	let removeAttendanceModal: RemoveAttendance;
	let removeAttendanceModalState = $state(false);
	function removeAttendanceModalToggle() {
		clearVars();
		removeAttendanceModalState = !removeAttendanceModalState;
	}
</script>

<UpdateAttendance bind:this={updateAttendanceModal} isModalOpen={updateAttendanceModalState} modalToggle={updateAttendanceModalToggle} />
<RemoveAttendance bind:this={removeAttendanceModal} isModalOpen={removeAttendanceModalState} modalToggle={removeAttendanceModalToggle} />

<Modal isOpen={isModalOpen} toggle={clearVars} header="View Attendances" size="lg">
	<ModalBody>
		<FormGroup floating label="School Year">
			<Input type="select" bind:value={selectedSchoolYear} on:change={getYears}>
				{#each schedules?.Schedules as schedule}
					<option>{schedule.SchoolYear.StartYear + "-" + schedule.SchoolYear.EndYear}</option>
				{/each}
			</Input>
		</FormGroup>
		<InputGroup>
			<FormGroup floating label="Year">
				<Input type="select" disabled={selectedSchoolYear==undefined} bind:value={selectedDate.Year} on:change={getMonths}>
					{#each dates.Years as year}
						<option>{year}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Month">
				<Input type="select" disabled={dates.Months==undefined} on:change={getAttendances} bind:value={selectedDate.Month}>
					{#each dates.Months as month }
						<option value={month}>{monthsName[month-1]}</option>
					{/each}
				</Input>
			</FormGroup>
		</InputGroup>
		<Table striped size="sm" responsive>
			<thead>
				<tr>
					<th scope="col" class="text-center">DATE</th>
					<th scope="col" class="text-center">STATE</th>
					<th scope="col" class="text-center">TIME IN HOUR</th>
					<th scope="col" class="text-center">TIME IN MINUTE</th>
					<th scope="col" class="text-center">TIME OUT HOUR</th>
					<th scope="col" class="text-center">TIME OUT MINUTE</th>
					<th scope="col" class="text-center">ACTION</th>
				</tr>
			</thead>
			<tbody>
				{#each attendances as attendance}
					<tr>
						<td class="text-center fw-bold">{attendance.Date.Day}</td>
						<td class="text-center fw-bold">
							{#if attendance.Suspended != "NOTSUSPENDED" }
								<Badge color="info">{attendance.Suspended}</Badge>
							{/if}
							<Badge color={badgeColor(attendance.State)}>{attendance.State}</Badge>
						</td>
						<td class="text-center">{attendance.TimeIn.Hour}</td>
						<td class="text-center">{attendance.TimeIn.Minute}</td>
						<td class="text-center">{attendance.TimeOut.Hour}</td>
						<td class="text-center">{attendance.TimeOut.Minute}</td>
						{#if attendance.State != "ABSENT" && attendance.State != "DAYOFF"}
							<td>
								<ButtonGroup size="sm">
									<Button color="info" on:click={() => {
										updateAttendanceModal.init(selectedEmployee, true, schedules?.Schedules, selectedSchoolYear as string, attendance);
										updateAttendanceModalToggle();}}>EDIT</Button>
									<Button color="danger" on:click={() => {
										removeAttendanceModal.init(selectedEmployee, attendance.Date);
										removeAttendanceModalToggle();
									}}>REMOVE</Button>
								</ButtonGroup>
							</td>
						{:else}
							<td></td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</Table>
	</ModalBody>
	<ModalFooter>
		<Button color="success" on:click={() => {updateAttendanceModal.init(selectedEmployee, false, schedules?.Schedules, selectedSchoolYear as string, undefined); updateAttendanceModalToggle();}}>Add Attendance</Button>
	</ModalFooter>
</Modal>
