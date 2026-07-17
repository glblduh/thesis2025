<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, FormGroup, Input, ModalFooter, InputGroup, Badge, ButtonGroup } from "@sveltestrap/sveltestrap";
	import type { AttendancesDates, DayDate, Attendance, ApiRes } from "./utils";
    import { monthsName } from "./utils";

	let { isModalOpen, modalToggle } = $props();
	let selectedDate = $state({}) as DayDate;
	let selectedSchoolYear: string | undefined = $state();
	let schoolYears = $state([]) as string[];
	let attendances = $state([]) as Attendance[];

	function clearVars() {
		selectedDate = {} as DayDate;
		selectedSchoolYear = undefined;
		schoolYears = [];
		attendances = {} as Attendance[];
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
</script>

<Modal isOpen={isModalOpen} toggle={clearVars} header="View Attendances" size="lg">
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
			<tr>
				<th scope="col" class="text-center">NAME</th>
				{#each {length: new Date(selectedDate.Year, selectedDate.Month, 0).getDate()}, day}
					<th scope="col" class="text-center">{day+1}</th>
				{/each}
			</tr>
		</thead>
		<tbody>
			{#each  as }
			{/each}
		</tbody>
		</Table>
	</ModalBody>
</Modal>
