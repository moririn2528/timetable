import { getClass, InputGraphType, sleep, server } from "./api";
import React from "react";
import { createRoot } from "react-dom/client";
import { TimetableType, date2str, TimetableUnit, timetableDecoder, TimetablePreviewDate, TimetableStyle, TimetableUnitStyle } from "./timetable";
import { array } from "typescript-json-decoder";

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
			return <div onClick={this.changeShow} style={{
				width: 20,
				height: 20,
				background: "url(./images/icon_001510_256.png)",
				borderRadius: "50%",
				backgroundSize: "cover",
			}}></div>;
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
				<div key={d.id.toString()} onClick={onClick}>
					{d.name}
				</div>
			);
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
			<div style={{
				display: "flex",
			}}>
				<div style={{
					display: "flex",
				}}>
					<div onClick={this.save}>クラス</div>: {this.selectedElements()}
					<InputButton callback={this.changedInput} />
				</div>
				<div style={{
					display: "flex",
				}}>{this.candidateElements()}</div>
			</div>
		);
	}
}

type TimetableProps = {
	classes: InputGraphType;
	date: Date;
};
type TimetableState = {
	selected_classes: SimpleData[];
	units: TimetableType[][];
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
		return <TimetableUnit key={d} units={this.state.units[d * 7 + p]} />;
	}

	printUnits() {
		const table_unit: JSX.Element[] = [];
		const D = this.D;
		const P = this.P;
		for (let i = 0; i < P; i++) {
			const table_row: JSX.Element[] = [];
			for (let j = 0; j < D; j++) {
				table_row.push(
					<td className="unit" key={j} style={TimetableUnitStyle}>
						{this.printUnit(j, i)}
					</ td>
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
					<table style={TimetableStyle}>
						{this.printHeader()}
						{this.printUnits()}
					</table>
				</div>
			</div>
		);
	}
}

export const ClassTimetableWithDate = (props: { classes: InputGraphType; date: Date }) => {
	const [date, setDate] = React.useState(props.date);
	return (
		<div>
			<TimetablePreviewDate today={props.date} date={date} setDate={setDate} />
			<Timetable classes={props.classes} date={date} />
		</div>
	);
};


