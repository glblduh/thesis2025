<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import 'bootstrap-icons/font/bootstrap-icons.css';
	import { Button, ButtonGroup, Table, Navbar, NavbarBrand, Icon, Input, Container, Row, Col, Dropdown, DropdownToggle, DropdownItem, DropdownMenu } from "@sveltestrap/sveltestrap";
	import { onMount, untrack } from "svelte";
	import AddEmployee from "./lib/AddEmployee.svelte";
	import RemoveEmployee from "./lib/RemoveEmployee.svelte";
	import Attendances from "./lib/Attendances.svelte";
	import Schedules from "./lib/Schedules.svelte";
	import Suspension from "./lib/AddSuspension.svelte"
    import AddSuspension from "./lib/AddSuspension.svelte";
    import MonthAttendances from "./lib/MonthAttendances.svelte";
    import LoginAuth from "./lib/LoginAuth.svelte";
    import { modFetch, logout, isAuthenticated, getUsername, getUserType } from "./lib/utils";

	let employees: Employee[] = $state([]);
	let selectedEmployee: number = $state(0);
	let search = $state("");
	let searchedEmployees: Employee[] = $derived(employees.filter(employee => {
		const employeeDetails = employee.IdNumber.toString() + employee.FirstName + employee.MiddleName + employee.LastName;
		return employeeDetails.toLocaleLowerCase().includes(search.toLocaleLowerCase());
	}));
	let authenticated = $state(false);
	$effect(() => {
		if (authenticated) {
			untrack(parseEmployees);
		}
	});

	interface Employee {
		IdNumber: number;
		LastName: string;
		FirstName: string;
		MiddleName: string;
		EmployeeType: string;
	}

	interface apiRes {
		Faculty: Employee[]
		Staff: Employee[]
	}

	async function parseEmployees() {
		const res = await modFetch("/api/getallemployees");
		const jsonRes: apiRes = await res.json();

		employees = [...jsonRes.Faculty || [], ...jsonRes.Staff || []]
	}

	function selectEmployee(selected: number) {
		selectedEmployee = selected;
	}

	onMount(async () => {
		authenticated = isAuthenticated();

		if (!authenticated) {
			loginModalToggle();
		}
	});

	let addEmployeeModalState = $state(false);
	function addEmployeeModalToggle() {
		addEmployeeModalState = !addEmployeeModalState;
	}

	let removeEmployeeModalState = $state(false);
	function removeEmployeeModalToggle() {
		removeEmployeeModalState = !removeEmployeeModalState;
	}

	let employeeAttendancesModal: Attendances = $state() as Attendances;
	let employeeAttendancesModalState = $state(false);
	function employeeAttendancesModalToggle() {
		if (!employeeAttendancesModalState) {
			employeeAttendancesModal.init(selectedEmployee);
		}
		employeeAttendancesModalState = !employeeAttendancesModalState;
	}

	let employeeSchedulesModal: Schedules = $state() as Schedules;
	let employeeSchedulesModalState = $state(false);
	function employeeSchedulesModalToggle() {
		if (!employeeSchedulesModalState) {
			employeeSchedulesModal.init(selectedEmployee);
		}
		employeeSchedulesModalState = !employeeSchedulesModalState;
	}

	let suspensionModal: Suspension = $state() as Suspension;
	let suspensionModalState = $state(false);
	function suspensionModalToggle() {
		if (!suspensionModalState) {
			suspensionModal.init();
		}
		suspensionModalState = !suspensionModalState;
	}

	let monthAttendancesModal: MonthAttendances = $state() as MonthAttendances;
	let monthAttendancesModalState = $state(false);
	function monthAttendancesModalToggle() {
		if (!monthAttendancesModalState) {
			monthAttendancesModal.getSchoolYears();
		}
		monthAttendancesModalState = !monthAttendancesModalState
	}

	let loginModal: LoginAuth = $state() as LoginAuth;
	let loginModalState = $state(false)
	function loginModalToggle() {
		if (loginModalState && !authenticated) {
			return
		}
		loginModalState = !loginModalState
	}
</script>

<main class="px-4">
	<LoginAuth isModalOpen={loginModalState} modalToggle={loginModalToggle} />

	{#if authenticated}
		<AddEmployee isModalOpen={addEmployeeModalState} modalToggle={addEmployeeModalToggle} refreshList={parseEmployees} />
		<RemoveEmployee isModalOpen={removeEmployeeModalState} modalToggle={removeEmployeeModalToggle} refreshList={parseEmployees} idNumber={selectedEmployee} />
		<Schedules bind:this={employeeSchedulesModal} isModalOpen={employeeSchedulesModalState} modalToggle={employeeSchedulesModalToggle} />
		<Attendances bind:this={employeeAttendancesModal} isModalOpen={employeeAttendancesModalState} modalToggle={employeeAttendancesModalToggle} />
		<AddSuspension bind:this={suspensionModal} isModalOpen={suspensionModalState} modalToggle={suspensionModalToggle} />
		<MonthAttendances bind:this={monthAttendancesModal} isModalOpen={monthAttendancesModalState} modalToggle={monthAttendancesModalToggle} />

		<Navbar fixed="top" sticky="top">
			<NavbarBrand href="/" class="fw-bold">Attendance Viewer</NavbarBrand>
			<Dropdown>
				<DropdownToggle outline color="primary" size="sm" class="shadow">{getUsername()}</DropdownToggle>
				<DropdownMenu>
					<DropdownItem disabled>{getUserType()}</DropdownItem>
					<DropdownItem divider />
					<DropdownItem on:click={logout} class="text-danger"><Icon name="box-arrow-left" /> Logout</DropdownItem>
				</DropdownMenu>
			</Dropdown>
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
								<td>{employee.IdNumber}</td>
								<td>{employee.EmployeeType}</td>
								<td>{employee.FirstName}</td>
								<td>{employee.MiddleName}</td>
								<td>{employee.LastName}</td>
								<td style="width: 1%;">
									<Dropdown style="width: 100%;">
										<DropdownToggle outline color="primary" size="sm"><Icon name="three-dots" class="fw-bold" /></DropdownToggle>
										<DropdownMenu>
											<DropdownItem on:click={() => {selectEmployee(employee.IdNumber); employeeAttendancesModalToggle();}} ><Icon name="card-list" /> Attendances</DropdownItem>
											<DropdownItem on:click={() => {selectEmployee(employee.IdNumber); employeeSchedulesModalToggle();}}><Icon name="calendar" /> Schedules</DropdownItem>
											<DropdownItem divider />
											<DropdownItem on:click={() => {selectEmployee(employee.IdNumber); removeEmployeeModalToggle();}} class="text-danger"><Icon name="trash" /> Remove</DropdownItem>
										</DropdownMenu>
									</Dropdown>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</Table>
		</div>
	{/if}
</main>
