import { server } from "./api";
import { decodeType, record, number, string, array } from "typescript-json-decoder";

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

	// const search = {
	// 	draw_candidates: function () {
	// 		const input = document.getElementById("");
	// 		// 検索結果を格納するための配列を用意
	// 		searchResult = [];

	// 		// 検索結果エリアの表示を空にする
	// 		$("#search-result__list").empty();
	// 		$(".search-result__hit-num").empty();

	// 		// 検索ボックスに値が入ってる場合
	// 		if (searchText != "") {
	// 			$(".target-area li").each(function () {
	// 				targetText = $(this).text();

	// 				// 検索対象となるリストに入力された文字列が存在するかどうかを判断
	// 				if (targetText.indexOf(searchText) != -1) {
	// 					// 存在する場合はそのリストのテキストを用意した配列に格納
	// 					searchResult.push(targetText);
	// 				}
	// 			});

	// 			// 検索結果をページに出力
	// 			for (var i = 0; i < searchResult.length; i++) {
	// 				$("<span>").text(searchResult[i]).appendTo("#search-result__list");
	// 			}

	// 			// ヒットの件数をページに出力
	// 			hitNum = "<span>検索結果</span>:" + searchResult.length + "件見つかりました。";
	// 			$(".search-result__hit-num").append(hitNum);
	// 		}
	// 	},
	// };

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
			console.log(await server.get("timetable/class?id=" + search_data.classes.join("_") + "&day=" + date_str));
			const res = await server.get("timetable/class?id=" + search_data.classes.join("_") + "&day=" + date_str).then(array(timetableDecoder));
			console.log(res);
			timetable.draw(res);
		},
	};

	window.addEventListener("load", async function (e) {
		api.get_timetable();
	});
})();
