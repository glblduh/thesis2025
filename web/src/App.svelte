<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, ButtonGroup, Table, Navbar, NavbarBrand, Icon, Input, Container, Row, Col, Dropdown, DropdownToggle, DropdownItem, DropdownMenu } from "@sveltestrap/sveltestrap";
	import { onMount } from "svelte";
	import AddEmployee from "./lib/AddEmployee.svelte";
	import RemoveEmployee from "./lib/RemoveEmployee.svelte";
	import Attendances from "./lib/Attendances.svelte";
	import Schedules from "./lib/Schedules.svelte";
	import Suspension from "./lib/AddSuspension.svelte"
    import AddSuspension from "./lib/AddSuspension.svelte";
    import MonthAttendances from "./lib/MonthAttendances.svelte";

	let employees: Employee[] = $state([]);
	let selectedEmployee: number = $state(0);
	let search = $state("");
	let searchedEmployees: Employee[] = $derived(employees.filter(employee =>
		employee.idNumber.toString().toLowerCase().includes(search.toLowerCase()) ||
		employee.firstName.toLowerCase().includes(search.toLowerCase()) ||
		employee.middleName.toLowerCase().includes(search.toLowerCase()) ||
		employee.lastName.toLowerCase().includes(search.toLowerCase())
	));

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

	let addEmployeeModalState = $state(false);
	function addEmployeeModalToggle() {
		addEmployeeModalState = !addEmployeeModalState;
	}

	let removeEmployeeModalState = $state(false);
	function removeEmployeeModalToggle() {
		removeEmployeeModalState = !removeEmployeeModalState;
	}

	let employeeAttendancesModal: Attendances;
	let employeeAttendancesModalState = $state(false);
	function employeeAttendancesModalToggle() {
		if (!employeeAttendancesModalState) {
			employeeAttendancesModal.init(selectedEmployee);
		}
		employeeAttendancesModalState = !employeeAttendancesModalState;
	}

	let employeeSchedulesModal: Schedules;
	let employeeSchedulesModalState = $state(false);
	function employeeSchedulesModalToggle() {
		if (!employeeSchedulesModalState) {
			employeeSchedulesModal.init(selectedEmployee);
		}
		employeeSchedulesModalState = !employeeSchedulesModalState;
	}

	let suspensionModal: Suspension;
	let suspensionModalState = $state(false);
	function suspensionModalToggle() {
		if (!suspensionModalState) {
			suspensionModal.init();
		}
		suspensionModalState = !suspensionModalState;
	}

	let monthAttendancesModal: MonthAttendances;
	let monthAttendancesModalState = $state(false);
	function monthAttendancesModalToggle() {
		if (!monthAttendancesModalState) {
			monthAttendancesModal.getSchoolYears();
		}
		monthAttendancesModalState = !monthAttendancesModalState
	}
</script>

<main class="px-4">
	<AddEmployee isModalOpen={addEmployeeModalState} modalToggle={addEmployeeModalToggle} refreshList={parseEmployees} />
	<RemoveEmployee isModalOpen={removeEmployeeModalState} modalToggle={removeEmployeeModalToggle} refreshList={parseEmployees} idNumber={selectedEmployee} />
	<Schedules bind:this={employeeSchedulesModal} isModalOpen={employeeSchedulesModalState} modalToggle={employeeSchedulesModalToggle} />
	<Attendances bind:this={employeeAttendancesModal} isModalOpen={employeeAttendancesModalState} modalToggle={employeeAttendancesModalToggle} />
	<AddSuspension bind:this={suspensionModal} isModalOpen={suspensionModalState} modalToggle={suspensionModalToggle} />
	<MonthAttendances bind:this={monthAttendancesModal} isModalOpen={monthAttendancesModalState} modalToggle={monthAttendancesModalToggle} />

	<Navbar fixed="top" sticky="top">
		<NavbarBrand href="/" class="fw-bold">Attendance Viewer</NavbarBrand>
	</Navbar>

	{#if searchedEmployees.length != 0 || search.length != 0}
		<Container fluid class="px-0">
			<Row cols={2}>
				<Col xs={12} md={6}>
					<Input type="text" placeholder="Search" bind:value={search} class="shadow" />
				</Col>
				<Col xs={12} md={6} class="mt-1">
					<div class="h-100 d-flex align-items-center justify-content-center justify-content-md-end">
						<ButtonGroup size="sm" class="shadow">
							<Button outline color="primary" on:click={monthAttendancesModalToggle}><Icon name="list-columns-reverse" class="fw-bold" /> All Attendances</Button>
							<Button outline color="primary" on:click={addEmployeeModalToggle}><Icon name="person-plus" class="fw-bold" /> Add Employee</Button>
							<Button outline color="primary" on:click={suspensionModalToggle}><Icon name="calendar-plus" class="fw-bold" /> Add Suspension</Button>
							<Button outline color="primary" on:click={parseEmployees}><Icon name="arrow-clockwise" class="fw-bold" /> Refresh</Button>
						</ButtonGroup>
					</div>
				</Col>
			</Row>
		</Container>
	{/if}

	<div class="mt-4 shadow rounded overflow-hidden">
		<Table responsive hover class="mb-0">
			<thead>
				<tr>
					{#if searchedEmployees.length != 0}
						<th scope="col">ID NUMBER</th>
						<th scope="col">TYPE</th>
						<th scope="col">FIRST NAME</th>
						<th scope="col">MIDDLE NAME</th>
						<th scope="col">LAST NAME</th>
						<th scope="col"></th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#if searchedEmployees.length == 0}
					<tr>
						<td class="fw-bold text-center">No Employee Found</td>
					</tr>
				{:else}
					{#each searchedEmployees as employee}
						<tr>
							<td>{employee.idNumber}</td>
							<td>{employee.employeeType}</td>
							<td>{employee.firstName}</td>
							<td>{employee.middleName}</td>
							<td>{employee.lastName}</td>
							<td style="width: 1%;">
								<Dropdown style="width: 100%;">
									<DropdownToggle outline color="primary" size="sm"><Icon name="three-dots" class="fw-bold" /></DropdownToggle>
									<DropdownMenu>
										<DropdownItem on:click={() => {selectEmployee(employee.idNumber); employeeAttendancesModalToggle();}} ><Icon name="card-list" /> Attendances</DropdownItem>
										<DropdownItem on:click={() => {selectEmployee(employee.idNumber); employeeSchedulesModalToggle();}}><Icon name="calendar" /> Schedules</DropdownItem>
										<DropdownItem divider />
										<DropdownItem on:click={() => {selectEmployee(employee.idNumber); removeEmployeeModalToggle();}} class="text-danger"><Icon name="trash" /> Remove</DropdownItem>
									</DropdownMenu>
								</Dropdown>
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</Table>
	</div>
</main>
