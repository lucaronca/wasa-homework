<script>
export default {
	props: {
		userId: String,
		isFollowedByLoggedInUser: Boolean,
	},
	emits: ['followUpdated'],
	methods: {
		async follow(willFollow) {
			try {
				if (willFollow) {
					await this.$axios.put(
						`/users/me/followings/${this.userId}`
					);
					this.$emit('followUpdated', true)
				} else {
					await this.$axios.delete(
						`/users/me/followings/${this.userId}`
					);
					this.$emit('followUpdated', false)
				}
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
	v-if="isFollowedByLoggedInUser"
	type="button"
	class="btn btn-sm btn-primary"
	@click.prevent="follow(false)"
>
	<svg class="feather">
		<use href="/feather-sprite-v4.29.0.svg#user-minus" />
	</svg>
	Unfollow user
</button>
<button
	v-else
	type="button"
	class="btn btn-sm btn-primary"
	@click.prevent="follow(true)"
>
	<svg class="feather">
		<use href="/feather-sprite-v4.29.0.svg#user-plus" />
	</svg>
	Follow user
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
