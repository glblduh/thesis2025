<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import { Button, ButtonGroup, Modal, Table } from "@sveltestrap/sveltestrap";
	import { onMount } from "svelte";
	import AddEmployee from "./lib/AddEmployee.svelte";
	import RemoveEmployee from "./lib/RemoveEmployee.svelte";
	import Attendances from "./lib/Attendances.svelte";
	import Schedules from "./lib/Schedules.svelte";

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

	onMount(async () => {
		parseEmployees();
	});

	let employees: Employee[] = [];

	let addEmployeeModalState = false;
	function addEmployeeModalToggle() {
		addEmployeeModalState = !addEmployeeModalState;
	}

	let removeEmployeeModalState = false;
	function removeEmployeeModalToggle() {
		removeEmployeeModalState = !removeEmployeeModalState;
	}

	let employeeAttendancesModalState = false;
	function employeeAttendancesModalToggle() {
		employeeAttendancesModalState = !employeeAttendancesModalState;
	}

	let employeeSchedulesModalState = false;
	function employeeSchedulesModalToggle() {
		employeeSchedulesModalState = !employeeSchedulesModalState;
	}
</script>

<main>
	<AddEmployee isModalOpen={addEmployeeModalState} modalToggle={addEmployeeModalToggle} />

	<div class="header">
		<h2>Attendance Viewer</h2>
	</div>

	<div>
		<ButtonGroup size="sm">
			<Button color="success" on:click={addEmployeeModalToggle}>Add Employee</Button>
			<Button color="info">Refresh</Button>
		</ButtonGroup>
	</div>

	<div class="employeeContainer">
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
								<Button color="primary">Attendances</Button>
								<Button color="primary">Schedules</Button>
								<Button color="danger">Remove</Button>
							</ButtonGroup>
						</td>
						<td>{employee.idNumber}</td>
						<td>{employee.employeeType}</td>
						<td>{employee.firstName}</td>
						<td>{employee.middleName}</td>
						<td>{employee.lastName}</td>
					</tr>
				{/each}
				<tr>
					<td>
						<ButtonGroup vertical size="sm">
							<Button color="primary">Attendances</Button>
							<Button color="primary">Schedules</Button>
							<Button color="danger">Remove</Button>
						</ButtonGroup>
					</td>
					<td>1</td>
					<td>gg</td>
					<td>gg</td>
					<td>gg</td>
					<td>gg</td>
				</tr>
			</tbody>
		</Table>
	</div>
</main>

<style>
	.employeeContainer {
		display: grid;
		grid-template-columns: 1;
		padding: 1% 1% 1% 1%;
	}
</style>
