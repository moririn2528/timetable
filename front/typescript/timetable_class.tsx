import { getClass, InputGraphType } from "./api";
import React from "react";
import { createRoot } from "react-dom/client";
import { TimetablePreviewDate, Timetable } from "./timetable";

(function () {
	const TimetableWithDate = (props: { classes: InputGraphType; date: Date }) => {
		const [date, setDate] = React.useState(props.date);
		return (
			<div>
				<TimetablePreviewDate today={props.date} date={date} setDate={setDate} />
				<Timetable classes={props.classes} date={date} />
			</div>
		);
	};

	window.addEventListener("load", async function (e) {
		// api.get_timetable();
		const classes = await getClass();
		const container = document.getElementById("root");
		if (container == null) {
			console.assert(`container is null: ${container}`);
			return;
		}
		const root = createRoot(container);
		root.render(<TimetableWithDate classes={classes} date={new Date(2021, 3, 12)} />);
	});
})();
