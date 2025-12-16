<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import { Button, Modal, ModalBody, ModalFooter } from "@sveltestrap/sveltestrap";

	let { isModalOpen, modalToggle, refreshList, idNumber } = $props();

	interface removeEmployeeBody {
		IdNumber: number;
	}

	async function removeEmployee() {
		let jsonBody: removeEmployeeBody = {
			IdNumber: idNumber
		}
		let removeEmployeeRes = await fetch("/api/removeemployee", {method: "DELETE", body: JSON.stringify(jsonBody)});
		refreshList();
		modalToggle();
	}
</script>

<Modal isOpen={isModalOpen} toggle={modalToggle} header="Remove Employee">
	<ModalBody>
		Are you sure to remove employee number <span class="fw-bold">{idNumber}</span>? This action is irreversible.
	</ModalBody>
	<ModalFooter>
		<Button color="secondary" on:click={modalToggle}>No</Button>
		<Button color="danger" on:click={removeEmployee}>Yes</Button>
	</ModalFooter>
</Modal>
