import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig(({ command, mode, ssrBuild }) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				"@": fileURLToPath(new URL("./src", import.meta.url)),
			},
		},
		...process.env.APP_BASE_PATH && { base: process.env.APP_BASE_PATH },
	};
	ret.define = {
		__API_URL__: JSON.stringify(
			process.env.API_URL || "http://localhost:3000"
		),
		__STATIC_FILES_URL__: JSON.stringify(
			process.env.STATIC_FILES_URL || "http://localhost:3000"
		),
	};
	return ret;
});
