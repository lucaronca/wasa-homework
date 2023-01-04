<script>
import { loggedInUser } from "../user";
import Feed from "../components/shared/Feed.vue";
import FollowsLists from "../components/follows/FollowsLists.vue";
import UsernameModal from "../components/cta/UsernameModal.vue";
import FollowButton from "../components/follows/FollowButton.vue";
import PhotoUploader from "../components/cta/PhotoUploader.vue";
import BanButton from "../components/cta/BanButton.vue";

export default {
	components: {
		Feed,
		FollowsLists,
		UsernameModal,
		FollowButton,
		PhotoUploader,
		BanButton,
	},
	data() {
		return {
			loading: false,
			id: null,
			user: null,
			loggedInUser,
			isFollowedByLoggedInUser: false,
		};
	},
	computed: {
		isSignedInUserProfile() {
			return this.id == "me";
		},
	},
	methods: {
		async refresh() {
			this.loading = true;
			await Promise.all([
				this.$refs.feed.fetchPhotos(),
				this.getUserInfo(),
			]);
			this.loading = false;
		},
		async photoFetcher(config) {
			try {
				const response = await this.$axios.get(
					`/users/${this.id}/photos`,
					config
				);
				return response.data;
			} catch (e) {
				if (e.response.status == 404) {
					return [];
				} else {
					console.error(e);
					if (e.response?.data) {
						window.alert(e.response.data);
					} else {
						window.alert("Api error");
					}
				}
			}
		},
		async getFollowsListsRef() {
			if (this.$refs.followsLists) {
				return this.$refs.followsLists
			} else {
				await this.$nextTick()
				return this.$refs.followsLists
			}
		},
		async getUserInfo() {
			try {
				const response = await this.$axios.get(`/users/${this.id}`);
				this.user = response.data;
				if (this.id !== "me" && !this.user.bannedForUser) {
					this.isFollowedByLoggedInUser = await (
						await this.getFollowsListsRef()
					).isFollowedByLoggedInUser();
				}
			} catch (e) {
				console.error(e);
				if (e.response?.data) {
					window.alert(e.response.data);
					if (e.response.status == 404) {
						this.$router.replace("/");
					} else {
						window.alert("Api error");
					}
				}
			}
		},
		onFollowUpdated(newFollower) {
			if (newFollower) {
				this.isFollowedByLoggedInUser = true;
				this.user.totalFollowers++;
			} else {
				this.isFollowedByLoggedInUser = false;
				this.user.totalFollowers--;
			}
		},
	},
	mounted() {
		this.id =
			this.$route.params.id == loggedInUser.id
				? "me"
				: this.$route.params.id;
		this.refresh();
	},
	async beforeRouteUpdate(to) {
		(await this.getFollowsListsRef())?.hideLists();
		this.id = to.params.id == loggedInUser.id ? "me" : to.params.id;
		this.refresh();
	},
};
</script>

<template>
	<div v-show="!loading" class="pt-3 pb-2 mb-3 border-bottom">
		<div class="header-item mb-2">
			<h1 class="h2">{{ user?.username }}'s Profile</h1>
			<h5 v-if="!user?.bannedForUser">
				Total photos:
				<span class="badge badge-secondary">{{
					user?.totalPhotos
				}}</span>
			</h5>
		</div>
		<div class="header-item mb-2">
			<FollowsLists
				v-if="!user?.bannedForUser"
				ref="followsLists"
				:userId="id"
				:totalFollowers="user?.totalFollowers"
				:totalFollowings="user?.totalFollowings"
			/>
			<div class="btn-group mb-0">
				<div v-if="isSignedInUserProfile" class="btn-group">
					<PhotoUploader @photoUploaded="refresh" />
					<UsernameModal @usernameUpdated="refresh" />
				</div>
				<div v-else class="btn-group">
					<FollowButton
						v-if="!user?.bannedForUser"
						:isFollowedByLoggedInUser="isFollowedByLoggedInUser"
						:userId="id"
						@followUpdated="onFollowUpdated"
					/>
					<BanButton
						:bannedForUser="user?.bannedForUser"
						:userId="id"
						@banUpdated="refresh"
					/>
				</div>
			</div>
		</div>
	</div>
	<Feed
		ref="feed"
		:fetcher="photoFetcher"
		@photoDeleted="refresh"
		@loadMorePhotosStarted="loading = true"
		@loadMorePhotosEnded="loading = false"
	/>
	<LoadingSpinner :loading="loading" marginTop="40vh" />
</template>

<style scoped>
.container {
	padding: 0;
}
h1,
h5 {
	margin-bottom: 0;
}
.badge-secondary {
	background-color: var(--bs-secondary);
}
.badge-light {
	color: var(--bs-dark);
	background-color: var(--bs-gray-100);
}
.header-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	width: 100%;
}
.header-item:last-child {
	flex-direction: column;
	align-items: baseline;
}
@media (min-width: 576px) {
	.header-item:last-child {
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
	}
}
</style>
