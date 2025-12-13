<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import { Button, ButtonGroup, Form, FormGroup, Input, InputGroup, Modal, ModalBody, ModalFooter } from "@sveltestrap/sveltestrap";

	export let isModalOpen: boolean | undefined;
	export let modalToggle: () => void | undefined;

	let formValidated = false;
	let idNumber: number;
	let firstName: string;
	let middleName: string;
	let lastName: string;

	async function addEmployeeToDB(e: SubmitEvent) {
		e.preventDefault();
	}
</script>

<Modal body autoFocus isOpen={isModalOpen} toggle={modalToggle} header="Add Employee">
	<ModalBody>
		<Form validated={formValidated} on:submit={addEmployeeToDB}>
			<InputGroup>
				<FormGroup floating label="ID Number">
					<Input placeholder="ID Number" type="number" required bind:value={idNumber} />
				</FormGroup>
				<FormGroup floating label="Employee Type">
					<Input placeholder="Employee Type" type="select">
						<option>Faculty</option>
						<option>Staff</option>
					</Input>
				</FormGroup>
			</InputGroup>
			<FormGroup floating label="First Name">
				<Input type="text" required bind:value={firstName} />
			</FormGroup>
			<FormGroup floating label="Middle Name">
				<Input type="text" required bind:value={middleName} />
			</FormGroup>
			<FormGroup floating label="Last Name">
				<Input type="text" required bind:value={lastName} />
			</FormGroup>
			<div class="text-end">
				<Button type="submit" color="primary" on:click={() => {formValidated=true;}}>Add</Button>
			</div>
		</Form>
	</ModalBody>
</Modal>
