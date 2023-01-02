<script>
export default {
	data() {
		return {
			username: '',
			loading: false,
			users: [],
		};
	},
	methods: {
		async search(showAll) {
			this.loading = true;
			try {
				let response = await this.$axios.get("/users", {
					params: {
						...(!showAll && {
							username: this.username,
						}),
					},
				});
				this.users = response.data;
			} catch (e) {
				console.error(e);
				if (e.response?.data) {
					window.alert(e.response.data);
				} else {
					window.alert("Api error");
				}
			}
			this.loading = false;
		},
	},
};
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
		>
			<h1 class="h2">Search users</h1>
		</div>
		<div class="container">
			<div class="row gx-3 gy-3">
				<div class="col-md-6">
					<button
						type="button"
						class="btn btn-secondary mr-0"
						@click.prevent="search(true)"
					>
						All users
					</button>
				</div>
				<div class="col-md-6">
					<div class="input-group mb-3 ml-auto">
						<div class="input-group-prepend">
							<span class="input-group-text">Username</span>
						</div>
						<input
							v-model="username"
							type="text"
							class="form-control"
							placeholder="Search by username"
						/>
						<button
							type="button"
							class="btn btn-primary"
							@click.prevent="search(false)"
						>
							Search
						</button>
					</div>
				</div>
			</div>
			<div class="row gx-3 gy-3">
				<table class="table">
					<thead>
						<tr>
							<th scope="col">Id #</th>
							<th scope="col">Profile link</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="user in users" :key="user.id">
							<th scope="row">{{ user.id }}</th>
							<td><ProfileLink :id="user.id" :name="user.username" /></td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<LoadingSpinner :loading="loading" />
	</div>
</template>

<style scoped>
.container {
	padding: 0;
}
.input-group {
	max-width: 500px;
	margin-left: auto;
}
.input-group-text,
button.btn {
	border-radius: 0;
}
</style>
