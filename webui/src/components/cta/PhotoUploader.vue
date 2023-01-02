<script>
export default {
	emits: ['photoUploaded'],
	methods: {
		uploadPhoto() {
			const input = document.createElement("input");
			input.type = "file";
			input.onchange = async (e) => {
				const file = e.target.files[0];
				try {
					await this.$axios({
						method: "POST",
						url: "/photos",
						data: file,
					});
					this.$emit('photoUploaded')
				} catch (e) {
					console.error(e);
					if (e.response?.data) {
						window.alert(e.response.data);
					} else {
						window.alert("Api error");
					}
				} finally {
					input.remove();
				}
			};
			input.click();
		},
	}
}
</script>
<template>
	<button
		type="button"
		class="btn btn-sm btn-primary"
		@click.prevent="uploadPhoto"
	>
		<svg class="feather">
			<use href="/feather-sprite-v4.29.0.svg#image" />
		</svg>
		Upload a photo
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
