<script>
export default {
	props: ["loading", "marginTop"]
}
</script>

<template>
	<div v-if="loading" class="spinner" :style="{marginTop: marginTop}">
		<div>
			<div class="spinner-border" role="status">
				<span class="visually-hidden">Loading...</span>
			</div>
		</div>
	</div>
	<div v-if="!loading"><slot /></div>
</template>

<style>
.spinner {
	display: flex;
    justify-content: center;
    align-items: center;
}
</style>
