<script>
import { loggedInUser } from '../../user';
export default {
	props: {
		userId: String,
		totalFollowers: Number,
		totalFollowings: Number,
	},
	data() {
		return {
			followersList: [],
			showFollowersList: false,
			followersListLoading: false,
			followingsList: [],
			showFollowingsList: false,
			followingsListLoading: false,
		}
	},
	watch: {
		async showFollowersList(newShowFollowersList) {
			if (newShowFollowersList) {
				this.showFollowingsList = false;
				this.followersListLoading = true;
				try {
					await this.updateFollowersList()
				} catch (err) {
					console.error(err);
				} finally {
					this.followersListLoading = false;
				}
			}
		},
		async showFollowingsList(newShowFollowingsList) {
			if (newShowFollowingsList) {
				this.showFollowersList = false
				this.followingsListLoading = true;
				try {
					this.followingsList = (
						await this.$axios.get(
							`/users/${this.userId}/followings`
						)
					).data;
				} catch (err) {
					console.error(err);
				} finally {
					this.followingsListLoading = false;
				}
			}
		},
	},
	methods: {
		async updateFollowersList() {
			this.followersList = (
				await this.$axios.get(`/users/${this.userId}/followers`)
			).data;
		},
		hideLists() {
			this.showFollowersList = false;
			this.showFollowingsList = false;
		},
		async isFollowedByLoggedInUser() {
			await this.updateFollowersList()
			return this.followersList.some(
				f => f.id === loggedInUser.id
			)
		}
	},
}
</script>
<template>
<div class="btn-group mb-2 mb-md-0">
	<div class="btn-group">
		<button
			type="button"
			class="btn btn-secondary btn-sm"
			:class="{
				'dropdown-toggle': totalFollowers !== 0,
			}"
			@click="
				totalFollowers !== 0 &&
					(showFollowersList = !showFollowersList)
			"
		>
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#list" />
			</svg>
			Followers
			<span class="badge badge-light">{{
				totalFollowers
			}}</span>
		</button>
		<div v-show="showFollowersList" class="dropdown-menu">
			<LoadingSpinner v-if="followersListLoading" :loading="followersListLoading" />
			<ProfileLink
				v-else
				v-for="follower in followersList"
				:key="follower.id"
				:id="follower.id"
				:name="follower.username"
			/>
		</div>
	</div>
	<div class="btn-group">
		<button
			type="button"
			class="btn btn-secondary btn-sm"
			:class="{
				'dropdown-toggle': totalFollowings !== 0,
			}"
			@click="
				totalFollowings !== 0 &&
					(showFollowingsList = !showFollowingsList)
			"
		>
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#list" />
			</svg>
			Following
			<span class="badge badge-light">{{
				totalFollowings
			}}</span>
		</button>
		<div v-show="showFollowingsList" class="dropdown-menu">
			<LoadingSpinner v-if="followingsListLoading" :loading="followingsListLoading" />
			<ProfileLink
				v-else
				v-for="following in followingsList"
				:key="following.id"
				:id="following.id"
				:name="following.username"
			/>
		</div>
	</div>
</div>
</template>
<style scoped>
.dropdown-menu {
	display: block;
	position: absolute;
	top: 40px;
}
.btn-group:last-child > button.btn {
	border-top-right-radius: 0.375rem;
	border-bottom-right-radius: 0.375rem;
}
button.btn {
	display: flex;
	align-items: center;
}
button.btn > svg {
	margin-right: 0.5rem;
}
button.btn .badge {
	margin-left: 0.5rem;
	margin-right: 0.25rem;
    padding: 0;
	position: static;
	margin-top: 1px;
}
</style>
