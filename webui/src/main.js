import { createApp } from "vue";

import App from "./App.vue";
import router from "./router";
import axios from "./services/axios.js";
import LoadingSpinner from "./components/shared/LoadingSpinner.vue";
import ProfileLink from "./components/shared/ProfileLink.vue";
import Photo from "./components/photo/Photo.vue";

import "./assets/dashboard.css";
import "./assets/main.css";

const app = createApp(App);

app.config.globalProperties.$axios = axios;

app.component("LoadingSpinner", LoadingSpinner);
app.component("ProfileLink", ProfileLink);
app.component("Photo", Photo);

app.use(router);

app.mount("#app");
