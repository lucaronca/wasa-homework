<script>
export default {
	props: {
		photoId: Number,
		image: String,
	},
	emits: ['commentsUpdated'],
	data() {
		return {
			showCommentModal: false,
			comment: '',
		};
	},
	watch: {
		showCommentModal(newShowCommentModal) {
			if (newShowCommentModal) {
				document.body.style.overflow = "hidden";
			} else {
				document.body.style.overflow = "auto";
			}
		},
	},
	methods: {
		async saveComment() {
			try {
				await this.$axios.post(`/photos/${this.photoId}/comments`, {
					content: this.comment,
				});
				this.$emit("commentsUpdated");
				this.showCommentModal = false;
				this.comment = ""
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
};
</script>
<template>
	<button
		type="button"
		class="btn btn-outline-primary btn-sm action-btn"
		@click.prevent="showCommentModal = !showCommentModal"
	>
		<span>
			Comment
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#message-circle" />
			</svg>
		</span>
	</button>
	<div v-if="showCommentModal" class="modal-backdrop fade show"></div>
	<div
		class="modal fade"
		:class="{ show: showCommentModal }"
		tabindex="-1"
		role="dialog"
	>
		<div class="modal-dialog modal-dialog-centered" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="exampleModalLongTitle">
						Comment a photo
					</h5>
					<span
						class="close"
						@click.prevent="showCommentModal = !showCommentModal"
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
					<div class="form-group">
						<textarea
							v-model="comment"
							class="form-control"
							placeholder="Write a comment here"
							rows="3"
						></textarea>
					</div>
				</div>
				<div class="modal-footer">
					<button
						type="button"
						class="btn btn-secondary"
						@click.prevent="showCommentModal = !showCommentModal"
					>
						Close
					</button>
					<button
						type="button"
						class="btn btn-primary"
						@click.prevent="saveComment"
					>
						Comment
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
