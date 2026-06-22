<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, Form, FormGroup, Input, ModalFooter } from "@sveltestrap/sveltestrap";
    import UpdateSchedule from "./UpdateSchedule.svelte";
    import RemoveSchedule from "./RemoveSchedule.svelte";
    import type { ApiRes, Schedule, SchoolYearRange, DayTimeRange } from "./utils";
    import { getSchedules } from "./utils";

	let { isModalOpen, modalToggle } = $props();

	let selectedSchoolYearSchedules: ApiRes | undefined = $state();
	let selectedSchoolYear: string | undefined = $state();
	let selectedEmployeeIdNumber: number = $state(0);
	let selectedSchedule = $state({}) as Schedule;

	function clearVars() {
		selectedSchoolYearSchedules = undefined;
		selectedSchoolYear = undefined;
		selectedEmployeeIdNumber = 0;
		selectedSchedule = {} as Schedule;
		modalToggle();
	}

	export async function init(idNumber: number) {
		addScheduleModal.setIdNumber(idNumber);
		removeScheduleModal.setIdNumber(idNumber);

		getSchedules(idNumber).then((resJson) => {
			selectedSchoolYearSchedules = resJson;
		});
	}

	function getSelectedYearSchedule(schoolYear: string): Schedule {
		let scheduleOut = {} as Schedule;
		let schoolYearSplitted = schoolYear.split("-");
		let yearRange: SchoolYearRange = {
			StartYear: Number(schoolYearSplitted[0]),
			EndYear: Number(schoolYearSplitted[1])
		}
		selectedSchoolYearSchedules?.Schedules.forEach((schedule) => {
			if (schedule.SchoolYear.StartYear == yearRange.StartYear && schedule.SchoolYear.EndYear == yearRange.EndYear) {
				scheduleOut = schedule;
			}
		})
		return scheduleOut;
	}

	function showSchedule() {
		selectedSchedule = getSelectedYearSchedule(selectedSchoolYear as string);
	}

	let addScheduleModalState = $state(false);
	let addScheduleModal: UpdateSchedule;
	function addScheduleModalToggle() {
		clearVars();
		addScheduleModalState = !addScheduleModalState;
	}

	let removeScheduleModal: RemoveSchedule;
	let removeScheduleModalState = $state(false);
	function removeScheduleModalToggle() {
		clearVars();
		removeScheduleModalState = !removeScheduleModalState;
	}
</script>

<UpdateSchedule bind:this={addScheduleModal} isModalOpen={addScheduleModalState} modalToggle={addScheduleModalToggle} />
<RemoveSchedule bind:this={removeScheduleModal} isModalOpen={removeScheduleModalState} modalToggle={removeScheduleModalToggle} refreshList={getSchedules} />

<Modal isOpen={isModalOpen} toggle={clearVars} header="View Schedules">
	<ModalBody>
		<FormGroup floating label="School Year">
			<Input type="select" bind:value={selectedSchoolYear} on:change={showSchedule}>
				{#each selectedSchoolYearSchedules?.Schedules as schoolYearSchedule }
					<option>{schoolYearSchedule.SchoolYear.StartYear + "-" + schoolYearSchedule.SchoolYear.EndYear}</option>
				{/each}
			</Input>
		</FormGroup>
		{#if selectedSchoolYear != undefined}
			<Table size="sm" striped>
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
					{#each Object.entries(selectedSchedule) as [day, time]}
						{#if day != "SchoolYear"}
							<tr>
								<td class="text-center fw-bold">{day.toUpperCase()}</td>
								<td class="text-center"><Input disabled type="switch" bind:checked={(time as DayTimeRange).DayOff} /></td>
								<td class="text-center">{(time as DayTimeRange).StartTimeHour}</td>
								<td class="text-center">{(time as DayTimeRange).StartTimeMinute}</td>
								<td class="text-center">{(time as DayTimeRange).EndTimeHour}</td>
								<td class="text-center">{(time as DayTimeRange).EndTimeMinute}</td>
							</tr>
						{/if}
					{/each}
				</tbody>
			</Table>
		{/if}
	</ModalBody>
	<ModalFooter>
		{#if selectedSchoolYear != undefined}
			<Button color="danger" disabled={selectedSchoolYear == undefined} on:click={() => {removeScheduleModal.setSchoolYear(selectedSchoolYear as string); removeScheduleModalToggle();}}>Remove Schedule</Button>
			<Button color="info" disabled={selectedSchoolYear == undefined} on:click={() => {addScheduleModal.setSchedule(getSelectedYearSchedule(selectedSchoolYear as string)); addScheduleModalToggle();}}>Edit Schedule</Button>
		{/if}
		<Button color="success" on:click={addScheduleModalToggle}>Add Schedule</Button>
	</ModalFooter>
</Modal>
