<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Modal, ModalBody, ModalFooter } from "@sveltestrap/sveltestrap";
    import { type DayDate, modFetch } from "./utils";

	let { isModalOpen, modalToggle } = $props();
	let selectedDate = $state({}) as DayDate;
	let selectedType = $state("") as string;

	export function init(date: DayDate, suspensionType: string) {
		selectedDate = date;
		selectedType = suspensionType;
	}

	interface ApiBody {
		Date: DayDate
		Type: string
	}

	async function removeSuspended() {
		let body: ApiBody = {
			Date: selectedDate,
			Type: selectedType
		}
		await modFetch("/api/removesuspended", {method: "DELETE", body: JSON.stringify(body)});
		modalToggle();
	}
</script>

<Modal isOpen={isModalOpen} toggle={modalToggle} header="Remove Event">
	<ModalBody>
		Are you sure to remove the suspension <span class="fw-bold">{selectedType}</span> on <span class="fw-bold">{selectedDate.Month}/{selectedDate.Day}/{selectedDate.Year}</span>? This action is irreversible.
	</ModalBody>
	<ModalFooter>
		<Button outline color="secondary" on:click={modalToggle}>Back</Button>
		<Button outline color="danger" on:click={removeSuspended}>Remove</Button>
	</ModalFooter>
</Modal>
