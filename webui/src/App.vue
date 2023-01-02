<script>
import { RouterLink, RouterView } from "vue-router";
import { loggedInUser } from "./user";

export default {
	async mounted() {
		try {
			if (!window.localStorage.getItem("$loggedInUserToken")) {
				this.$router.replace("/login")
			} else {
				const resp = await this.$axios.get("/users/me");
				loggedInUser.id = resp.data.id;
				loggedInUser.username = resp.data.username;
			}
		} catch (e) {
			console.error(e);
			if (e.response?.data) {
				window.alert(e.response.data);
			} else {
				window.alert("Login failed!");
			}
		} finally {
			this.ready = true
		}
	},
	data() {
		return {
			ready: false,
			loggedInUser,
		};
	},
	methods: {
		logout() {
			loggedInUser.id = null;
			loggedInUser.username = null;
			window.document.title = 'Wasa Photo'
			this.$router.replace({ name: "Login" });
		},
	},
};
</script>

<template>
	<span v-if="ready">
		<header
			class="navbar navbar-dark sticky-top flex-md-nowrap p-0 shadow"
		>
			<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">
				Wasa Photo<span v-if="loggedInUser?.username"> - {{ loggedInUser.username }}</span>
			</a>
			<button
				class="navbar-toggler position-absolute d-md-none collapsed"
				type="button"
				data-bs-toggle="collapse"
				data-bs-target="#sidebarMenu"
				aria-controls="sidebarMenu"
				aria-expanded="false"
				aria-label="Toggle navigation"
			>
				<span class="navbar-toggler-icon"></span>
			</button>
		</header>

		<div class="container-fluid">
			<div class="row">
				<nav
					v-if="loggedInUser.id"
					id="sidebarMenu"
					class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse"
				>
					<div class="position-sticky pt-3 sidebar-sticky">
						<ul class="nav flex-column">
							<li class="nav-item">
								<RouterLink to="/" class="nav-link">
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#rss"
										/>
									</svg>
									Stream
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink
									:to="{
										name: 'Profile',
										params: { id: loggedInUser.id },
									}"
									class="nav-link"
								>
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#user"
										/>
									</svg>
									Profile
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink to="/search" class="nav-link">
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#search"
										/>
									</svg>
									Search
								</RouterLink>
							</li>
							<li class="nav-item">
								<a
									href="#"
									class="nav-link"
									@click.prevent="logout"
								>
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#log-out"
										/>
									</svg>
									Log out
								</a>
							</li>
						</ul>
					</div>
				</nav>

				<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
					<RouterView />
				</main>
			</div>
		</div>
	</span>
	<div v-else>Loading...</div>
</template>

<style>
	:root {
		--primary: #0B9095;
		--bs-primary: var(--primary);
	}

	.navbar {
		background-color: #00777C;
	}

	.btn-primary {
		--bs-btn-bg: var(--primary);
		--bs-btn-border-color: var(--primary);
		--bs-btn-hover-bg: var(--primary);
		--bs-btn-hover-border-color: var(--primary);
	}

	.btn-outline-primary {
		--bs-btn-color: var(--primary);
		--bs-btn-border-color: var(--primary);
		--bs-btn-active-bg: var(--primary);
		--bs-btn-active-border-color: var(--primary);
		--bs-btn-focus-shadow-rgb: 11, 144, 149;
		--bs-btn-hover-bg: var(--primary);
		--bs-btn-hover-border-color: var(--primary);
	}

	.form-control:focus {
		border-color: rgb(11 144 149 / 25%);
		box-shadow: 0 0 0 0.25rem rgb(11 144 149 / 25%);
	}
</style>
