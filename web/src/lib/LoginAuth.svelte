<script lang="ts">
	import "bootstrap/dist/css/bootstrap.min.css";
	import "bootstrap/dist/js/bootstrap.bundle.min.js";
	import { Alert, Button, Form, FormGroup, Input, Modal, ModalBody } from "@sveltestrap/sveltestrap";

	let { isModalOpen, modalToggle } = $props();

	interface loginAuthBody {
		Username: string;
		Password: string;
	}

	interface loginAuthRes {
		Username: string;
		Type: string;
		Key: string;
	}

	let formValidated = $state(false);
	let login = $state({}) as loginAuthBody;
	let status = $state(0);

	function toggleValidate() {
		formValidated = !formValidated;
	}

	async function validateAuth(e: SubmitEvent) {
		e.preventDefault();
		status = 0;

		let res = await fetch("/auth/validate", {method: "POST", body: JSON.stringify(login)})
		status = res.status

		if (!res.ok) {
			return
		}

		let jsonRes: loginAuthRes = await res.json()
		sessionStorage.setItem("username", jsonRes.Username);
		sessionStorage.setItem("usertype", jsonRes.Type)
		sessionStorage.setItem("userkey", jsonRes.Key);

		window.location.reload();
	}
</script>

<Modal body autoFocus isOpen={isModalOpen} header="Login">
	<ModalBody>
		<Alert fade isOpen={status != 0 && status != 200} color="warning">
			{#if status == 404 }
				User not found
			{/if}
			{#if status == 401 }
				Invalid crendentials
			{/if}
		</Alert>
		<Form validated={formValidated} on:submit={validateAuth}>
			<FormGroup floating label="Username">
				<Input type="text" required bind:value={login.Username} />
			</FormGroup>
			<FormGroup floating label="Password">
				<Input type="password" required bind:value={login.Password} />
			</FormGroup>
			<div class="d-flex justify-content-center">
				<Button outline type="submit" color="primary" on:click={toggleValidate}>Login</Button>
			</div>
		</Form>
	</ModalBody>
</Modal>
