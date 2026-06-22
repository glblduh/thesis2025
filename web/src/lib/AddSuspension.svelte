<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, Input, ModalFooter, Form, Icon, FormGroup, InputGroup, Row, Col } from "@sveltestrap/sveltestrap";
    import type { SuspendedDay, DayDate } from "./utils";
    import RemoveSuspension from "./RemoveSuspension.svelte";

	let { isModalOpen, modalToggle } = $props();
	let suspensions = $state([]) as SuspendedDay[];
	let inputDateString = $state("");
	let inputDateFormatted = $derived.by(() => {
		let splittedDate = inputDateString.split("-");
		return {
			Year: Number(splittedDate[0]),
			Month: Number(splittedDate[1]),
			Day: Number(splittedDate[2])
		};
	});
	let inputType = $state("SUSPENSION");
	let formValidated = $state(false);

	function clearVars() {
		suspensions = [];
		inputDateString = "";
		inputType = "SUSPENSION";
		formValidated = false;
		modalToggle();
	}

	export async function init() {
		let getAllSuspended = await fetch("/api/getallsuspended");
		suspensions = await getAllSuspended.json();
	}

	function toggleValidate() {
		formValidated = !formValidated;
	}

	interface ApiBody {
		Date: DayDate,
		Type: string
	}

	async function updateSuspended(e: Event) {
		e.preventDefault();

		let body: ApiBody = {
			Date: inputDateFormatted,
			Type: inputType
		};
		await fetch("/api/updatesuspended", {method: "POST", body: JSON.stringify(body)});
		init();
	}

	let removeSuspensionModal: RemoveSuspension;
	let removeSuspensionModalState = $state(false);
	function removeSuspensionModalToggle() {
		clearVars();
		removeSuspensionModalState = !removeSuspensionModalState;
	}
</script>

<RemoveSuspension bind:this={removeSuspensionModal} isModalOpen={removeSuspensionModalState} modalToggle={removeSuspensionModalToggle} />

<Modal isOpen={isModalOpen} toggle={clearVars} header="Suspensions">
	<ModalBody>
		<Form validated={formValidated} on:submit={updateSuspended}>
			<InputGroup>
				<FormGroup floating label="Date">
					<Input type="date" required bind:value={inputDateString} placeholder="date placeholder"/>
				</FormGroup>
				<FormGroup floating label="Suspension Type">
					<Input type="select" bind:value={inputType}>
						<option value="SUSPENSION">SUSPENSION</option>
						<option value="HOLIDAY">HOLIDAY</option>
					</Input>
				</FormGroup>
			</InputGroup>
			<div class="text-end">
				<Button type="submit" color="info" on:click={toggleValidate}><Icon name="plus-lg" class="fw-bold"/></Button>
			</div>
		</Form>
		{#if suspensions.length != 0}
			<Table responsive striped>
				<thead>
					<tr>
						<th scope="col">YEAR</th>
						<th scope="col">MONTH</th>
						<th scope="col">DAY</th>
						<th scope="col">TYPE</th>
						<th scope="col">REMOVE?</th>
					</tr>
					{#each suspensions as suspension }
						<tr>
							<td>{suspension.Date.Year}</td>
							<td>{suspension.Date.Month}</td>
							<td>{suspension.Date.Day}</td>
							<td>{suspension.Type}</td>
							<td><Button color="danger" on:click={() => {
								removeSuspensionModal.init(suspension.Date);
								removeSuspensionModalToggle();
							}}>REMOVE</Button></td>
						</tr>
					{/each}
				</thead>
			</Table>
		{/if}
	</ModalBody>
</Modal>
