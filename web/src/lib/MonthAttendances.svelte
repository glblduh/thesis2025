<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, FormGroup, Input, ModalFooter, InputGroup, Badge, ButtonGroup } from "@sveltestrap/sveltestrap";
	import type { DayDate, MonthAttendances } from "./utils";
    import { monthsName, badgeColor } from "./utils";

	let { isModalOpen, modalToggle } = $props();
	let selectedDate = $state({}) as DayDate;
	let selectedSchoolYear: string | undefined = $state();
	let schoolYears = $state([]) as string[];
	let attendances = $state({}) as MonthAttendances;
	let showHeader = $state(false);

	function clearVars() {
		selectedDate = {} as DayDate;
		selectedSchoolYear = undefined;
		schoolYears = [];
		attendances = {} as MonthAttendances;
		showHeader = false;
		modalToggle();
	}

	export async function getSchoolYears() {
		fetch("/api/getallschoolyears").then((res) => {
			res.json().then((resJson) => {
				schoolYears = resJson
			})
		})
	}

	async function getAttendances() {
		fetch("/api/getallmonthattendances/" + selectedSchoolYear + "/" + selectedDate.Year + "/" + selectedDate.Month).then((res) => {
			res.json().then((resJson) => {
				attendances = resJson
			})
		})
	}

	function checkShowHeader() {
		if (!showHeader) {
			showHeader = true;
		}
	}
</script>

<Modal isOpen={isModalOpen} toggle={clearVars} header="View All Attendances" size="lg">
	<ModalBody>
		<InputGroup>
			<FormGroup floating label="School Year">
				<Input type="select" bind:value={selectedSchoolYear}>
					{#each schoolYears as schoolYear}
						<option>{schoolYear}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Year">
				<Input type="select" disabled={selectedSchoolYear==undefined} bind:value={selectedDate.Year}>
						<option>{selectedSchoolYear?.split("-")[0]}</option>
						<option>{selectedSchoolYear?.split("-")[1]}</option>
				</Input>
			</FormGroup>
			<FormGroup floating label="Month">
				<Input type="select" disabled={selectedDate.Year==undefined} on:change={getAttendances} bind:value={selectedDate.Month}>
					{#each {length: 12}, month}
						<option value={month+1}>{monthsName[month]}</option>
					{/each}
				</Input>
			</FormGroup>
		</InputGroup>
		<Table striped size="sm" responsive>
			<thead>
				{#if attendances.AttendancesEmpty != undefined && !attendances.AttendancesEmpty}
					<tr>
						<th scope="col" class="text-center">NAME</th>
						{#each {length: new Date(selectedDate.Year, selectedDate.Month, 0).getDate()}, day}
							<th scope="col" class="text-center">{day+1}</th>
						{/each}
					</tr>
				{/if}
			</thead>
			<tbody>
				{#each attendances.Attendances?.Faculty as attendance}
					{#if attendance.Attendances?.length != 0}
						<tr class="table-primary">
							<td>{attendance.EmployeeInfo.LastName + ", " + attendance.EmployeeInfo.FirstName}</td>
							{#each attendance.Attendances as dayAttendance}
								<td>
									<Badge color={badgeColor(dayAttendance.State)}>{dayAttendance.State}</Badge>
									{#if dayAttendance.State == "ATTENDED"}
										<Badge color="info">IN: {dayAttendance.TimeIn.Hour + ":" + dayAttendance.TimeIn.Minute}</Badge>
										<Badge color="info">OUT: {dayAttendance.TimeOut.Hour + ":" + dayAttendance.TimeOut.Minute}</Badge>
									{/if}
									{#if dayAttendance.State == "NOOUT"}
										<Badge color="info">IN: {dayAttendance.TimeIn.Hour + ":" + dayAttendance.TimeIn.Minute}</Badge>
									{/if}
								</td>
							{/each}
						</tr>
					{/if}
				{/each}
				{#each attendances.Attendances?.Staff as attendance}
					{#if attendance.Attendances?.length != 0}
						<tr class="table-secondary">
							<td>{attendance.EmployeeInfo.LastName + ", " + attendance.EmployeeInfo.FirstName}</td>
							{#each attendance.Attendances as dayAttendance}
								<td>
									<Badge color={badgeColor(dayAttendance.State)}>{dayAttendance.State}</Badge>
									{#if dayAttendance.State == "ATTENDED"}
										<Badge color="info">IN: {dayAttendance.TimeIn.Hour + ":" + dayAttendance.TimeIn.Minute}</Badge>
										<Badge color="info">OUT: {dayAttendance.TimeOut.Hour + ":" + dayAttendance.TimeOut.Minute}</Badge>
									{/if}
									{#if dayAttendance.State == "NOOUT"}
										<Badge color="info">IN: {dayAttendance.TimeIn.Hour + ":" + dayAttendance.TimeIn.Minute}</Badge>
									{/if}
								</td>
							{/each}
						</tr>
					{/if}
				{/each}
			</tbody>
		</Table>
	</ModalBody>
</Modal>
