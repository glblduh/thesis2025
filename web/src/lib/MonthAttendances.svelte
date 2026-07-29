<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, FormGroup, Input, ModalFooter, InputGroup, Badge, ButtonGroup } from "@sveltestrap/sveltestrap";
	import type { DayDate, MonthAttendances } from "./utils";
    import { monthsName, badgeColor, modFetch } from "./utils";
    import { untrack } from "svelte";

	let { isModalOpen, modalToggle } = $props();
	let selectedDate = $state({}) as DayDate;

	let selectedSchoolYear: string | undefined = $state();
	$effect(() => {
		if (selectedSchoolYear == undefined) {
			return
		}

		selectedDate.Year = untrack(() => years[0]);
		untrack(() => getAttendances());
	});

	let schoolYears = $state([]) as string[];
	let attendances = $state({}) as MonthAttendances;
	let years = $state([]) as number[];

	function clearVars() {
		selectedDate = {} as DayDate;
		selectedSchoolYear = undefined;
		schoolYears = [];
		attendances = {} as MonthAttendances;
		years = [];
		modalToggle();
	}

	export async function getSchoolYears() {
		modFetch("/api/getallschoolyears").then((res) => {
			res.json().then((resJson) => {
				schoolYears = resJson
			})
		})
	}

	async function getAttendances() {

		if (selectedSchoolYear == undefined || selectedDate.Year == undefined || selectedDate.Month == undefined) {
			return
		}

		modFetch("/api/getallmonthattendances/" + selectedSchoolYear + "/" + selectedDate.Year + "/" + selectedDate.Month).then((res) => {
			res.json().then((resJson) => {
				attendances = resJson
			})
		})
	}

	function getYears() {
		let splitted = selectedSchoolYear?.split("-") as string[];
		years = [Number(splitted[0]), Number(splitted[1])];
	}
</script>

<Modal isOpen={isModalOpen} toggle={clearVars} header="View All Attendances" size="lg">
	<ModalBody>
		<InputGroup>
			<FormGroup floating label="School Year">
				<Input type="select" bind:value={selectedSchoolYear} on:change={getYears}>
					{#each schoolYears as schoolYear}
						<option>{schoolYear}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Year">
				<Input type="select" disabled={selectedSchoolYear==undefined} bind:value={selectedDate.Year} on:change={getAttendances}>
					{#each years as year}
						<option>{year}</option>
					{/each}
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
		{#if attendances.AttendancesEmpty != undefined && !attendances.AttendancesEmpty}
			<hr />
			<Table size="sm" responsive>
				<thead>
					<tr>
						<th scope="col" class="text-center">NAME</th>
						<th scope="col" class="text-center">TYPE</th>
						{#each {length: new Date(selectedDate.Year, selectedDate.Month, 0).getDate()}, day}
							<th scope="col" class="text-center">{day+1}</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each attendances.Attendances as attendance}
						{#if attendance.Attendances?.length != 0}
							<tr>
								<td>{attendance.EmployeeInfo.LastName + ", " + attendance.EmployeeInfo.FirstName}</td>
								<td>{attendance.EmployeeInfo.EmployeeType}</td>
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
		{/if}
	</ModalBody>
</Modal>
