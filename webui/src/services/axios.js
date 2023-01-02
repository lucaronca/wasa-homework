import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
});

instance.interceptors.request.use((config) => {
	const token = window.localStorage.getItem('$loggedInUserToken')
	config.headers = {
		...config.headers,
	  	Authorization: `Bearer ${token}`,
	}
	return config
})

export default instance;
