<script>
import { RouterLink } from "vue-router";
import { format, formatDistance } from "date-fns";

import { loggedInUser } from "../../user";
import CommentModal from "./CommentModal.vue";
import DeletePhotoModal from "./DeletePhotoModal.vue";

export default {
	components: {
		CommentModal,
		DeletePhotoModal,
	},
	props: {
		id: Number,
		image: String,
		ownerUsername: String,
		ownerId: Number,
		totalLikes: Number,
		userLiked: Boolean,
		totalComments: Number,
		uploadDate: String,
	},
	emits: [
		"likeToggled",
		"likesCountUpdated",
		"commentsCountUpdated",
		"photoDeleted",
	],
	computed: {
		ownedByUser() {
			return this.ownerId === loggedInUser.id;
		},
	},
	data() {
		return {
			likesList: [],
			showLikesList: false,
			likesListLoading: false,
			commentsList: [],
			showCommentsList: false,
			commentsListLoading: false,
			imageLoading: false,
			loggedInUser,
		};
	},
	watch: {
		async showLikesList(newShowLikesList) {
			if (newShowLikesList) {
				this.likesListLoading = true;
				try {
					this.likesList = (
						await this.$axios.get(`/photos/${this.id}/likes`)
					).data;
				} catch (err) {
					console.error(err);
				} finally {
					this.likesListLoading = false;
				}
			}
		},
		showCommentsList(newShowCommentsList) {
			if (newShowCommentsList) {
				setTimeout(() => {
					document.querySelector(`#photo-${this.id} .comments`).scrollIntoView({ block: 'center' })
				})
				this.updateCommentsList();
			}
		},
		showDeletePhotoModal(newShowDeletePhotoModal) {
			if (newShowDeletePhotoModal) {
				document.body.style.overflow = "hidden";
			} else {
				document.body.style.overflow = "auto";
			}
		},
		image: {
			immediate: true,
			handler() {
				this.imageLoading = true
			}
		}
	},
	methods: {
		async toggleLike() {
			try {
				if (this.userLiked) {
					await this.$axios.delete(`/photos/${this.id}/likes/me`);
				} else {
					await this.$axios.put(`/photos/${this.id}/likes/me`);
				}
				this.$emit("likeToggled");

				const newLikesCount = (
					await this.$axios.get(`/photos/${this.id}/likes`)
				).data.length;
				this.$emit("likesCountUpdated", newLikesCount);
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async deleteComment(commentId) {
			try {
				await this.$axios.delete(
					`/photos/${this.id}/comments/${commentId}`
				);
				this.onCommentsListChanged();
			} catch (e) {
				console.error(e);
			}
		},
		async onCommentsListChanged() {
			const newCount = await this.updateCommentsList();
			this.showCommentsList = newCount > 0;
		},
		async updateCommentsList() {
			this.commentsListLoading = true;
			try {
				this.commentsList = (
					await this.$axios.get(`/photos/${this.id}/comments`)
				).data;

				const count = this.commentsList.length
				if (count !== this.totalComments) {
					this.$emit("commentsCountUpdated", count);
				}
				return count
			} catch (err) {
				console.error(err);
			} finally {
				this.commentsListLoading = false;
			}
		},
		formatDate(date, relative = false) {
			if (relative) {
				return formatDistance(new Date(date), new Date(), {
					includeSeconds: true,
					addSuffix: true,
				});
			}
			return format(new Date(date), "yy/MM/dd HH:mm");
		},
	},
};
</script>

<template>
	<div class="card" :id="`photo-${id}`">
		<div class="card-header">
			<div>
				<RouterLink
					class="badge badge-primary user-badge"
					:to="{ name: 'Profile', params: { id: ownerId } }"
					><svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#user" />
					</svg>
					{{ ownerUsername }}
				</RouterLink>
			</div>
			<div>
				<span
					class="badge badge-secondary"
					:class="{ 'dropdown-toggle': totalLikes !== 0 }"
					@click="
						totalLikes !== 0 && (showLikesList = !showLikesList)
					"
				>
					Likes {{ totalLikes }}
				</span>
				<div v-show="showLikesList" class="dropdown-menu">
					<LoadingSpinner :loading="likesListLoading" />
					<ProfileLink
						v-for="like in likesList"
						:key="like.owner.id"
						:id="like.owner.id"
						:name="like.owner.username"
					/>
				</div>
				<span
					class="badge badge-secondary"
					:class="{ 'dropdown-toggle': totalComments !== 0 }"
					@click="
						totalComments !== 0 &&
							(showCommentsList = !showCommentsList)
					"
				>
					Comments {{ totalComments }}
				</span>
			</div>
		</div>
		<div v-show="imageLoading" class="load-wrapper">
			<div class="activity"></div>
		</div>
		<img v-show="!imageLoading" class="card-img-top" :src="image" alt="Card image cap" @load="imageLoading = false" />
		<div class="btn-group" role="group">
			<button
				v-if="!ownedByUser"
				type="button"
				class="btn btn-outline-primary btn-sm action-btn"
				:class="{ active: userLiked }"
				@click.prevent="toggleLike"
			>
				Like
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#thumbs-up" />
				</svg>
			</button>
			<CommentModal
				:photoId="id"
				:image="image"
				@commentsUpdated="onCommentsListChanged"
			/>
			<DeletePhotoModal
				v-if="ownedByUser"
				:photoId="id"
				:image="image"
				@photoDeleted="$emit('photoDeleted')"
			/>
		</div>
		<div class="collapse comments" :class="{ show: showCommentsList }">
			<div class="card-body">
				<div v-if="commentsListLoading">
					<LoadingSpinner :loading="commentsListLoading" marginTop="1rem" />
				</div>
				<ul class="list-group">
					<li
						class="list-group-item"
						v-for="comment in commentsList"
						:key="comment.id"
					>
						<div class="media">
							<div class="media-body">
								<div>
									<b class="mt-0">{{
										comment.owner.username
									}}</b>
									-
									<small class="text-muted">{{
										formatDate(comment.date, true)
									}}</small>
									<span
										v-if="
											loggedInUser.id === comment.owner.id
										"
										class="close"
										@click="deleteComment(comment.id)"
									>
										<svg class="feather" stroke="#000000">
											<use
												href="/feather-sprite-v4.29.0.svg#x"
											/>
										</svg>
									</span>
								</div>
								{{ comment.content }}
							</div>
						</div>
					</li>
				</ul>
			</div>
		</div>
		<div class="card-footer">
			<small class="text-muted">
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#calendar" />
				</svg>
				<div>Uploaded: {{ formatDate(uploadDate) }}</div>
			</small>
		</div>
	</div>
</template>

<style scoped>
.card-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
.card-header .user-badge {
	text-decoration: none;
}
.card-header .user-badge:hover {
	color: var(--bs-white);
}
.card-header .user-badge .feather {
	width: 9px;
	height: 9px;
}
.card-img-top {
	border-radius: 0;
}
.action-btn {
	width: 50%;
	border-radius: 0;
	display: flex;
	justify-content: center;
	align-items: center;
}
.action-btn.action-btn:hover {
	background-color: var(--bs-btn-bg);
	color: var(--bs-btn-color);
	box-shadow: var(--bs-btn-focus-box-shadow);
}
.action-btn.active:hover {
	background-color: var(--bs-btn-active-bg);
	color: var(--bs-btn-active-color);
}
.action-btn:focus,
.action-btn:active {
	background-color: var(--bs-btn-bg);
	color: var(--bs-btn-color);
	box-shadow: none;
}
.action-btn.active:focus,
.action-btn.active:active {
	background-color: var(--bs-btn-active-bg);
	color: var(--bs-btn-active-color);
	box-shadow: none;
}
.action-btn > svg {
	margin-left: 4px;
}
.badge-primary {
	background-color: var(--bs-primary);
}
.badge-secondary {
	background-color: var(--bs-secondary);
}
.badge-secondary:first-child {
	margin-right: 4px;
}
.dropdown-toggle {
	cursor: pointer;
}
.dropdown-toggle::after {
	vertical-align: 0.1em;
}
.dropdown-toggle:hover {
	background-color: var(--bs-primary);
}
.dropdown-menu {
	display: block;
}
.dropdown-item {
	display: flex;
	align-items: center;
}
.dropdown-item > svg {
	margin-right: 4px;
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
.media .media-body > div {
	display: flex;
	align-items: center;
}
.media .media-body b {
	margin-right: 4px;
}
.media .media-body small {
	margin-left: 4px;
	line-height: 11px;
}
.media .media-body .close {
	margin-left: auto;
}
.card-footer small {
	display: flex;
	align-items: center;
}
.card-footer small svg {
	margin-right: 4px;
}
.load-wrapper{
  position: relative;
  height: 200px;
  width: 100%;
  background-color: rgb(211,211,211);
  z-index: 44;
  overflow: hidden;
  border-radius: 5px;
}
.activity{
  position: absolute;
  left: -45%;
  height: 100%;
  width: 45%;
  background-image: linear-gradient(to left, rgba(251,251,251, .05), rgba(251,251,251, .3), rgba(251,251,251, .6), rgba(251,251,251, .3), rgba(251,251,251, .05));
  animation: loading 0.4s infinite;
  z-index: 45;
}

@keyframes loading {
  0%{
    left: -45%;
  }
  100%{
    left: 100%;
  }
}
</style>
