export const server = {
	get: async function (path: string) {
		const port = location.port;
		console.log("http://localhost:" + port + "/api/" + path);
		const res = await fetch("http://localhost:" + port + "/api/" + path);
		if (res.status != 200) {
			console.error(res);
			return;
		}
		return res.json();
	},
};
export const sleep = (waitTime: number) => new Promise((resolve) => setTimeout(resolve, waitTime));
