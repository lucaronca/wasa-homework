<script setup>
import { RouterLink } from 'vue-router'
</script>
<script>
	export default {
		props: {
			id: Number,
			name: String
		}
	}
</script>
<template>
	<RouterLink
		class="dropdown-item"
		:to="{ name: 'Profile', params: { id }}"
	>
		<svg class="feather">
			<use
				href="/feather-sprite-v4.29.0.svg#user"
			/>
		</svg>
		{{ name }}
	</RouterLink>
</template>
<style scoped>
.feather {
	margin-right: 0.5rem;
}
.dropdown-item {
	display: flex;
    align-items: center;
}
</style>
