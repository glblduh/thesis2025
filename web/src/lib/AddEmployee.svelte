<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import { Button, Form, FormGroup, Input, InputGroup, Modal, ModalBody } from "@sveltestrap/sveltestrap";

	let { refreshList, isModalOpen, modalToggle } = $props();

	let formValidated = $state(false);
	let idNumber: number | undefined = $state();
	let employeeType: string | undefined = $state();
	let firstName: string | undefined = $state();
	let middleName: string | undefined = $state();
	let lastName: string | undefined = $state();

	function toggleValidate() {
		formValidated = !formValidated;
	}

	function clearVars() {
		formValidated = false;
		idNumber = undefined;
		employeeType = undefined;
		firstName = undefined;
		middleName = undefined;
		lastName = undefined;
		modalToggle();
	}

	function isFaculty(): boolean {
		let returnBool: boolean = false;
		switch(employeeType) {
			case "Faculty":
				returnBool = true;
				break;
			case "Staff":
				returnBool = false;
				break;
		}
		return returnBool;
	}

	interface addEmployeeBody {
		IdNumber: number | undefined;
		IsFaculty: boolean | undefined;
		FirstName: string | undefined;
		MiddleName: string | undefined;
		LastName: string | undefined;
	}

	async function addEmployeeToDB(e: SubmitEvent) {
		e.preventDefault();

		let jsonBody: addEmployeeBody = {
			IdNumber: idNumber,
			IsFaculty: isFaculty(),
			FirstName: firstName,
			MiddleName: middleName,
			LastName: lastName
		};

		let addEmployeeRes = await fetch("/api/addemployee", {method: "POST", body: JSON.stringify(jsonBody)})
		refreshList();
		clearVars();
	}
</script>

<Modal body autoFocus isOpen={isModalOpen} toggle={clearVars} header="Add Employee">
	<ModalBody>
		<Form validated={formValidated} on:submit={addEmployeeToDB}>
			<InputGroup>
				<FormGroup floating label="ID Number">
					<Input placeholder="ID Number" type="number" required bind:value={idNumber} />
				</FormGroup>
				<FormGroup floating label="Employee Type">
					<Input placeholder="Employee Type" type="select" bind:value={employeeType}>
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
				<Button type="submit" color="primary" on:click={toggleValidate}>Add</Button>
			</div>
		</Form>
	</ModalBody>
</Modal>
