import { server, jsdate, TeacherType } from "./api";
import { decodeType, record, number, string, array, undef, dict } from "typescript-json-decoder";
import React from "react";
import { TimetableType, TimetableMoveType, date2str, timetableDecoder, timetableMoveDecoder, TimetableUnit, TimetableUnitProps, TimetablePreviewDate, TimetableStyle, TimetableUnitStyle } from "./timetable";

const PERIOD = 7; // 1日の最大コマ数
type BanUnit = {
	date: Date;
	period: number;
};

type TeacherAvoidType = decodeType<typeof teacherAvoidDecoder>;
const teacherAvoidDecoder = record({
	day: jsdate,
	avoid: array(number),
});

type ChangeAvoidType = {
	day: Date;
	frame: number;
	avoid: number;
}

type TeacherTimetableUnitProps = TimetableUnitProps & {
	avoid: number;
	color?: string;
};

const TeacherTimetableUnit = (props: TeacherTimetableUnitProps) => {
	let col = "#ffffff";
	switch (props.avoid) {
		case 1:
			col = "##ffffc8";
			break;
		case 2:
			col = "#ffff82";
			break;
		case 3:
			col = "#ffcc82";
			break;
		case 4:
			col = "#ff9682";
			break;
		case 5:
			col = "#ff7882";
			break;
		case 7:
			col = "#b6e1ff";
			break;
		case 8:
			col = "#777777";
			break;
		case 9:
			col = "#cccccc";
			break;
	}

	return <TimetableUnit units={props.units} onClick={props.onClick} color1={props.color} color2={col} />;
};

const SearchChange = (props: { teacher_id: number; select_units: BanUnit[]; date: Date; clear_units: () => void; display_changes: (param: TimetableMoveType[]) => void }) => {
	const day_weak_strs = ["月", "火", "水", "木", "金", "土"];
	const [change_units, setChangeUnits] = React.useState<TimetableMoveType[]>([]);
	const [change_avoids, setChangeAvoids] = React.useState<ChangeAvoidType[]>([]);
	const [change_error, setChangeError] = React.useState<string | undefined>();
	const [on_load, setOnLoad] = React.useState<Boolean>(false);
	const [changed, setChangedFlag] = React.useState<Boolean>(false);

	const calc = () => {
		setOnLoad(true);
		server
			.get(
				`timetable/change?duration_id=1&teacher_id=${props.teacher_id}&day=${date2str(props.date)}&ban_units=${props.select_units
					.map((val) => date2str(val.date) + "A" + String(val.period))
					.join("A")}`
			)
			.then(array(timetableMoveDecoder))
			.then((data) => {
				setChangeUnits(data);
				setChangeError(undefined);
				setOnLoad(false);
				setChangedFlag(false);
			})
			.catch((e) => {
				setChangeUnits([]);
				setChangeError(String(e));
				setOnLoad(false);
			});
	};

	const buttonDisplay = () => {
		props.display_changes(change_units);
	};
	const changeTimetable = () => {
		setChangedFlag((prev) => {
			if (prev) {
				return true;
			}
			server
				.post("timetable/change", change_units)
				.then(() => {
					console.log("変更されました");
				})
				.catch((e) => {
					console.log("変更エラー", e);
				});
			return true;
		});
	};

	return (
		<div>
			{!on_load && props.select_units.length > 0 ? (
				<div>
					<button onClick={calc}>計算</button>
					<button onClick={props.clear_units}>クリア</button>
					<button onClick={buttonDisplay}>描画</button>
					<button onClick={changeTimetable}>変更</button>
				</div>
			) : (
				<div>
					<button onClick={buttonDisplay}>描画</button>
				</div>
			)}
			{change_error != undefined ? "エラー: " + change_error : ""}
			<br />
			結果: {on_load ? "計算中" : ""}
			<br />
			{change_units.map((val, index) => {
				const tim = val.timetable;
				return (
					<div key={index}>
						クラス名: {tim.class_name}, 教科: {tim.subject_name}, 先生: {tim.teacher_name.join(",")} 日時: {date2str(tim.day)}({day_weak_strs[tim.day.getDay() - 1]}){" "}
						{(tim.frame_id % PERIOD) + 1} 限 変更日時: {date2str(val.day)}({day_weak_strs[val.day.getDay() - 1]}) {(val.frame_id % PERIOD) + 1}限
					</div>
				);
			})}
		</div>
	);
};

type TimetableProps = {
	teacher: number;
	date: Date;
};
type TimetableState = {
	units: TimetableType[][];
	avoids: number[];
	selected_units: BanUnit[];
	change_units: TimetableMoveType[];
	change_avoids: ChangeAvoidType[];
	changed: boolean;
};

class TeacherTimetable extends React.Component<TimetableProps, TimetableState> {
	static day_weak_strs = ["月", "火", "水", "木", "金", "土"];
	static D = 6;
	static P = PERIOD;

	constructor(props: TimetableProps) {
		super(props);
		const D = TeacherTimetable.D;
		const P = TeacherTimetable.P;
		this.state = {
			units: new Array(D * P).fill([]),
			avoids: new Array(D * P).fill(0),
			selected_units: [],
			change_units: [],
			change_avoids: [],
			changed: false,
		};
		this.setUnits = this.setUnits.bind(this);
		this.setAvoids = this.setAvoids.bind(this);
		this.printUnit = this.printUnit.bind(this);
		this.getTimetable();
	}

	getDaystr(day: number) {
		if (day == 0) return "日";
		return TeacherTimetable.day_weak_strs[day - 1];
	}

	printHeader() {
		return (
			<thead>
				<tr>
					<th key="a"></th>
					{TeacherTimetable.day_weak_strs.map((day_str, idx) => {
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

	printUnit(d: number, p: number, change_from: boolean, change_to: boolean) {
		// day, period
		const date = new Date(this.props.date);
		let dif = d + 1 - date.getDay();
		if (dif < 0) dif += 7;
		date.setDate(date.getDate() + dif);
		const onClick = () => {
			this.setState((prevstate) => {
				const select_units = [...prevstate.selected_units];
				const unit: BanUnit = {
					date: date,
					period: p,
				};
				const idx = select_units.findIndex((u) => {
					return u.date.getTime() == date.getTime() && u.period == p;
				});
				if (idx != -1) {
					// erase select_units[idx]
					const m = select_units.length - 1;
					[select_units[idx], select_units[m]] = [select_units[m], select_units[idx]];
					select_units.pop();
				} else {
					select_units.push(unit);
				}
				return {
					selected_units: select_units,
				};
			});
		};
		let color: string | undefined = undefined;
		if (change_from && change_to) {
			color = "#fcf8b3";
		} else if (change_from) {
			color = "#ffbeba";
		} else if (change_to) {
			color = "#b5e1ff";
		} else if (this.state.selected_units.findIndex((val) => val.date.getTime() == date.getTime() && val.period == p) != -1) {
			color = "#baefb3";
		}

		return <TeacherTimetableUnit key={d} units={this.state.units[d * 7 + p]} avoid={this.state.avoids[d * 7 + p]} onClick={onClick} color={color} />;
	}

	getDayIndex(day: Date) {
		const D = TeacherTimetable.D;
		let d = (day.getTime() - this.props.date.getTime()) / (24 * 60 * 60 * 1000);
		if (d < 0 || 7 <= d) return -1;
		d += this.props.date.getDay() - 1;
		if (d >= D) d -= D;
		return d;
	}

	printUnits() {
		const table_unit: JSX.Element[] = [];
		const D = TeacherTimetable.D;
		const P = TeacherTimetable.P;
		const start_date = new Date(this.props.date);
		start_date.setDate(start_date.getDate() - start_date.getDay() + 1);
		const change_from: number[] = [];
		const change_to: number[] = [];

		for (let i = 0; i < this.state.change_units.length; i++) {
			const u = this.state.change_units[i];
			if (!u.timetable.teacher_id.includes(this.props.teacher)) continue;
			let d = this.getDayIndex(u.timetable.day);
			let d2 = u.timetable.day.getDay() - 1;
			let p = u.timetable.frame_id % P;
			if (0 <= d && d < D) change_from.push(d2 * P + p);
			d = this.getDayIndex(u.day);
			d2 = u.day.getDay() - 1;
			p = u.frame_id % P;
			if (0 <= d && d < D) change_to.push(d2 * P + p);
		}
		for (let i = 0; i < P; i++) {
			const table_row: JSX.Element[] = [];
			for (let j = 0; j < D; j++) {
				table_row.push(
					<td className="unit" key={j} style={TimetableUnitStyle}>
						{this.printUnit(j, i, change_from.includes(j * P + i), change_to.includes(j * P + i))}
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

	setUnits(timetables: TimetableType[]) {
		const D = TeacherTimetable.D;
		const P = TeacherTimetable.P;
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
	setAvoids(in_avoids: TeacherAvoidType[]) {
		const D = TeacherTimetable.D;
		const P = TeacherTimetable.P;
		let avoids: number[] = new Array(D * P).fill(0);
		const lim = new Date(this.props.date);
		lim.setDate(lim.getDate() + 7);
		lim.setHours(0, 0, 0, 0);
		for (let i = 0; i < in_avoids.length; i++) {
			const av_day = in_avoids[i];
			const d = av_day.day.getDay() - 1;
			console.assert(d != -1);
			console.assert(this.props.date <= av_day.day && av_day.day < lim);
			for (let j = 0; j < P; j++) {
				avoids[d * P + j] = av_day.avoid[j];
			}
		}
		this.setState({ avoids: avoids });
	}
	getTimetable() {
		server
			.get("timetable/teacher?id=" + String(this.props.teacher) + "&day=" + date2str(this.props.date))
			.then(array(timetableDecoder))
			.then(this.setUnits);
		server
			.get("teacher/avoid?id=" + String(this.props.teacher) + "&day=" + date2str(this.props.date))
			.then(array(teacherAvoidDecoder))
			.then(this.setAvoids);
	}

	componentDidUpdate(prevprops: TimetableProps) {
		if (this.props.date !== prevprops.date) {
			this.getTimetable();
		}
		if (this.props.teacher !== prevprops.teacher) {
			this.getTimetable();
		}
	}

	render(): React.ReactNode {
		const clear_units = () =>
			this.setState({
				selected_units: [],
				changed: false,
			});
		const display_changes = (move: TimetableMoveType[]) => {
			this.setState((prev) => {
				const bef = prev.change_units;
				if (JSON.stringify(bef) == JSON.stringify(move)) {
					return {
						change_units: [],
					};
				} else {
					return {
						change_units: move,
					};
				}
			});
		};
		return (
			<div>
				<div id="timetable">
					<table style={TimetableStyle}>
						{this.printHeader()}
						{this.printUnits()}
					</table>
				</div>
				<SearchChange teacher_id={this.props.teacher} date={this.props.date} select_units={this.state.selected_units} clear_units={clear_units} display_changes={display_changes} />
			</div>
		);
	}
}

export const TeacherTimetableWithDate = (props: { teachers: TeacherType[]; date: Date }) => {
	const [date, setDate] = React.useState(props.date);
	const [teacher, setTeacher] = React.useState(props.teachers[0].id);
	const onChange: React.ChangeEventHandler<HTMLSelectElement> = (e) => {
		setTeacher(parseInt(e.currentTarget.value, 10));
	};
	return (
		<div>
			<TimetablePreviewDate today={props.date} date={date} setDate={setDate} />
			先生:
			<select name="teacher" onChange={onChange}>
				{props.teachers.map((teacher, index) => {
					return (
						<option value={teacher.id} key={index}>
							{teacher.name}
						</option>
					);
				})}
			</select>
			<TeacherTimetable teacher={teacher} date={date} />
		</div>
	);
};
