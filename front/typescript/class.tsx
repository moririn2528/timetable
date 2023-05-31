import { getClass, InputGraphType, server, sleep } from "./api";
import React, { useEffect } from "react";

let nodes: {
	elems: HTMLDivElement[];
	ids: number[];
	id2index: Map<number, number>;
} = {
	elems: [],
	ids: [],
	id2index: new Map(),
};
let edges: {
	vertex: [number, number][];
} = {
	vertex: [],
};
let node_moving: {
	target: HTMLDivElement | null;
	x: number;
	y: number;
} = {
	target: null,
	x: 0,
	y: 0,
};

type Pos = {
	x: number;
	y: number;
}

const NodeComponent = React.forwardRef((props: { id: number, name: string, pos: Pos, setPos: (pos: Pos) => void }, ref: React.ForwardedRef<HTMLDivElement>) => {
	const [diff, setDiff] = React.useState<Pos>({ x: 0, y: 0 });
	const [moving, setMoving] = React.useState(false);
	const onMouseDown = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		if (e.target !== e.currentTarget) return;
		console.log("down", props.name, e.target, e.currentTarget);
		setDiff({ x: e.pageX - props.pos.x, y: e.pageY - props.pos.y });
		setMoving(true);
	}
	const onMouseMove = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		if (!moving) return;
		console.log("move", props.id, props.name, props.pos, e.pageX, e.pageY, diff);
		props.setPos({ x: e.pageX - diff.x, y: e.pageY - diff.y });
	}
	const onMouseUp = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		setMoving(false);
	}
	return (
		<div style={{
			left: props.pos.x,
			top: props.pos.y,
			height: "20px",
			userSelect: "none",
			padding: "10px",
			background: "#fff",
			border: "2px solid #000000",
			boxShadow: "0 0 2px 1px rgba(0, 0, 0, 0.2)",
			position: "absolute",
			zIndex: "auto",
			whiteSpace: "nowrap",
		}} onMouseDown={onMouseDown} onMouseMove={onMouseMove} onMouseUp={onMouseUp} ref={ref}>
			{props.name}
		</div >
	);
})

function str2posMap(data_str: string) {
	const a = data_str.split("A");
	const id_arr = a[0].split(",");
	const pos_arr = a[1].split(",");
	console.assert(id_arr.length == pos_arr.length);
	const poss = pos_arr.map((str) => {
		const xy = str.split(":");
		return { x: Number(xy[0]), y: Number(xy[1]) };
	})
	return new Map(id_arr.map((str, i) => [Number(str), poss[i]]));
}

function posMap2str(poss: Map<number, Pos>) {
	const id_arr = Array.from(poss.keys()).map((id) => id.toString());
	const pos_arr = Array.from(poss.values()).map((pos) => `${pos.x}:${pos.y}`);
	return id_arr.join(",") + "A" + pos_arr.join(",");
}

function loadPosMap() {
	let data_str = sessionStorage.getItem("nodes");
	if (data_str != null) {
		return str2posMap(data_str);
	}
	data_str = localStorage.getItem("nodes");
	if (data_str != null) {
		return str2posMap(data_str);
	}
	return new Map<number, Pos>();
}

function savePosTemp(poss: Map<number, Pos>) {
	const data_str = posMap2str(poss);
	sessionStorage.setItem("nodes", data_str);
}

function savePos(poss: Map<number, Pos>) {
	const data_str = posMap2str(poss);
	localStorage.setItem("nodes", data_str);
	console.log("saved", data_str);
}

export function GraphComponent(props: { classes: InputGraphType }) {
	const canvasRef = React.useRef<HTMLCanvasElement>(null);
	const [poss, setPoss] = React.useState<Map<number, Pos>>(
		() => {
			const pos_map = loadPosMap();
			for (let i = 0; i < props.classes.nodes.length; i++) {
				const node = props.classes.nodes[i];
				if (!pos_map.has(node.id)) {
					pos_map.set(node.id, { x: 0, y: 0 });
				}
			}
			return pos_map;
		}
	);
	const node_refs = new Map(props.classes.nodes.map((node) => [node.id, React.useRef<HTMLDivElement>(null)]))
	const [diffs, setDiffs] = React.useState<Map<number, Pos>>(
		new Map(props.classes.nodes.map((node) => [node.id, { x: 0, y: 0 }]))
	);
	const [moving, setMoving] = React.useState(false);

	useEffect(() => {
		savePosTemp(poss);
	}, [poss]);
	const draw_edges = () => {
		const canvas = canvasRef.current;
		if (canvas == null) throw "canvas is null";
		const ctx = canvas.getContext("2d");
		if (ctx == null) throw "ctx is null";
		ctx.clearRect(0, 0, canvas.clientWidth, canvas.clientHeight);
		ctx.strokeStyle = "black";
		ctx.fillStyle = "black";
		ctx.lineWidth = 2;
		const corner_x = canvas.getBoundingClientRect().left;
		const corner_y = canvas.getBoundingClientRect().top;
		for (let i = 0; i < props.classes.edges.length; i++) {
			const edge = props.classes.edges[i];
			const from_pos = poss.get(edge.from);
			const to_pos = poss.get(edge.to);
			const from_ref = node_refs.get(edge.from);
			const to_ref = node_refs.get(edge.to);
			if (from_pos == null || to_pos == null) throw "pos is null";
			if (from_ref == null || to_ref == null) throw "ref is null";
			if (from_ref.current == null || to_ref.current == null) throw "edge ref.current is null";
			const from_x = from_pos.x + from_ref.current.clientWidth / 2 - corner_x;
			const from_y = from_pos.y + from_ref.current.clientHeight / 2 - corner_y;
			const to_x = to_pos.x + to_ref.current.clientWidth / 2 - corner_x;
			const to_y = to_pos.y + to_ref.current.clientHeight / 2 - corner_y;
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
	}
	useEffect(draw_edges, [props.classes, poss]);
	useEffect(() => {
		const canvas_resize = () => {
			const canvas = canvasRef.current;
			if (canvas == null) throw "canvas is null";
			canvas.width = canvas.clientWidth;
			canvas.height = canvas.clientHeight;
			draw_edges();
		}
		canvas_resize();
		window.addEventListener("resize", canvas_resize);
		return () => {
			window.removeEventListener("resize", canvas_resize);
		}
	}, []);

	const onMouseDown = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		if (e.target !== e.currentTarget) return;
		console.log("down all", e.target, e.currentTarget);
		node_moving.x = e.pageX;
		node_moving.y = e.pageY;
		setDiffs(
			new Map(props.classes.nodes.map((node) => {
				const pos = poss.get(node.id);
				if (pos == null) throw "pos is null";
				return [node.id, { x: e.pageX - pos.x, y: e.pageY - pos.y }];
			}))
		)
		setMoving(true);
	}

	const onMouseMove = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		const mx = e.pageX, my = e.pageY;
		if (!moving) return;
		setPoss((prev) => {
			const pos_map = new Map(prev);
			prev.forEach((pos, id) => {
				const diff = diffs.get(id);
				if (diff == null) throw "diff is null";
				pos_map.set(id, { x: mx - diff.x, y: my - diff.y });
			})
			return pos_map;
		})
	}

	const onMouseUp = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		setMoving(false);
	}

	const onClickSavePos = () => {
		savePos(poss);
	}

	return (
		<div style={{
			width: "100%",
			height: "calc(100% - 40px)",
		}} onMouseDown={onMouseDown} onMouseMove={onMouseMove} onMouseUp={onMouseUp}>
			<div onClick={onClickSavePos} style={{
				width: 50,
				userSelect: "none",
				padding: 2,
				border: "1px solid #000000",
				textAlign: "center",
			}}>save</div>
			<canvas style={{
				width: "100%",
				height: "calc(100% - 50px)",
				zIndex: -1,
				position: "absolute",
			}} ref={canvasRef}></canvas>
			{props.classes.nodes.map((node) => {
				const pos = poss.get(node.id);
				if (pos == null) throw "pos is null";
				return <NodeComponent id={node.id} name={node.name} pos={pos} setPos={(pos) => {
					setPoss((prev) => {
						const pos_map = new Map(prev);
						pos_map.set(node.id, pos);
						return pos_map;
					})
				}} ref={node_refs.get(node.id)!} key={node.id} />
			})}

		</div>
	)
}
