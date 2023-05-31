import { decodeType, record, number, string, array } from "typescript-json-decoder";

export const server = {
	get: async function (path: string) {
		console.log(location.origin + "/api/" + path);
		const res = await fetch(location.origin + "/api/" + path);
		if (!res.ok) {
			console.error(res);
			throw `response error, code: ${res.status}`;
		}
		return res.json();
	},

	post: async function (path: string, data: any) {
		console.log(location.origin + "/api/" + path);
		const res = await fetch(location.origin + "/api/" + path, {
			method: "POST",
			mode: "same-origin",
			cache: "no-cache",
			credentials: "same-origin",
			headers: {
				"Content-Type": "application/json",
			},
			redirect: "follow",
			referrerPolicy: "no-referrer",
			body: JSON.stringify(data),
		});
		if (!res.ok) {
			console.error(res);
			throw `response error, code: ${res.status}`;
		}
		return res.json();
	},
};
export const sleep = (waitTime: number) => new Promise((resolve) => setTimeout(resolve, waitTime));

export const jsdate = function (val: any) {
	return new Date(string(val));
};

export type InputGraphType = decodeType<typeof inputGraphDecoder>;
const inputGraphDecoder = record({
	nodes: array({ id: number, name: string }),
	edges: array({ from: number, to: number }),
});
export const getClass = () => {
	return server.get("class").then(inputGraphDecoder);
};


export type TeacherType = decodeType<typeof inputTeacherDecoder>;
const inputTeacherDecoder = record({
	id: number,
	name: string,
});
export async function getTeachers() {
	return server.get("teacher").then(array(inputTeacherDecoder));
}
