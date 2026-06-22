<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Modal, ModalBody, ModalFooter } from "@sveltestrap/sveltestrap";
    import { type DayDate } from "./utils";

	let { isModalOpen, modalToggle } = $props();
	let selectedDate = $state({}) as DayDate;

	export function init(date: DayDate) {
		selectedDate = date;
	}

	interface ApiBody {
		Date: DayDate
	}

	async function removeSuspended() {
		let body: ApiBody = {
			Date: selectedDate
		}
		await fetch("/api/removesuspended", {method: "DELETE", body: JSON.stringify(body)});
		modalToggle();
	}
</script>

<Modal isOpen={isModalOpen} toggle={modalToggle} header="Remove Suspension">
	<ModalBody>
		Are you sure to remove the suspension on <span class="fw-bold">{selectedDate.Month}/{selectedDate.Day}/{selectedDate.Year}</span>? This action is irreversible.
	</ModalBody>
	<ModalFooter>
		<Button color="secondary" on:click={modalToggle}>Back</Button>
		<Button color="danger" on:click={removeSuspended}>Remove</Button>
	</ModalFooter>
</Modal>
