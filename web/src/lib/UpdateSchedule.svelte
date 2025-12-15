<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Modal, ModalBody, FormGroup, Input, Table, InputGroup } from "@sveltestrap/sveltestrap";

	let { isModalOpen, modalToggle } = $props();

	interface schoolYearRange {
		StartYear: number,
		EndYear: number
	}

	interface DayTimeRange {
		Change: boolean,
		StartTimeHour: number,
		StartTimeMinute: number,
		EndTimeHour: number,
		EndTimeMinute: number
	}

	interface Schedule {
		SchoolYear: schoolYearRange,
		Monday: DayTimeRange,
		Tuesday: DayTimeRange,
		Wednesday: DayTimeRange,
		Thursday: DayTimeRange,
		Friday: DayTimeRange,
		Saturday: DayTimeRange,
		Sunday: DayTimeRange
	}

	let schedules = $state({
		SchoolYear: {StartYear: 0, EndYear: 0},
		Monday: {Change: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Tuesday: {Change: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Wednesday: {Change: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Thursday: {Change: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Friday: {Change: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Saturday: {Change: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Sunday: {Change: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
	}) as Schedule;
	let editSchedule = $state(false);

	export function setSchedule(selectedSchedule: Schedule) {
		editSchedule = true;
		schedules = selectedSchedule;
	}

	function clearVars() {
		editSchedule = false;
		modalToggle();
	}

	function setHeader(): string {
		let header = "Add Schedule";
		if (editSchedule) {
			header = "Edit Schedule";
		}
		return header;
	}
</script>

<Modal isOpen={isModalOpen} toggle={clearVars} header={setHeader()}>
	<ModalBody>
		<InputGroup>
			<FormGroup floating label="School Year Start">
				<Input type="number" disabled={editSchedule} bind:value={schedules.SchoolYear.StartYear} />
			</FormGroup>
			<FormGroup floating label="School Year End">
				<Input type="number" disabled={editSchedule} bind:value={schedules.SchoolYear.EndYear} />
			</FormGroup>
		</InputGroup>
		<Table size="sm" striped>
			<thead>
				<tr>
					<th scope="col" class="text-center">CHANGE?</th>
					<th scope="col" class="text-center">DAY</th>
					<th scope="col" class="text-center">START HOUR</th>
					<th scope="col" class="text-center">START MINUTE</th>
					<th scope="col" class="text-center">END HOUR</th>
					<th scope="col" class="text-center">END MINUTE</th>
				</tr>
			</thead>
			<tbody>
				{#each Object.entries(schedules) as [day, time]}
					{#if day != "SchoolYear"}
						<tr>
							<td class="align-center"><Input type="checkbox" bind:checked={(time as DayTimeRange).Change} /></td>
							<td class="text-center fw-bold">{day.toUpperCase()}</td>
							<td><Input type="number" disabled={!(time as DayTimeRange).Change} min=0 max=23 bind:value={(time as DayTimeRange).StartTimeHour} /></td>
							<td><Input type="number" disabled={!(time as DayTimeRange).Change} min=0 max=59 bind:value={(time as DayTimeRange).StartTimeMinute} /></td>
							<td><Input type="number" disabled={!(time as DayTimeRange).Change} min=0 max=23 bind:value={(time as DayTimeRange).EndTimeHour} /></td>
							<td><Input type="number" disabled={!(time as DayTimeRange).Change} min=0 max=59 bind:value={(time as DayTimeRange).EndTimeMinute} /></td>
						</tr>
					{/if}
				{/each}
			</tbody>
		</Table>
	</ModalBody>
</Modal>
