<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, Table, Modal, ModalBody, Input, ModalFooter, Form, Icon, FormGroup, InputGroup, Row, Col } from "@sveltestrap/sveltestrap";
    import type { SuspendedDay } from "./utils";

	let { isModalOpen, modalToggle } = $props();
	let suspensions = $state([]) as SuspendedDay[];

	export async function init() {
		let getAllSuspended = await fetch("/api/getallsuspended");
		suspensions = await getAllSuspended.json();
	}

	async function updateSuspended() {
		await fetch("/api/updatesuspended")
	}
</script>

<Modal isOpen={isModalOpen} toggle={modalToggle} header="Suspensions">
	<ModalBody>
		<Form>
			<InputGroup>
				<FormGroup floating label="Date">
					<Input type="date" placeholder="date placeholder"/>
				</FormGroup>
				<FormGroup floating label="Suspension Type">
					<Input type="select">
						<option>GOVT SUSPENSION</option>
						<option>HOLIDAY</option>
					</Input>
				</FormGroup>
				<Button color="success"><Icon name="plus-lg" class="fw-bold"/></Button>
			</InputGroup>
		</Form>
		<Table responsive striped>
			<thead>
				<tr>
					<th scope="col">MONTH</th>
					<th scope="col">DAY</th>
					<th scope="col">YEAR</th>
					<th scope="col">TYPE</th>
				</tr>
				{#each suspensions as suspension }
					<tr>
						<td>{suspension.Date.Month}</td>
						<td>{suspension.Date.Day}</td>
						<td>{suspension.Date.Year}</td>
						<td>{suspension.Type}</td>
					</tr>
				{/each}
			</thead>
		</Table>
	</ModalBody>
</Modal>
