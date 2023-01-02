<script>
import { loggedInUser } from '../user'

export default {
	mounted() {
		window.localStorage.clear()
	},
	data() {
		return {
            username: '',
        }
	},
	methods: {
		async login() {
			try {
				const response = await this.$axios.post('/session', {
					name: this.username.trim(),
				})
				const { identifier } = response.data
				window.localStorage.setItem('$loggedInUserToken', identifier)

				const resp = await this.$axios.get('/users/me')
				loggedInUser.id = resp.data.id;
				loggedInUser.username = resp.data.username;
				window.document.title = `Wasa Photo - ${loggedInUser.username}`
				this.$router.replace('/')
			} catch (e) {
				console.error(e);
				if (e.response?.data) {
					window.alert(e.response.data);
				} else {
					window.alert("Api error");
				}
			}
		},
	},
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
		>
			<h1 class="h2">Login</h1>
		</div>
		<div class="alert alert-info" role="alert">
			Write a username to log in as that user, no password is required.<br />
			If you login with a new username, a newly user will be created.
		</div>
		<div class="input-group mb-3">
			<div class="input-group-prepend">
				<span class="input-group-text"
					>Username</span
				>
			</div>
			<input
				v-model="username"
				type="text"
				class="form-control"
				placeholder="Write your username"
			/>
			<button
					type="button"
					class="btn btn-primary"
					@click.prevent="login"
				>
					Login
			</button>
		</div>
	</div>
</template>

<style scoped>
.input-group, .alert {
	max-width: 600px;
}
.input-group-text {
	border-radius: 0;
}

</style>
