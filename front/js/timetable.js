const server = {
	get: async function (path) {
		const port = location.port;
		console.log(port);
		const res = await fetch("http://localhost:" + port + "/api/" + path);
		if (res.status != 200) {
			console.error(res);
			return;
		}
		return res.json();
	},
};
const sleep = (waitTime) => new Promise((resolve) => setTimeout(resolve, waitTime));

(function () {
	let search_data = {
		classes: [0],
	};
	let data = {
		day: new Date(2021, 3, 1),
		classes: [],
	};
	let units = {
		elems: [],
		objects: [], // [i][j]: frame_id=i, j 個目の object
	};

	const search = {
		draw_candidates: function () {
			const input = document.getElementById("");
			// 検索結果を格納するための配列を用意
			searchResult = [];

			// 検索結果エリアの表示を空にする
			$("#search-result__list").empty();
			$(".search-result__hit-num").empty();

			// 検索ボックスに値が入ってる場合
			if (searchText != "") {
				$(".target-area li").each(function () {
					targetText = $(this).text();

					// 検索対象となるリストに入力された文字列が存在するかどうかを判断
					if (targetText.indexOf(searchText) != -1) {
						// 存在する場合はそのリストのテキストを用意した配列に格納
						searchResult.push(targetText);
					}
				});

				// 検索結果をページに出力
				for (var i = 0; i < searchResult.length; i++) {
					$("<span>").text(searchResult[i]).appendTo("#search-result__list");
				}

				// ヒットの件数をページに出力
				hitNum = "<span>検索結果</span>:" + searchResult.length + "件見つかりました。";
				$(".search-result__hit-num").append(hitNum);
			}
		},
	};

	const timetable = {
		init: function () {
			const frame_num = 42;
			units.elems = new Array(frame_num);
			units.objects = new Array(frame_num);
			day_weak_strs = ["月", "火", "水", "木", "金", "土"];
			const tt = document.getElementById("timetable");
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

			const add_t_col = function (t) {
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

		draw: function (timetables) {
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
					units.elems[i].subject.innerText = units.objects[i].length;
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
			to_str = function (x) {
				return ("00" + String(x)).slice(-2);
			};
			date_str = `${1900 + data.day.getYear()}-${to_str(data.day.getMonth() + 1)}-${to_str(data.day.getDate())}`;
			const res = await server.get("timetable/class?id=" + search_data.classes.join("_") + "&day=" + date_str);
			console.log(res);
			timetable.draw(res);
		},
	};

	window.addEventListener("load", async function (e) {
		data.day = new Date(2021, 3, 1);
		api.get_timetable();
	});
})();
