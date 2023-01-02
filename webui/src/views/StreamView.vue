<script>
import Feed from "../components/shared/Feed.vue";

export default {
	components: {
		Feed,
	},
	data() {
		return {
			loading: false,
		};
	},
	methods: {
		async photoFetcher(config) {
			try {
				const response = await this.$axios.get("/users/me/stream", config);
				return response.data;
			} catch (e) {
				console.error(e)
			}
		},
	},
	mounted() {
		this.$refs.feed.fetchPhotos();
	},
};
</script>

<template>
	<div
		class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
	>
		<h1 class="h2">Stream</h1>
		<small class="text-muted"></small>
	</div>
	<Feed
		ref="feed"
		:fetcher="photoFetcher"
		@fetchPhotosStarted="loading = true"
		@loadMorePhotosStarted="loading = true"
		@fetchPhotosEnded="loading = false"
		@loadMorePhotosEnded="loading = false"
	/>

	<LoadingSpinner :loading="loading" marginTop="33vh" />
</template>

<style>
.container {
	padding: 0;
}
</style>
