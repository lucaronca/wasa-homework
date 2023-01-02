import {createRouter, createWebHashHistory} from 'vue-router'
import StreamView from '../views/StreamView.vue'
import ProfileView from '../views/ProfileView.vue'
import SearchView from '../views/SearchView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/login', component: LoginView, 'name': 'Login'},
		{path: '/', component: StreamView},
		{path: '/profile/:id', component: ProfileView, 'name': 'Profile'},
		{path: '/search', component: SearchView, 'name': 'Search'},
	]
})

export default router
