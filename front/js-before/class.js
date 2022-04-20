const server = {
	get: async function (path) {
		const port = location.port;
		console.log(port);
		const res = await fetch("http://localhost:" + port + "/api/" + path);
		if (res.status != 200) {
			console.error(res);
			return;
		}
		return res.json();
	},
};
const sleep = (waitTime) => new Promise((resolve) => setTimeout(resolve, waitTime));

(function () {
	let nodes = {
		elems: [],
		ids: [],
		id2index: {},
	};
	let edges = {
		vertex: [],
	};
	let node_moving = {
		target: null,
		x: 0,
		y: 0,
		width: 0,
	};

	const node_event = {
		str2pos: function (data_str) {
			const a = data_str.split("A");
			const ids = a[0];
			const poss = a[1];
			const id_arr = ids.split(",");
			const pos_arr = poss.split(",");
			console.assert(id_arr.length == pos_arr.length);
			for (let i = 0; i < id_arr.length; i++) {
				const id = id_arr[i];
				const pos_str = pos_arr[i].split(":");
				const x = pos_str[0];
				const y = pos_str[1];
				console.assert(id in nodes.id2index);
				const idx = nodes.id2index[id];
				const elem = nodes.elems[idx];
				elem.style.left = `${x}px`;
				elem.style.top = `${y}px`;
			}
		},
		pos2str: function () {
			pos_arr = [];
			for (let i = 0; i < nodes.elems.length; i++) {
				const elem = nodes.elems[i];
				const rect = elem.getBoundingClientRect();
				const x = rect.left;
				pos_arr.push(String(rect.left) + ":" + String(rect.top));
			}
			return nodes.ids.join(",") + "A" + pos_arr.join(",");
		},

		load_pos: function () {
			let data_str = sessionStorage.getItem("nodes");
			if (data_str != null) {
				this.str2pos(data_str);
				return;
			}
			data_str = localStorage.getItem("nodes");
			if (data_str != null) {
				this.str2pos(data_str);
			}
		},
		temp_save_pos: function () {
			let data_str = this.pos2str();
			sessionStorage.setItem("nodes", data_str);
		},
		save_pos: function () {
			let data_str = this.pos2str();
			localStorage.setItem("nodes", data_str);
			console.log("saved", data_str);
		},

		draw_edges: function () {
			const canvas = document.getElementById("graph_canvas");
			const field = document.getElementById("graph_field");
			canvas.width = field.clientWidth;
			canvas.height = field.clientHeight;
			const field_rect = field.getBoundingClientRect();
			const corner_x = field_rect.left;
			const corner_y = field_rect.top;
			const ctx = canvas.getContext("2d");
			ctx.clearRect(0, 0, canvas.width, canvas.height);
			ctx.strokeStyle = "black";
			ctx.fillStyle = "black";
			ctx.lineWidth = 2;
			for (let i = 0; i < edges.vertex.length; i++) {
				const edge = edges.vertex[i];
				const f = edge[0];
				const t = edge[1];
				const from_rect = nodes.elems[f].getBoundingClientRect();
				const to_rect = nodes.elems[t].getBoundingClientRect();
				const from_x = from_rect.left + from_rect.width / 2 - corner_x;
				const from_y = from_rect.top + from_rect.height / 2 - corner_y;
				const to_x = to_rect.left + to_rect.width / 2 - corner_x;
				const to_y = to_rect.top + to_rect.height / 2 - corner_y;
				const center_x = (from_x + to_x) / 2;
				const center_y = (from_y + to_y) / 2;
				const dis = Math.sqrt((to_x - from_x) ** 2 + (to_y - from_y) ** 2);
				const vec_x = (to_x - from_x) / dis;
				const vec_y = (to_y - from_y) / dis;
				const tri_hei = Math.min(10, dis / 4);

				ctx.beginPath();
				ctx.moveTo(from_x, from_y);
				ctx.lineTo(to_x, to_y);
				ctx.stroke();
				ctx.beginPath();
				ctx.moveTo(center_x + tri_hei * vec_x, center_y + tri_hei * vec_y);
				ctx.lineTo(center_x - tri_hei * vec_y, center_y + tri_hei * vec_x);
				ctx.lineTo(center_x + tri_hei * vec_y, center_y - tri_hei * vec_x);
				ctx.closePath();
				ctx.fill();
			}
		},

		mouse_down: function (e) {
			const target = e.currentTarget;
			//const targetW = target.width;
			const target_rect = target.getBoundingClientRect();

			node_moving.target = target;
			node_moving.x = target_rect.left - e.pageX;
			node_moving.y = target_rect.top - e.pageY;
			//target.style.width = `${targetW}px`;
			target.classList.add("moving");
			window.addEventListener("mousemove", node_event.mouse_move);
			window.addEventListener("mouseup", node_event.mouse_up);
		},

		mouse_move: function (e) {
			const target = node_moving.target;

			let x = e.pageX + node_moving.x;
			let y = e.pageY + node_moving.y;
			if (x < 0) x = 0;
			if (y < 0) y = 0;
			target.style.left = `${x}px`;
			target.style.top = `${y}px`;
			//console.log(x,y);
			node_event.draw_edges();
		},

		mouse_up: function () {
			const target = node_moving.target;
			//console.log(target.style.left,target.style.top);
			target.classList.remove("moving");
			window.removeEventListener("mousemove", node_event.mouse_move);
			window.removeEventListener("mouseup", node_event.mouse_up);
			node_event.temp_save_pos();
			node_event.draw_edges();
		},
	};
	const board_event = {
		mouse_down: function (e) {
			if (e.target.classList.contains("node")) {
				return;
			}
			const target = e.currentTarget;
			target.classList.add("moving");
			node_moving.target = target;
			node_moving.x = e.pageX;
			node_moving.y = e.pageY;
			window.addEventListener("mousemove", board_event.mouse_move);
			window.addEventListener("mouseup", board_event.mouse_up);
		},

		mouse_move: function (e) {
			const mx = e.pageX,
				my = e.pageY;
			for (let i = 0; i < nodes.elems.length; i++) {
				const target = nodes.elems[i];
				const target_rect = target.getBoundingClientRect();
				let x = mx - node_moving.x + target_rect.left;
				let y = my - node_moving.y + target_rect.top;
				target.style.left = `${x}px`;
				target.style.top = `${y}px`;
			}
			(node_moving.x = mx), (node_moving.y = my);
			node_event.draw_edges();
		},

		mouse_up: function () {
			const target = node_moving.target;
			target.classList.remove("moving");
			window.removeEventListener("mousemove", board_event.mouse_move);
			window.removeEventListener("mouseup", board_event.mouse_up);
			node_event.temp_save_pos();
			node_event.draw_edges();
		},
	};

	window.addEventListener("load", async function (e) {
		const res = await server.get("class");
		//console.log(res);
		const parent_div = this.document.getElementById("graph_field");
		parent_div.addEventListener("mousedown", board_event.mouse_down);
		for (let i = 0; i < res.nodes.length; i++) {
			const node = res.nodes[i];
			const id = node.id;
			const name = node.name;
			const div = document.createElement("div");
			div.className = "node";
			div.innerText = name;
			div.addEventListener("mousedown", node_event.mouse_down);
			parent_div.appendChild(div);
			nodes.id2index[id] = nodes.elems.length;
			nodes.elems.push(div);
			nodes.ids.push(id);
		}
		node_event.load_pos();

		for (let i = 0; i < res.edges.length; i++) {
			const edge = res.edges[i];
			const from = nodes.id2index[edge.from];
			const to = nodes.id2index[edge.to];
			edges.vertex.push([from, to]);
		}
		node_event.draw_edges();
	});

	document.getElementById("nodes_pos_save").addEventListener("click", async function (e) {
		const elem = e.currentTarget;
		node_event.save_pos();
		elem.classList.add("clicked_color");
		await sleep(2000);
		elem.classList.remove("clicked_color");
	});
})();
