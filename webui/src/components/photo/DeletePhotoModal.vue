<script>
export default {
	props: {
		photoId: Number,
		image: String,
	},
	emits: ['photoDeleted'],
	data() {
		return {
			showDeleteModal: false,
		}
	},
	methods: {
		async deletePhoto() {
			try {
				await this.$axios.delete(`/photos/${this.photoId}`);
				this.showDeleteModal = false;
				this.$emit("photoDeleted");
			} catch (e) {
				console.error(e);
			}
		},
	}
}
</script>
<template>
<button
	type="button"
	class="btn btn-outline-primary btn-sm action-btn"
	@click.prevent="showDeleteModal = !showDeleteModal"
>
	<span>
		Delete
		<svg class="feather">
			<use href="/feather-sprite-v4.29.0.svg#trash-2" />
		</svg>
	</span>
</button>
<div v-if="showDeleteModal" class="modal-backdrop fade show"></div>
	<div
		class="modal fade"
		:class="{ show: showDeleteModal }"
		tabindex="-1"
		role="dialog"
	>
		<div class="modal-dialog modal-dialog-centered" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="exampleModalLongTitle">
						Delete a photo
					</h5>
					<span
						class="close"
						@click.prevent="showDeleteModal = !showDeleteModal"
					>
						<svg class="feather" stroke="#000000">
							<use href="/feather-sprite-v4.29.0.svg#x" />
						</svg>
					</span>
				</div>
				<div class="modal-body">
					<img
						:src="image"
						class="rounded mx-auto d-block img-thumbnail mb-3"
						alt="Image to comment"
					/>
					Are you sure you want to delete this photo?
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						@click.prevent="showDeleteModal = !showDeleteModal"
					>
						Close
					</button>
					<button
						type="button"
						class="btn btn-primary"
						@click.prevent="deletePhoto"
					>
						Delete
					</button>
				</div>
			</div>
		</div>
	</div>
</template>
<style scoped>
.action-btn {
	width: 50%;
	border-radius: 0;
}
.modal.show {
	display: block;
}
.modal:not(.show) {
	display: none;
}
.close {
	cursor: pointer;
}
</style>
