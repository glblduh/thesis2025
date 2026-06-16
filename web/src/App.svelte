<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, ButtonGroup, Table, Navbar, NavbarBrand, Icon } from "@sveltestrap/sveltestrap";
	import { onMount } from "svelte";
	import AddEmployee from "./lib/AddEmployee.svelte";
	import RemoveEmployee from "./lib/RemoveEmployee.svelte";
	import Attendances from "./lib/Attendances.svelte";
	import Schedules from "./lib/Schedules.svelte";
	import Suspension from "./lib/AddSuspension.svelte"
    import AddSuspension from "./lib/AddSuspension.svelte";

	let employees: Employee[] = [];
	let selectedEmployee: number = 0;

	interface Employee {
		idNumber: number;
		lastName: string;
		firstName: string;
		middleName: string;
		employeeType: string;
	}

	async function parseEmployees() {
		let getAllEmployees = await fetch("/api/getallemployees");
		let jsonAllEmployees = await getAllEmployees.json();
		let facultyArray = jsonAllEmployees.Faculty;
		let staffArray = jsonAllEmployees.Staff;

		employees.length = 0;

		if (facultyArray != null) {
			for (let i = 0; i < facultyArray.length; i++) {
				let currentFaculty = jsonAllEmployees.Faculty[i];

				let parsedEmployee: Employee = {
					idNumber: currentFaculty.IdNumber,
					employeeType: currentFaculty.EmployeeType,
					firstName: currentFaculty.FirstName,
					middleName: currentFaculty.MiddleName,
					lastName: currentFaculty.LastName,
				};

				employees.push(parsedEmployee);
			}
		}

		if (staffArray != null) {
			for (let i = 0; i < staffArray.length; i++) {
				let currentFaculty = jsonAllEmployees.Staff[i];

				let parsedEmployee: Employee = {
					idNumber: currentFaculty.IdNumber,
					employeeType: currentFaculty.EmployeeType,
					firstName: currentFaculty.FirstName,
					middleName: currentFaculty.MiddleName,
					lastName: currentFaculty.LastName,
				};

				employees.push(parsedEmployee);
			}
		}
	}

	function selectEmployee(selected: number) {
		selectedEmployee = selected;
	}

	onMount(async () => {
		parseEmployees();
	});

	let addEmployeeModalState = false;
	function addEmployeeModalToggle() {
		addEmployeeModalState = !addEmployeeModalState;
	}

	let removeEmployeeModalState = false;
	function removeEmployeeModalToggle() {
		removeEmployeeModalState = !removeEmployeeModalState;
	}

	let employeeAttendancesModal: Attendances;
	let employeeAttendancesModalState = false;
	function employeeAttendancesModalToggle() {
		if (!employeeAttendancesModalState) {
			employeeAttendancesModal.init(selectedEmployee);
		}
		employeeAttendancesModalState = !employeeAttendancesModalState;
	}

	let employeeSchedulesModal: Schedules;
	let employeeSchedulesModalState = false;
	function employeeSchedulesModalToggle() {
		if (!employeeSchedulesModalState) {
			employeeSchedulesModal.init(selectedEmployee);
		}
		employeeSchedulesModalState = !employeeSchedulesModalState;
	}

	let suspensionModal: Suspension;
	let suspensionModalState = false;
	function suspensionModalToggle() {
		if (!suspensionModalState) {
			suspensionModal.init();
		}
		suspensionModalState = !suspensionModalState;
	}
</script>

<main>
	<AddEmployee isModalOpen={addEmployeeModalState} modalToggle={addEmployeeModalToggle} refreshList={parseEmployees} />
	<RemoveEmployee isModalOpen={removeEmployeeModalState} modalToggle={removeEmployeeModalToggle} refreshList={parseEmployees} idNumber={selectedEmployee} />
	<Schedules bind:this={employeeSchedulesModal} isModalOpen={employeeSchedulesModalState} modalToggle={employeeSchedulesModalToggle} />
	<Attendances bind:this={employeeAttendancesModal} isModalOpen={employeeAttendancesModalState} modalToggle={employeeAttendancesModalToggle} />
	<AddSuspension bind:this={suspensionModal} isModalOpen={suspensionModalState} modalToggle={suspensionModalToggle} />

	<Navbar fixed="top" sticky="top">
		<NavbarBrand href="/" class="fw-bold">Attendance Viewer</NavbarBrand>
		<ButtonGroup size="sm">
			<Button color="success" on:click={addEmployeeModalToggle}><Icon name="person-plus-fill" class="fw-bold" /> Add Employee</Button>
			<Button color="success" on:click={suspensionModalToggle}><Icon name="calendar-plus-fill" class="fw-bold" /> Add Suspension</Button>
			<Button color="info" on:click={parseEmployees}><Icon name="arrow-clockwise" class="fw-bold" /> Refresh</Button>
		</ButtonGroup>
	</Navbar>

	<div style="padding: .5%;">
		<Table responsive striped>
			<thead>
				<tr>
					<th scope="col">ACTION</th>
					<th scope="col">ID NUMBER</th>
					<th scope="col">TYPE</th>
					<th scope="col">FIRST NAME</th>
					<th scope="col">MIDDLE NAME</th>
					<th scope="col">LAST NAME</th>
				</tr>
			</thead>
			<tbody>
				{#each employees as employee}
					<tr>
						<td>
							<ButtonGroup vertical size="sm">
								<Button color="primary" on:click={() => {selectEmployee(employee.idNumber); employeeAttendancesModalToggle();}}>Attendances</Button>
								<Button color="primary" on:click={() => {selectEmployee(employee.idNumber); employeeSchedulesModalToggle();}}>Schedules</Button>
								<Button color="danger" on:click={() => {selectEmployee(employee.idNumber); removeEmployeeModalToggle();}}>Remove</Button>
							</ButtonGroup>
						</td>
						<td>{employee.idNumber}</td>
						<td>{employee.employeeType}</td>
						<td>{employee.firstName}</td>
						<td>{employee.middleName}</td>
						<td>{employee.lastName}</td>
					</tr>
				{/each}
			</tbody>
		</Table>
	</div>
</main>

<style>
</style>
