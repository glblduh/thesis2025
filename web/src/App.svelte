<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableHead,
		TableHeadCell,
		Radio,
		TableSearch,
		TableBodyRow,
	} from "flowbite-svelte";
	import { onMount } from "svelte";

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
</script>

<main>
	<div class="header">
		<h2>Attendance Viewer</h2>
	</div>

	<div class="employeeSelection">
		<div class="employeeContainer">
			<Table class="table table-striped caption-top">
				<caption> Employee List </caption>
				<TableHead>
					<TableHeadCell>SELECT</TableHeadCell>
					<TableHeadCell>ID NUMBER</TableHeadCell>
					<TableHeadCell>LAST NAME</TableHeadCell>
					<TableHeadCell>FIRST NAME</TableHeadCell>
					<TableHeadCell>MIDDLE NAME</TableHeadCell>
					<TableHeadCell>TYPE</TableHeadCell>
				</TableHead>
				<TableBody>
					{#each employees as employee}
						<TableBodyRow>
							<TableBodyCell><Radio /></TableBodyCell>
							<TableBodyCell>{employee.idNumber}</TableBodyCell>
							<TableBodyCell>{employee.lastName}</TableBodyCell>
							<TableBodyCell>{employee.firstName}</TableBodyCell>
							<TableBodyCell>{employee.middleName}</TableBodyCell>
							<TableBodyCell>{employee.employeeType}</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		</div>
		<div class="selectionAction">
			<span class="selectionActionTitle">Select an action</span>
			<div class="selectionActionButtons">
				<button>Add Employee</button>
				<button>Show Schedule</button>
				<button>Show Attendance</button>
				<button>Remove Employee</button>
			</div>
		</div>
	</div>
</main>

<style>
	.employeeSelection {
		display: grid;
		grid-template-columns: 4fr 1fr;
	}
	.employeeContainer {
		display: grid;
		grid-template-columns: 1;
	}
	.selectionAction {
		display: grid;
		grid-template-columns: 1;
	}
	.selectionActionButtons {
		display: grid;
		grid-template-columns: 1;
		grid-gap: 5%;
	}
</style>
