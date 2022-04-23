import { decodeType, record, number, string, array } from "typescript-json-decoder";

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

export type InputGraphType = decodeType<typeof inputGraphDecoder>;
const inputGraphDecoder = record({
	nodes: array({ id: number, name: string }),
	edges: array({ from: number, to: number }),
});
export const getClass = () => {
	return server.get("class").then(inputGraphDecoder);
};
