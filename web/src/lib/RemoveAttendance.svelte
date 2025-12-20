<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Modal, ModalBody, ModalFooter } from "@sveltestrap/sveltestrap";
    import { monthsName, type DayDate } from "./utils";

	let { isModalOpen, modalToggle } = $props();
	let selectedEmployee: number = $state(0);
	let selectedDate = $state({}) as DayDate;

	export function init(idNumber: number, date: DayDate) {
		selectedEmployee = idNumber;
		selectedDate = date;
	}

	interface ApiBody {
		IdNumber: number,
		Date: DayDate
	}

	async function removeAttendance() {
		let body: ApiBody = {
			IdNumber: selectedEmployee,
			Date: selectedDate
		}
		await fetch("/api/removeattendance", {method: "DELETE", body: JSON.stringify(body)});
		modalToggle();
	}
</script>

<Modal isOpen={isModalOpen} toggle={modalToggle} header="Remove Attendance">
	<ModalBody>
		Are you sure to remove the attendance of <span class="fw-bold">{selectedEmployee}</span> for <span class="fw-bold">{monthsName[selectedDate.Month-1]} {selectedDate.Day}, {selectedDate.Year}</span>? This action is irreversible.
	</ModalBody>
	<ModalFooter>
		<Button color="secondary" on:click={modalToggle}>Back</Button>
		<Button color="danger" on:click={removeAttendance}>Remove</Button>
	</ModalFooter>
</Modal>
