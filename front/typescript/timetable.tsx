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

export type TimetableUnitProps = {
	units: TimetableType[];
	color1?: string;
	color2?: string;
	onClick?: () => void;
};

export const TimetableUnit = (props: TimetableUnitProps) => {
	const style1: React.CSSProperties = {
		backgroundColor: props.color1,
	};
	const style2: React.CSSProperties = {
		backgroundColor: props.color2,
	};
	let values: string[] = [];
	if (props.units.length == 1) {
		const u = props.units[0];
		values = [u.subject_name, u.class_name, u.teacher_name.join(",")];
	} else {
		const l = props.units.length;
		values = [l == 0 ? "" : String(l), "", ""];
	}
	return (
		<table onClick={props.onClick}>
			<tbody>
				<tr>
					<td className="subject" style={style1}>
						{values[0]}
					</td>
				</tr>
				<tr>
					<td className="class" style={style2}>
						{values[1]}
					</td>
				</tr>
				<tr>
					<td className="teacher" style={style2}>
						{values[2]}
					</td>
				</tr>
			</tbody>
		</table>
	);
};

export function date2str(date: Date) {
	const to_str = function (x: number) {
		return ("00" + String(x)).slice(-2);
	};
	return `${date.getFullYear()}-${to_str(date.getMonth() + 1)}-${to_str(date.getDate())}`;
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
