<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Modal, ModalBody, FormGroup, Input, Table, InputGroup, Form, ModalFooter } from "@sveltestrap/sveltestrap";
	import type { ApiRes, Schedule, SchoolYearRange, DayTimeRange } from "./utils";

	let { isModalOpen, modalToggle } = $props();

	interface AddScheduleBody {
		IdNumber: number,
		Schedule: Schedule
	}

	const defaultSchedulesState: Schedule = {
		SchoolYear: {StartYear: 0, EndYear: 0},
		Monday: {DayOff: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Tuesday: {DayOff: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Wednesday: {DayOff: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Thursday: {DayOff: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Friday: {DayOff: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Saturday: {DayOff: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
		Sunday: {DayOff: false, StartTimeHour: 0, StartTimeMinute: 0, EndTimeHour: 0, EndTimeMinute: 0},
	};

	let schedules = $state(defaultSchedulesState) as Schedule;
	let editSchedule = $state(false);
	let idNumber: number = 0;

	export function setIdNumber(id: number) {
		idNumber = id;
	}

	export function setSchedule(selectedSchedule: Schedule) {
		editSchedule = true;
		schedules = selectedSchedule;
	}

	function clearVars() {
		schedules = defaultSchedulesState;
		editSchedule = false;
		idNumber = 0;
		modalToggle();
	}

	function setHeader(): string {
		let header = "Add Schedule";
		if (editSchedule) {
			header = "Edit Schedule";
		}
		return header;
	}

	async function addSchedule() {
		let requestBody: AddScheduleBody = {
			IdNumber: idNumber,
			Schedule: schedules
		};
		let apiRes = await fetch("/api/updateschedule", {method: "POST", body: JSON.stringify(requestBody)});
		clearVars();
	}
</script>

<Modal isOpen={isModalOpen} toggle={clearVars} header={setHeader()} size="lg">
	<ModalBody>
		<InputGroup>
			<FormGroup floating label="School Year Start">
				<Input type="number" disabled={editSchedule} bind:value={schedules.SchoolYear.StartYear} />
			</FormGroup>
			<FormGroup floating label="School Year End">
				<Input type="number" disabled={editSchedule} bind:value={schedules.SchoolYear.EndYear} />
			</FormGroup>
		</InputGroup>
		<Table size="sm" striped responsive>
			<thead>
				<tr>
					<th scope="col" class="text-center">DAY</th>
					<th scope="col" class="text-center">OFF?</th>
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
							<td class="text-center fw-bold">{day.toUpperCase()}</td>
							<td><Input type="checkbox" bind:checked={(time as DayTimeRange).DayOff} /></td>
							<td>
								<Input bind:disabled={(time as DayTimeRange).DayOff} type="select" bind:value={(time as DayTimeRange).StartTimeHour}>
									{#each Array.from(Array(24).keys()) as hour }
										<option>{hour}</option>
									{/each}
								</Input>
							</td>
							<td>
								<Input bind:disabled={(time as DayTimeRange).DayOff} type="select" bind:value={(time as DayTimeRange).StartTimeMinute}>
									{#each Array.from(Array(60).keys()) as minute }
										<option>{minute}</option>
									{/each}
								</Input>
							</td>
							<td>
								<Input bind:disabled={(time as DayTimeRange).DayOff} type="select" bind:value={(time as DayTimeRange).EndTimeHour}>
									{#each Array.from(Array(24).keys()) as hour }
										<option>{hour}</option>
									{/each}
								</Input>
							</td>
							<td>
								<Input bind:disabled={(time as DayTimeRange).DayOff} type="select" bind:value={(time as DayTimeRange).EndTimeMinute}>
									{#each Array.from(Array(60).keys()) as minute }
										<option>{minute}</option>
									{/each}
								</Input>

							</td>
						</tr>
					{/if}
				{/each}
			</tbody>
		</Table>
	</ModalBody>
	<ModalFooter>
		<Button color="secondary" on:click={clearVars}>Back</Button>
		<Button color="success" type="submit" on:click={addSchedule}>Submit</Button>
	</ModalFooter>
</Modal>
