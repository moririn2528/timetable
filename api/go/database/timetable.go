package database

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"timetable/errors"
	"timetable/usecase"
)

type DatabaseAny struct {
}
type DatabaseTimetable struct {
}

func (*DatabaseAny) GetDurationId(date time.Time) (int, error) {
	date_str := date.Format("2006-01-02")
	row := db.QueryRow(
		"SELECT id FROM duration WHERE start_date <= ? AND ? <= end_date",
		date_str, date_str,
	)
	if row.Err() != nil {
		return -1, errors.ErrorWrap(row.Err())
	}
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, errors.ErrorWrap(err)
	}
	return id, nil
}

func (*DatabaseTimetable) GetNomalTimetable(duration_id int, class_ids []int, teacher_id int) ([]usecase.NormalTimetable, error) {
	// class_id, teacher_id は -1 を許す
	var res []usecase.NormalTimetable

	list2str := func(l []int) string {
		var res []string
		for _, v := range l {
			res = append(res, strconv.Itoa(v))
		}
		return strings.Join(res, ",")
	}
	// 検索条件
	conditions := []string{fmt.Sprintf("duration_id = %d", duration_id)}
	if len(class_ids) > 0 {
		conditions = append(conditions, fmt.Sprintf(
			"class_id in (%s)", list2str(class_ids),
		))
	}
	if teacher_id != -1 {
		conditions = append(conditions, fmt.Sprintf(
			"(teacher_id = %d OR add_teacher_id = %d)", teacher_id, teacher_id,
		))
	}

	// sql
	rows, err := db.Query(strings.Join([]string{
		"SELECT tb.id, c.id, c.name, d.id, d.name, tb.frame_id,",
		"s.id, s.name, t.id, t.name, t2.id, t2.name",
		"FROM (SELECT * FROM normal_timetable WHERE",
		strings.Join(conditions, " AND "),
		") AS tb",
		"LEFT JOIN classroom AS c ON tb.class_id = c.id",
		"LEFT JOIN duration AS d ON tb.duration_id = d.id",
		"LEFT JOIN subject AS s ON tb.subject_id = s.id",
		"LEFT JOIN teacher AS t ON tb.teacher_id = t.id",
		"LEFT JOIN teacher AS t2 ON tb.add_teacher_id = t2.id",
	}, " "))
	if err != nil {
		return res, errors.ErrorWrap(err)
	}
	var res_temp []usecase.NormalTimetable
	for rows.Next() {
		var t usecase.NormalTimetable
		var teach_id int
		var teach_name string
		var teach2_id sql.NullInt64
		var teach2_name sql.NullString
		err := rows.Scan(&t.Id, &t.ClassId, &t.ClassName, &t.DurationId,
			&t.DurationName, &t.FrameId,
			&t.SubjectId, &t.SubjectName, &teach_id, &teach_name, &teach2_id, &teach2_name)
		t.FrameDayWeek = t.FrameId/usecase.PERIOD + 1
		t.FramePeriod = t.FrameId % usecase.PERIOD
		t.TeacherIds = []int{teach_id}
		t.TeacherNames = []string{teach_name}
		if teach2_id.Valid {
			if !teach2_name.Valid {
				return nil, errors.NewError("addtional teacher name is null")
			}
			t.TeacherIds = append(t.TeacherIds, int(teach2_id.Int64))
			t.TeacherNames = append(t.TeacherNames, teach2_name.String)
		}
		if err != nil {
			return res, errors.ErrorWrap(err)
		}
		res_temp = append(res_temp, t)
	}

	// 重複取り除く
	// duration, frame 同じかつ (クラス同じか教師同じ) のとき重複
	type key1 struct {
		c, d, f int
	}
	type key2 struct {
		d, f int
		t    [2]int
	}
	m1 := make(map[key1]struct{})
	m2 := make(map[key2]struct{})
	for i := len(res_temp) - 1; i >= 0; i-- {
		t := res_temp[i]
		k1 := key1{t.ClassId, t.DurationId, t.FrameId}
		k2 := key2{t.DurationId, t.FrameId, [2]int{t.TeacherIds[0], -1}}
		if len(t.TeacherIds) > 1 {
			k2.t[1] = t.TeacherIds[1]
		}
		if k2.t[0] < k2.t[1] {
			k2.t[0], k2.t[1] = k2.t[1], k2.t[0]
		}
		_, exist1 := m1[k1]
		_, exist2 := m2[k2]
		if exist1 || exist2 {
			continue
		}
		res = append(res, t)
		m1[k1] = struct{}{}
		m2[k2] = struct{}{}
	}
	return res, nil
}

// GetTimetable:
// class_ids は空を許す、その時は条件として入れない
// teacher_id は -1 を許す、その時は条件として入れない
func (dc *DatabaseTimetable) GetTimetable(
	duration_id int, class_ids []int, teacher_id int,
	start_day time.Time, end_day time.Time,
) ([]usecase.Timetable, error) {
	var res []usecase.Timetable

	// 通常時の時間割取得
	normal_timetable, err := dc.GetNomalTimetable(duration_id, class_ids, teacher_id)
	if err != nil {
		return res, errors.ErrorWrap(err)
	}
	sort.Slice(normal_timetable, func(i int, j int) bool {
		return normal_timetable[i].FrameId < normal_timetable[j].FrameId
	})
	// 曜日で分割
	var normal [7][]usecase.NormalTimetable // id sorted
	j := 0
	for _, v := range normal_timetable {
		for ; j < usecase.FRAME_NUM; j++ {
			if v.FrameId == j {
				break
			}
		}
		if j >= usecase.FRAME_NUM {
			return res, errors.NewError("index error")
		}
		day_week := j/usecase.PERIOD + 1 // 曜日
		normal[day_week] = append(normal[day_week], v)
	}
	start_day = time.Date(
		start_day.Year(), start_day.Month(), start_day.Day(),
		0, 0, 0, 0, start_day.Location(),
	)

	// 追加時間割を取得
	// sql 条件生成
	list2str := func(l []int) string {
		var res []string
		for _, v := range l {
			res = append(res, strconv.Itoa(v))
		}
		return strings.Join(res, ",")
	}
	conditions := []string{fmt.Sprintf("duration_id = %d", duration_id)}
	if len(class_ids) > 0 {
		conditions = append(conditions, fmt.Sprintf(
			"class_id in (%s)", list2str(class_ids),
		))
	}
	if teacher_id != -1 {
		conditions = append(conditions, fmt.Sprintf(
			"(teacher_id = %d OR add_teacher_id = %d)", teacher_id, teacher_id,
		))
	}
	conditions = append(conditions,
		fmt.Sprintf("\"%s\" <= day", start_day.Format("2006-01-02")),
		fmt.Sprintf("day <= \"%s\"", end_day.Format("2006-01-02")),
	)
	// sql
	rows, err := db.Query(strings.Join([]string{
		"SELECT tb.id, c.id, c.name, d.id, d.name, tb.frame_id,",
		"s.id, s.name, t.id, t.name, tb.day, t2.id, t2.name",
		"FROM (SELECT * FROM additional_timetable WHERE",
		strings.Join(conditions, " AND "),
		") AS tb",
		"LEFT JOIN classroom AS c ON tb.class_id = c.id",
		"LEFT JOIN duration AS d ON tb.duration_id = d.id",
		"LEFT JOIN subject AS s ON tb.subject_id = s.id",
		"LEFT JOIN teacher AS t ON tb.teacher_id = t.id",
		"LEFT JOIN teacher AS t2 ON tb.add_teacher_id = t2.id",
	}, " "))
	if err != nil {
		return res, errors.ErrorWrap(err)
	}
	var additional_timetable []usecase.Timetable
	var additional_timetable_temp []usecase.Timetable
	for rows.Next() {
		var t usecase.Timetable
		var id int
		var teach_id int
		var teach_name string
		var teach2_id sql.NullInt64
		var teach2_name sql.NullString
		err := rows.Scan(&id, &t.ClassId, &t.ClassName, &t.DurationId,
			&t.DurationName, &t.FrameId,
			&t.SubjectId, &t.SubjectName, &teach_id, &teach_name,
			&t.Day, &teach2_id, &teach2_name,
		)
		t.FrameDayWeek = t.FrameId/usecase.PERIOD + 1
		t.FramePeriod = t.FrameId % usecase.PERIOD
		t.TeacherIds = []int{teach_id}
		t.TeacherNames = []string{teach_name}
		if teach2_id.Valid {
			if !teach2_name.Valid {
				return nil, errors.NewError("addtional teacher name is null")
			}
			t.TeacherIds = append(t.TeacherIds, int(teach2_id.Int64))
			t.TeacherNames = append(t.TeacherNames, teach2_name.String)
		}
		t.Id = id*1e6 + t.Day.Year()%100*1e4 + (int(t.Day.Month())+20)*100 + t.Day.Day()
		if err != nil {
			return res, errors.ErrorWrap(err)
		}
		additional_timetable_temp = append(additional_timetable_temp, t)
	}

	// 重複削除
	// duration, frame, day 同じかつ (クラス同じか教師同じ) のとき重複
	// day と frame の曜日が異なるものは削除

	type key1 struct {
		c, d, f, day int
	}
	type key2 struct {
		d, f, day int
		t         [2]int
	}
	day_hash := func(t time.Time) int {
		return t.Year()*400 + t.YearDay()
	}
	m1 := make(map[key1]struct{})
	m2 := make(map[key2]struct{})
	for i := len(additional_timetable_temp) - 1; i >= 0; i-- {
		t := additional_timetable_temp[i]
		dayweek := t.FrameId/usecase.PERIOD + 1
		if int(t.Day.Weekday()) != dayweek {
			continue
		}
		k1 := key1{t.ClassId, t.DurationId, t.FrameId, day_hash(t.Day)}
		k2 := key2{t.DurationId, t.FrameId, day_hash(t.Day), [2]int{t.TeacherIds[0], -1}}
		if len(t.TeacherIds) > 1 {
			k2.t[1] = t.TeacherIds[1]
		}
		if k2.t[0] < k2.t[1] {
			k2.t[0], k2.t[1] = k2.t[1], k2.t[0]
		}
		_, exist1 := m1[k1]
		_, exist2 := m2[k2]
		if exist1 || exist2 {
			continue
		}
		additional_timetable = append(additional_timetable, t)
		m1[k1] = struct{}{}
		m2[k2] = struct{}{}
	}

	// コマで sort

	sort.Slice(additional_timetable, func(i int, j int) bool {
		t1 := additional_timetable[i]
		t2 := additional_timetable[j]
		return t1.Day.Before(t2.Day) || (t1.Day.Equal(t2.Day) && t1.FrameId < t2.FrameId)
	})

	// 通常時間割の一部削除を取得

	type delNormalTableKey struct {
		id  int
		day time.Time
	}
	del_normal_set := make(map[delNormalTableKey]struct{}) // index: day, value: id list
	rows, err = db.Query("SELECT id,day FROM deleted_normal_timetable")
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	for rows.Next() {
		var id int
		var d time.Time
		err := rows.Scan(&id, &d)
		if err != nil {
			return nil, errors.ErrorWrap(err)
		}
		del_normal_set[delNormalTableKey{
			id:  id,
			day: d,
		}] = struct{}{}
	}

	// 祝日取得

	holidays, err := getHolidays()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}

	// 時間割を合わせる

	x := 0
	hi := 0
	for d := start_day; !d.After(end_day); d = d.AddDate(0, 0, 1) {
		y := 0
		day_week := int(d.Weekday())
		if hi < len(holidays) && holidays[hi].Equal(d) {
			hi++
			continue
		}
		add_normal := func(id int) {
			for ; y < len(normal[day_week]) && normal[day_week][y].FrameId < id; y++ {
				_, del := del_normal_set[delNormalTableKey{
					id:  normal[day_week][y].Id,
					day: d,
				}]
				if del {
					continue
				}
				t := &usecase.Timetable{
					NormalTimetable: normal[day_week][y],
					Day:             d,
				}
				t.Id = t.Id*1e6 + t.Day.Year()%100*1e4 + int(t.Day.Month())*100 + t.Day.Day()
				res = append(res, *t)
			}
			if y < len(normal[day_week]) && normal[day_week][y].FrameId == id {
				y++
			}
		}
		for ; x < len(additional_timetable); x++ {
			if d.Format("2006-01-02") != additional_timetable[x].Day.Format("2006-01-02") {
				break
			}
			add_normal(additional_timetable[x].FrameId)
			res = append(res, additional_timetable[x])
		}
		add_normal(math.MaxInt)
	}
	return res, nil
}
