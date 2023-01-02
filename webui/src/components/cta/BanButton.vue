<script>
export default {
	props: {
		userId: String,
		bannedForUser: Boolean,
	},
	emits: ['banUpdated'],
	methods: {
		async ban(willBan) {
			try {
				if (willBan) {
					await this.$axios.put(
						`/users/me/bans/${this.userId}`
					);
				} else {
					await this.$axios.delete(
						`/users/me/bans/${this.userId}`
					);
				}
				this.$emit('banUpdated')
			} catch (e) {
				if (e.response?.data) {
					window.alert(e.response.data);
				} else {
					window.alert("Api error");
				}
			}
		},
	}
}
</script>
<template>
	<button
		v-if="bannedForUser"
		type="button"
		class="btn btn-sm btn-outline-secondary"
		@click.prevent="ban(false)"
	>
		<svg class="feather">
			<use href="/feather-sprite-v4.29.0.svg#user-check" />
		</svg>
		Unban user
	</button>
	<button
		v-else
		type="button"
		class="btn btn-sm btn-outline-secondary"
		@click.prevent="ban(true)"
	>
		<svg class="feather">
			<use href="/feather-sprite-v4.29.0.svg#user-x" />
		</svg>
		Ban user
	</button>
</template>
<style scoped>
button.btn {
	display: flex;
	align-items: center;
}
button.btn > svg {
	margin-right: 0.5rem;
}
</style>
