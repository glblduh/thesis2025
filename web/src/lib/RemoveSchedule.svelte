<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Modal, ModalBody, ModalFooter } from "@sveltestrap/sveltestrap";

	let { isModalOpen, modalToggle, refreshList } = $props();
	let idNumber: number = 0;
	let schoolYear: string = $state("");

	export function setIdNumber(employeeIdNumber: number) {
		idNumber = employeeIdNumber;
	}

	export function setSchoolYear(selectedSchoolYear: string) {
		schoolYear = selectedSchoolYear;
	}

	interface apiBody {
		IdNumber: number,
		SchoolYear: string;
	}

	async function removeEmployee() {
		let body: apiBody = {
			IdNumber: idNumber,
			SchoolYear: schoolYear
		};
		await fetch("/api/removeschedule", {method: "DELETE", body: JSON.stringify(body)});
		refreshList(idNumber);
		modalToggle();
	}
</script>

<Modal isOpen={isModalOpen} toggle={modalToggle} header="Remove Schedule">
	<ModalBody>
		Are you sure to remove schedule for school year <span class="fw-bold">{schoolYear}</span>? This action is irreversible.
	</ModalBody>
	<ModalFooter>
		<Button color="secondary" on:click={modalToggle}>Back</Button>
		<Button color="danger" on:click={removeEmployee}>Remove</Button>
	</ModalFooter>
</Modal>
