<script>
export default {
	emits: ['usernameUpdated'],
	data() {
		return {
			showUsernameModal: false,
			username: "",
		};
	},
	methods: {
		async updateUsername() {
			try {
				await this.$axios({
					method: "PATCH",
					url: "/users/me",
					data: [
						{
							op: "replace",
							path: "/username",
							value: this.username.trim(),
						},
					],
					headers: {
						"Content-Type": "application/json-patch+json",
					},
				});
				this.showUsernameModal = false;
				this.$emit('usernameUpdated')
			} catch (e) {
				console.error(e);
				if (e.response?.data) {
					window.alert(
						e.response.data === "Payload not valid"
							? "Username not valid"
							: e.response.data
					);
				}
			}
		},
	}
}
</script>
<template>
	<button
		type="button"
		class="btn btn-sm btn-outline-secondary ch-username"
		@click="showUsernameModal = !showUsernameModal"
	>
		Change username
	</button>
	<div v-if="showUsernameModal" class="modal-backdrop fade show"></div>
	<div
		class="modal fade"
		:class="{ show: showUsernameModal }"
		tabindex="-1"
		role="dialog"
	>
		<div class="modal-dialog modal-dialog-centered" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="exampleModalLongTitle">
						Set a new username
					</h5>
					<span
						class="close"
						@click.prevent="showUsernameModal = !showUsernameModal"
					>
						<svg class="feather" stroke="#000000">
							<use href="/feather-sprite-v4.29.0.svg#x" />
						</svg>
					</span>
				</div>
				<div class="modal-body">
					<div class="input-group mb-3">
						<div class="input-group-prepend">
							<span class="input-group-text">Username</span>
						</div>
						<input
							v-model="username"
							type="text"
							class="form-control"
							placeholder="New username"
						/>
					</div>
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						@click.prevent="showUsernameModal = !showUsernameModal"
					>
						Close
					</button>
					<button
						type="button"
						class="btn btn-primary"
						@click.prevent="updateUsername"
					>
						Save
					</button>
				</div>
			</div>
		</div>
	</div>
</template>
<style scoped>
.modal.show {
	display: block;
}
.modal:not(.show) {
	display: none;
}
.close {
	cursor: pointer;
}
.input-group-prepend > span {
	border-radius: 0;
}
.btn.ch-username.btn.ch-username {
	border-top-right-radius: 4px;
	border-bottom-right-radius: 4px;
}
</style>
