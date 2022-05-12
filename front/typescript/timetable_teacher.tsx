import { server, jsdate } from "./api";
import { decodeType, record, number, string, array, undef } from "typescript-json-decoder";
import React, { useState } from "react";
import { createRoot } from "react-dom/client";
import { TimetableType, TimetableMoveType, date2str, timetableDecoder, timetableMoveDecoder, TimetableUnit, TimetableUnitProps, TimetablePreviewDate } from "./timetable";

(function () {
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

		return (
			<div
				style={{
					background: col,
				}}
			>
				<TimetableUnit units={props.units} onClick={props.onClick} />
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
		change_units: TimetableMoveType[];
		change_error?: string;
	};

	class TeacherTimetable extends React.Component<TimetableProps, TimetableState> {
		static day_weak_strs = ["月", "火", "水", "木", "金", "土"];
		static D = 6;
		static P = 7;

		constructor(props: TimetableProps) {
			super(props);
			const D = TeacherTimetable.D;
			const P = TeacherTimetable.P;
			this.state = {
				units: new Array(D * P).fill([]),
				change_units: [],
				avoids: new Array(D * P).fill(0),
			};
			this.setUnits = this.setUnits.bind(this);
			this.setAvoids = this.setAvoids.bind(this);
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
			return <TeacherTimetableUnit key={d} units={this.state.units[d * 7 + p]} avoid={this.state.avoids[d * 7 + p]} onClick={onClick} />;
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
					{this.printChanges()}
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
		root.render(<TimetableWithDate teachers={teachers} date={new Date(2021, 3, 12)} />);
	});
})();
