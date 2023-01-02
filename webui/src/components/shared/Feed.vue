<script>
export default {
	props: {
		fetcher: Function,
	},
	data() {
		return {
			photos: [],
			totalPhotos: 0,
		}
	},
	computed: {
		showLoadMoreBtn() {
			return this.photos.length < this.totalPhotos
		},
		staticFilesURL() {
			return __STATIC_FILES_URL__
		}
	},
	methods: {
		async fetchPhotos() {
			this.$emit('fetchPhotosStarted')
			const resp = await this.fetcher({})
			this.photos = resp.entries
			this.totalPhotos = resp.totalCount
			this.$emit('fetchPhotosEnded')
		},
		async loadMorePhotos() {
			this.$emit('loadMorePhotosStarted')
			try {
				const resp = await this.fetcher(
					{
						params: {
							offset: this.photos.length,
						}
					},
				)
				this.photos = [...this.photos, ...resp.entries];
				this.totalPhotos = resp.totalCount
				setTimeout(() => {
					document.querySelector('.photo:last-child').scrollIntoView()
				})
			} catch (e) {
				console.error(e);
				if (e.response?.data) {
					window.alert(e.response.data);
				} else {
					window.alert("Api error");
				}
			} finally {
				this.$emit('loadMorePhotosEnded')
			}
		},
	}
}
</script>
<template>
	<div class="container mb-3">
		<div class="row gx-3 gy-3">
			<div
				v-for="photo in photos"
				:key="photo.id"
				class="col-sm-6 col-lg-4 col-xl-3 photo"
			>
				<Photo
					:id="photo.id"
					:image="staticFilesURL + photo.url"
					:ownerUsername="photo.owner.username"
					:ownerId="photo.owner.id"
					:totalLikes="photo.totalLikes"
					:userLiked="photo.userLiked"
					:totalComments="photo.totalComments"
					:uploadDate="photo.uploadDate"
					@likeToggled="photo.userLiked = !photo.userLiked"
					@likesCountUpdated="photo.totalLikes = $event"
					@commentsCountUpdated="photo.totalComments = $event"
					@photoDeleted="$emit('photoDeleted')"
				/>
			</div>
		</div>
		<div v-if="showLoadMoreBtn" class="row gx-3 gy-3 mt-3">
			<div class="btn-group mx-auto" style="width: 200px">
				<button
					type="button"
					class="btn btn-primary"
					@click.prevent="loadMorePhotos"
				>
					Load more
				</button>
			</div>
		</div>
	</div>
</template>
<style></style>
