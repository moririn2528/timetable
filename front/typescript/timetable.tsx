import { server, InputGraphType, sleep, jsdate } from "./api";
import { decodeType, record, number, string, array, undef } from "typescript-json-decoder";
import React, { useState } from "react";
import { createRoot } from "react-dom/client";

export type TimetableType = decodeType<typeof timetableDecoder>;
export const timetableDecoder = record({
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
	teacher_id: array(number),
	teacher_name: array(string),
	place_id: number,
	day: jsdate,
});
export type TimetableMoveType = decodeType<typeof timetableMoveDecoder>;
export const timetableMoveDecoder = record({
	timetable: timetableDecoder,
	day: jsdate,
	frame_id: number,
});

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
			if (this.state.setting.find((c) => c.id === elem.id) == null) {
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

export type TimetableUnitProps = {
	units: TimetableType[];
	onClick: () => void;
};

export const TimetableUnit = (props: TimetableUnitProps) => {
	if (props.units.length == 1) {
		const u = props.units[0];
		return (
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
		);
	} else {
		const l = props.units.length;
		return (
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
		);
	}
};

export function date2str(date: Date) {
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
	change_error?: string;
};

export class Timetable extends React.Component<TimetableProps, TimetableState> {
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
		this.setUnits = this.setUnits.bind(this);
	}

	getDaystr(day: number) {
		if (day == 0) return "日";
		return this.day_weak_strs[day - 1];
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

	printUnit(d: number, p: number) {
		// day, period
		const onClick = () => {
			const change_id = this.state.units[d * 7 + p][0].id;
			server
				.get(`timetable/change?change_id=${change_id}&duration_id=1&day=${date2str(this.props.date)}`)
				.then(array(timetableMoveDecoder))
				.then((data) => this.setState({ change_units: data, change_error: undefined }))
				.catch((e) => {
					this.setState({ change_error: String(e) });
				});
		};
		return <TimetableUnit key={d} units={this.state.units[d * 7 + p]} onClick={onClick} />;
	}

	printUnits() {
		const table_unit: JSX.Element[] = [];
		const D = this.D;
		const P = this.P;
		for (let i = 0; i < P; i++) {
			const table_row: JSX.Element[] = [];
			for (let j = 0; j < D; j++) {
				table_row.push(
					<td className="unit" key={j}>
						{this.printUnit(j, i)}
					</td>
				);
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
		if (this.state.change_error != null) {
			return <div className="changes">サーバーエラー: {this.state.change_error}</div>;
		}
		return (
			<div className="changes">
				{changes.map((d) => {
					const tim = d.timetable;
					return (
						<div key={d.timetable.id}>
							クラス名: {tim.class_name}, 教科: {tim.subject_name}, 先生: {tim.teacher_name.join(",")} 日時: {date2str(tim.day)}({this.getDaystr(tim.day.getDay())}){" "}
							{tim.frame_period + 1} 限 変更日時: {date2str(d.day)}({this.getDaystr(d.day.getDay())}) {(d.frame_id % 7) + 1}限
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

	componentDidUpdate(prevprops: TimetableProps) {
		if (this.props.date !== prevprops.date) {
			console.log(this.state.selected_classes);
			this.getTimetable(this.state.selected_classes);
		}
	}

	render(): React.ReactNode {
		const updateSelected = (data: SimpleData[]) => {
			this.setState({
				selected_classes: data,
			});
			this.getTimetable(data);
		};
		return (
			<div>
				<div id="timetable_info">
					<Search name="select_classes" data={this.props.classes.nodes} updateSelected={updateSelected} />
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

export const TimetablePreviewDate = (props: { today: Date; date: Date; setDate: (param: Date) => void }) => {
	const setNextWeek = () => {
		const d = new Date(props.date);
		d.setDate(d.getDate() + 7 - d.getDay());
		props.setDate(d);
	};
	const setLastWeek = () => {
		let d = new Date(props.date);
		d.setDate(d.getDate() - 7);
		if (d < props.today) d = new Date(props.today);
		props.setDate(d);
	};
	return (
		<div className="flex">
			<div onClick={setLastWeek}>{"先週<"}</div>
			<div onClick={setNextWeek}>{">来週"}</div>
		</div>
	);
};
