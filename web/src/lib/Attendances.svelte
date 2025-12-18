<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, Form, FormGroup, Input, ModalFooter, InputGroup } from "@sveltestrap/sveltestrap";
	import type { ApiRes, Schedule, SchoolYearRange, DayTimeRange } from "./utils";
    import { getSchedules } from "./utils";

	let { isModalOpen, modalToggle } = $props();
	let schoolYears: ApiRes | undefined = $state();
	let selectedSchoolYear: string | undefined = $state();

	function clearVars() {
		schoolYears = undefined;
		selectedSchoolYear = undefined;
		modalToggle();
	}

	export async function init(idNumber: number) {
		getSchedules(idNumber).then((resJson) => {
			schoolYears = resJson;
		});
	}

</script>

<Modal isOpen={isModalOpen} toggle={clearVars} header="View Attendances">
	<ModalBody>
		<InputGroup>
			<FormGroup floating label="School Year">
				<Input type="select" bind:value={selectedSchoolYear}>
					{#each schoolYears?.Schedules as schoolYear}
						<option>{schoolYear.SchoolYear.StartYear + "-" + schoolYear.SchoolYear.EndYear}</option>
					{/each}
				</Input>
			</FormGroup>
			<FormGroup floating label="Month">
				<Input type="select" disabled={selectedSchoolYear==undefined}>
					{#each ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"] as month, index }
						<option value={index+1}>{month}</option>
					{/each}
				</Input>
			</FormGroup>
		</InputGroup>
	</ModalBody>
	<ModalFooter>

	</ModalFooter>
</Modal>
