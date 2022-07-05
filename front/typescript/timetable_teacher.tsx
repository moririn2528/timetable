import { server, jsdate } from "./api";
import { decodeType, record, number, string, array, undef } from "typescript-json-decoder";
import React, { useState } from "react";
import { createRoot } from "react-dom/client";
import { TimetableType, TimetableMoveType, date2str, timetableDecoder, timetableMoveDecoder, TimetableUnit, TimetableUnitProps, TimetablePreviewDate } from "./timetable";

(function () {
	const PERIOD = 7; // 1日の最大コマ数
	type BanUnit = {
		date: Date;
		period: number;
	};

	type TeacherType = decodeType<typeof inputTeacherDecoder>;
	const inputTeacherDecoder = record({
		id: number,
		name: string,
	});
	async function getTeachers() {
		return server.get("teacher").then(array(inputTeacherDecoder));
	}
	type TeacherAvoidType = decodeType<typeof teacherAvoidDecoder>;
	const teacherAvoidDecoder = record({
		day: jsdate,
		avoid: array(number),
	});

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

	const SearchChange = (props: { teacher_id: number; select_units: BanUnit[]; date: Date; clear_units: () => void }) => {
		const day_weak_strs = ["月", "火", "水", "木", "金", "土"];
		const [change_units, setChangeUnits] = React.useState<TimetableMoveType[]>([]);
		const [change_error, setChangeError] = React.useState<string | undefined>();
		const [on_load, setOnLoad] = React.useState<Boolean>(false);

		const calc = () => {
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
				})
				.catch((e) => {
					setChangeUnits([]);
					setChangeError(String(e));
					setOnLoad(false);
				});
		};

		return (
			<div>
				{!on_load && props.select_units.length > 0 ? (
					<div>
						<button onClick={calc}>計算</button>
						<button onClick={props.clear_units}>クリア</button>
					</div>
				) : (
					""
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

		printUnit(d: number, p: number) {
			// day, period
			const date = new Date(this.props.date);
			date.setDate(date.getDate() - date.getDay() + d + 1);
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

			return (
				<TeacherTimetableUnit
					key={d}
					units={this.state.units[d * 7 + p]}
					avoid={this.state.avoids[d * 7 + p]}
					onClick={onClick}
					color={this.state.selected_units.findIndex((val) => val.date.getTime() == date.getTime() && val.period == p) != -1 ? "#baefb3" : undefined}
				/>
			);
		}

		printUnits() {
			const table_unit: JSX.Element[] = [];
			const D = TeacherTimetable.D;
			const P = TeacherTimetable.P;
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
			return (
				<div>
					<div id="timetable">
						<table>
							{this.printHeader()}
							{this.printUnits()}
						</table>
					</div>
					<SearchChange
						teacher_id={this.props.teacher}
						date={this.props.date}
						select_units={this.state.selected_units}
						clear_units={() =>
							this.setState({
								selected_units: [],
							})
						}
					/>
				</div>
			);
		}
	}

	const TimetableWithDate = (props: { teachers: TeacherType[]; date: Date }) => {
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

	window.addEventListener("load", async function (e) {
		// api.get_timetable();

		const teachers = await getTeachers();
		const container = document.getElementById("root");
		if (container == null) {
			console.assert(`container is null: ${container}`);
			return;
		}
		const root = createRoot(container);
		root.render(<TimetableWithDate teachers={teachers} date={new Date(2021, 3, 13)} />);
	});
})();
