import { getClass, server, InputGraphType, sleep } from "./api";
import { decodeType, record, number, string, array, undef } from "typescript-json-decoder";
import React from "react";
import ReactDOM from "react-dom";
import { createRoot } from "react-dom/client";
import { HotUpdateChunk } from "webpack";

const jsdate = function (val: any) {
	return new Date(string(val));
};

type TimetableType = decodeType<typeof timetableDecoder>;
const timetableDecoder = record({
	id: number,
	class_id: number,
	class_name: string,
	duration_id: number,
	duration_name: string,
	frame_id: number,
	frame_day_week: number,
	frame_period: number,
	subject_id: number,
	subject_name: string,
	teacher_id: number,
	teacher_name: string,
	place_id: number,
	day: jsdate,
});
type TimetableMoveType = decodeType<typeof timetableMoveDecoder>;
const timetableMoveDecoder = record({
	timetable: timetableDecoder,
	day: jsdate,
	frame_id: number,
});

(function () {
	let search_data = {
		classes: [0],
	};
	let data = {
		day: new Date(2020, 3, 1),
		classes: [0],
	};
	interface UnitElemType {
		subject: HTMLTableCellElement;
		class: HTMLTableCellElement;
		teacher: HTMLTableCellElement;
	}
	let units: { elems: UnitElemType[]; objects: TimetableType[][] };
	units = {
		elems: [],
		objects: [], // [i][j]: frame_id=i, j 個目の object
	};

	const timetable = {
		init: function () {
			const frame_num = 42;
			units.elems = new Array(frame_num);
			units.objects = new Array(frame_num);
			const day_weak_strs = ["月", "火", "水", "木", "金", "土"];
			const tt = document.getElementById("timetable");
			if (tt == null) throw "getElementById timetable error";
			tt.innerHTML = "";
			const tt_head = document.createElement("tr");
			tt.appendChild(tt_head);
			tt_head.appendChild(document.createElement("th"));
			for (let i = 0; i < day_weak_strs.length; i++) {
				const day = new Date(data.day.getTime());
				let dif = i + 1 - day.getDay();
				if (dif < 0) {
					dif += 7;
				}
				day.setDate(day.getDate() + dif);
				const t = document.createElement("th");
				t.innerText = `${day.getMonth() + 1}/${day.getDate()}(${day_weak_strs[i]})`;
				tt_head.appendChild(t);
				if (dif == 0) {
					t.style.backgroundColor = "#95f9ef";
				}
			}

			const add_t_col = function (t: HTMLTableElement) {
				const tr = document.createElement("tr");
				t.appendChild(tr);
				const td = document.createElement("td");
				tr.appendChild(td);
				return td;
			};
			for (let i = 0; i < 7; i++) {
				const t_line = document.createElement("tr");
				tt.appendChild(t_line);
				const th = document.createElement("th");
				th.innerText = `${i + 1}限`;
				t_line.appendChild(th);
				for (let j = 0; j < 6; j++) {
					const t_wrap = document.createElement("td");
					t_wrap.classList.add("unit");
					t_line.appendChild(t_wrap);
					const t = document.createElement("table");
					t_wrap.appendChild(t);
					const t_subject = add_t_col(t);
					t_subject.classList.add("subject");
					const t_class = add_t_col(t);
					t_class.classList.add("class");
					const t_teacher = add_t_col(t);
					t_teacher.classList.add("teacher");
					const id = j * 7 + i;
					units.elems[id] = {
						subject: t_subject,
						class: t_class,
						teacher: t_teacher,
					};
				}
			}
			for (let i = 0; i < frame_num; i++) {
				units.objects[i] = new Array();
			}
		},

		draw: function (timetables: TimetableType[]) {
			this.init();
			for (let i = 0; i < timetables.length; i++) {
				const obj = timetables[i];
				const id = obj.frame_id;
				units.objects[id].push(obj);
			}

			console.log(units);
			for (let i = 0; i < units.objects.length; i++) {
				if (units.objects[i].length == 0) {
					continue;
				}
				if (units.objects[i].length >= 2) {
					units.elems[i].subject.innerText = String(units.objects[i].length);
					continue;
				}
				console.log(i);
				const obj = units.objects[i][0];
				const elem = units.elems[i];
				elem.subject.innerText = obj.subject_name;
				elem.class.innerText = obj.class_name;
				elem.teacher.innerText = obj.teacher_name;
			}
		},
	};

	const api = {
		get_timetable: async function () {
			console.assert(search_data.classes.length > 0);
			const to_str = function (x: number) {
				return ("00" + String(x)).slice(-2);
			};
			const date_str = `${data.day.getFullYear()}-${to_str(data.day.getMonth() + 1)}-${to_str(data.day.getDate())}`;
			const res = await server.get("timetable/class?id=" + search_data.classes.join("_") + "&day=" + date_str).then(array(timetableDecoder));
			console.log(res);
			timetable.draw(res);
		},
	};

	// react
	type InputButtonProps = {
		callback: (param: { show?: boolean; text?: string }) => void;
	};
	type InputButtonSates = {
		show: boolean;
		text: string;
	};
	class InputButton extends React.Component<InputButtonProps, InputButtonSates> {
		private input_ref = React.createRef<HTMLInputElement>();
		constructor(props: InputButtonProps) {
			super(props);
			this.state = {
				show: false,
				text: "",
			};
			this.handleChange = this.handleChange.bind(this);
			this.changeShow = this.changeShow.bind(this);
		}

		handleChange(event: React.ChangeEvent<HTMLInputElement>) {
			const val = event.currentTarget.value;
			this.setState({
				text: val,
			});
			this.props.callback({ text: val });
		}

		changeShow() {
			this.setState(
				(prev) => {
					return {
						show: !prev.show,
					};
				},
				async () => {
					if (this.state.show) {
						const elem = this.input_ref.current;
						if (elem != null) elem.focus();
					}
					if (!this.state.show) await sleep(1000);
					this.props.callback({ show: this.state.show });
				}
			);
		}

		render(): React.ReactNode {
			if (this.state.show) {
				return <input type="search" value={this.state.text} ref={this.input_ref} onChange={this.handleChange} onBlur={this.changeShow} />;
			} else {
				return <div className="plus_icon" onClick={this.changeShow}></div>;
			}
		}
	}

	interface SimpleData {
		id: number;
		name: string;
	}

	type SearchProps = {
		name: string;
		data: SimpleData[];
		updateSelected: (param: SimpleData[]) => void;
	};
	type SearchState = {
		focus: boolean;
		setting: SimpleData[];
		candidates: SimpleData[];
	};
	class Search extends React.Component<SearchProps, SearchState> {
		constructor(props: SearchProps) {
			super(props);
			const load_str = localStorage.getItem(props.name);
			const load_setting: SimpleData[] = [];
			if (load_str != null && load_str != "") {
				const load_ids = load_str.split(":").map((str) => Number(str));
				for (let i = 0; i < load_ids.length; i++) {
					const id = load_ids[i];
					const elem = props.data.find((d) => d.id === id);
					if (elem != null) load_setting.push(elem);
				}
			}
			this.state = {
				focus: false,
				setting: load_setting,
				candidates: [],
			};
			this.props.updateSelected(load_setting);
			this.addElement = this.addElement.bind(this);
			this.deleteElement = this.deleteElement.bind(this);
			this.changedInput = this.changedInput.bind(this);
			this.save = this.save.bind(this);
		}

		addElement(elem: SimpleData) {
			this.setState(
				function (prevstate) {
					return {
						setting: prevstate.setting.concat([elem]),
					};
				},
				() => {
					this.props.updateSelected(this.state.setting);
				}
			);
		}
		deleteElement(id: Number) {
			this.setState(
				function (prevstate) {
					return {
						setting: prevstate.setting.filter((d) => d.id != id),
					};
				},
				() => {
					this.props.updateSelected(this.state.setting);
				}
			);
		}

		selectedElements() {
			return this.state.setting.map((d) => {
				const onClick = () => this.deleteElement(d.id);
				return (
					<div className="flex" key={d.id.toString()} onClick={onClick}>
						{d.name}
					</div>
				);
				// <div className="delete" onClick={onClick} />
			});
		}

		candidateElements() {
			if (!this.state.focus) return null;
			return this.state.candidates.map((d) => {
				const onClick = () => this.addElement(d);
				return (
					<div className="candidate" key={d.id.toString()} onClick={onClick}>
						{d.name}
					</div>
				);
			});
		}

		search(phrase: string) {
			const score = (name: string) => {
				let s = 0;
				for (let i = 0; i + s < phrase.length; i++) {
					let t = 0;
					for (let j = 0; j < name.length && i + t < phrase.length; j++) {
						if (name[j] == phrase[i + t]) t++;
					}
					s = Math.max(s, t);
				}
				return s;
			};
			const data: SimpleData[] = [];
			// candidates を抜く
			for (let i = 0; i < this.props.data.length; i++) {
				const elem = this.props.data[i];
				if (this.state.candidates.find((c) => c.id === elem.id) == null) {
					data.push({ ...elem });
				}
			}
			data.sort((a, b: SimpleData) => {
				// score で降順 sort
				return score(b.name) - score(a.name);
			});
			this.setState({
				candidates: data.slice(0, Math.min(5, data.length)),
			});
		}

		changedInput(param: { show?: boolean; text?: string }) {
			if (param.text != null) this.search(param.text);
			if (param.show != null) this.setState({ focus: param.show });
		}

		save() {
			const ids_str = this.state.setting.map((d) => d.id.toString()).join(":");
			localStorage.setItem(this.props.name, ids_str);
			console.log(`saved name: ${this.props.name}, value: ${ids_str}`);
		}

		render(): React.ReactNode {
			return (
				<div className="select_class">
					<div className="selected flex">
						<div onClick={this.save}>クラス</div>: {this.selectedElements()}
						<InputButton callback={this.changedInput} />
					</div>
					<div className="candidates flex">{this.candidateElements()}</div>
				</div>
			);
		}
	}

	type TimetableUnitProps = {
		units: TimetableType[];
		onClick: () => void;
	};
	const TimetableUnit = (props: TimetableUnitProps) => {
		if (props.units.length == 1) {
			const u = props.units[0];
			return (
				<td className="unit">
					<table onClick={props.onClick}>
						<tbody>
							<tr>
								<td className="subject">{u.subject_name}</td>
							</tr>
							<tr>
								<td className="class">{u.class_name}</td>
							</tr>
							<tr>
								<td className="teacher">{u.teacher_name}</td>
							</tr>
						</tbody>
					</table>
				</td>
			);
		} else {
			const l = props.units.length;
			return (
				<td className="unit">
					<table>
						<tbody>
							<tr>
								<td className="subject">{l == 0 ? null : l}</td>
							</tr>
							<tr>
								<td className="class">{}</td>
							</tr>
							<tr>
								<td className="teacher">{}</td>
							</tr>
						</tbody>
					</table>
				</td>
			);
		}
	};

	function date2str(date: Date) {
		const to_str = function (x: number) {
			return ("00" + String(x)).slice(-2);
		};
		return `${date.getFullYear()}-${to_str(date.getMonth() + 1)}-${to_str(date.getDate())}`;
	}

	type TimetableProps = {
		classes: InputGraphType;
		date: Date;
	};
	type TimetableState = {
		selected_classes: SimpleData[];
		units: TimetableType[][];
		change_units: TimetableMoveType[];
	};
	class Timetable extends React.Component<TimetableProps, TimetableState> {
		day_weak_strs = ["月", "火", "水", "木", "金", "土"];
		D = 6;
		P = 7;
		constructor(props: TimetableProps) {
			super(props);
			this.state = {
				selected_classes: [],
				units: new Array(this.D * this.P).fill([]),
				change_units: [],
			};
			this.getTimetable = this.getTimetable.bind(this);
			this.setUnits = this.setUnits.bind(this);
		}

		printHeader() {
			return (
				<thead>
					<tr>
						<th key="a"></th>
						{this.day_weak_strs.map((day_str, idx) => {
							const date = new Date(this.props.date);
							let dif = idx + 1 - date.getDay();
							if (dif < 0) {
								dif += 7;
							}
							date.setDate(date.getDate() + dif);
							return (
								<th key={idx} style={{ backgroundColor: dif == 0 ? "#95f9ef" : "#ffffff" }}>
									{date.getMonth() + 1}/{date.getDate()}({day_str})
								</th>
							);
						})}
					</tr>
				</thead>
			);
		}

		printUnits() {
			const table_unit: JSX.Element[] = [];
			const D = this.D;
			const P = this.P;
			for (let i = 0; i < P; i++) {
				const table_row: JSX.Element[] = [];
				for (let j = 0; j < D; j++) {
					const onClick = () => {
						const change_id = this.state.units[j * 7 + i][0].id;
						server
							.get(`timetable/change?change_id=${change_id}&duration_id=1&day=${date2str(this.props.date)}`)
							.then(array(timetableMoveDecoder))
							.then((data) => this.setState({ change_units: data }));
					};
					table_row.push(<TimetableUnit key={j} units={this.state.units[j * 7 + i]} onClick={onClick} />);
				}
				table_unit.push(
					<tr key={i}>
						<th>{`${i + 1}限`}</th>
						{table_row}
					</tr>
				);
			}

			return <tbody>{table_unit}</tbody>;
		}

		printChanges() {
			const changes = this.state.change_units;
			return (
				<div className="changes">
					{changes.map((d) => {
						const tim = d.timetable;
						return (
							<div>
								クラス名: {tim.class_name}, 教科: {tim.subject_name}, 先生: {tim.teacher_name}
								日時: {date2str(tim.day)} {tim.frame_period} 限 変更日時: {date2str(d.day)} {d.frame_id / 7}限
							</div>
						);
					})}
				</div>
			);
		}

		setUnits(timetables: TimetableType[]) {
			const D = this.D;
			const P = this.P;
			let units: TimetableType[][] = new Array(D * P);
			for (let i = 0; i < D * P; i++) {
				units[i] = new Array();
			}
			const lim = new Date(this.props.date);
			lim.setDate(lim.getDate() + 7);
			lim.setHours(0, 0, 0, 0);
			for (let i = 0; i < timetables.length; i++) {
				const obj = timetables[i];
				if (obj.day < lim) units[obj.frame_id].push(obj);
			}
			this.setState({ units: units });
		}
		getTimetable(selected_classes: SimpleData[]) {
			if (selected_classes.length == 0) {
				this.setUnits([]);
				return;
			}
			server
				.get(
					"timetable/class?id=" +
						selected_classes
							.map((d) => {
								return d.id.toString();
							})
							.join("_") +
						"&day=" +
						date2str(this.props.date)
				)
				.then(array(timetableDecoder))
				.then(this.setUnits);
		}

		render(): React.ReactNode {
			return (
				<div>
					<div id="timetable_info">
						<Search name="select_classes" data={this.props.classes.nodes} updateSelected={this.getTimetable} />
					</div>
					<div id="timetable">
						<table>
							{this.printHeader()}
							{this.printUnits()}
						</table>
					</div>
					{this.printChanges()}
				</div>
			);
		}
	}

	window.addEventListener("load", async function (e) {
		// api.get_timetable();
		const classes = await getClass();
		const container = document.getElementById("root");
		if (container == null) {
			console.assert(`container is null: ${container}`);
			return;
		}
		const root = createRoot(container);
		root.render(<Timetable classes={classes} date={new Date(2021, 4, 1)} />);
	});
})();
