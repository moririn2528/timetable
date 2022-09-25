package database

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
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

func parseNormal2id(id int, date time.Time) int {
	return id*1e6 + date.Year()%100*1e4 + int(date.Month())*100 + date.Day()
}
func parseAdditional2id(id int, date time.Time) int {
	return id*1e6 + date.Year()%100*1e4 + (int(date.Month())+20)*100 + date.Day()
}
func parse2DatabaseId(id int) int {
	return id / 1e6
}
func isNormalId(id int) bool {
	id /= 100
	return id%100 < 20
}

// 時刻をローカル時間で 0:00 に設定する
func correntDayTime(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0, usecase.JST,
	)
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
	for rows.Next() {
		var t usecase.NormalTimetable
		var teach_id int
		var teach_name string
		var teach2_id sql.NullInt64
		var teach2_name sql.NullString
		err := rows.Scan(&t.Id, &t.ClassId, &t.ClassName, &t.DurationId,
			&t.DurationName, &t.FrameId,
			&t.SubjectId, &t.SubjectName, &teach_id, &teach_name, &teach2_id, &teach2_name)
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
		res = append(res, t)
	}

	// 重複取り除く
	// duration, frame 同じかつ (クラス同じか教師同じ) のとき重複
	// TODO: これは整合性チェックで行うべき、sql error として返すようにする
	// type key1 struct {
	// 	c, d, f int
	// }
	// type key2 struct {
	// 	d, f int
	// 	t    [2]int
	// }
	// m1 := make(map[key1]struct{})
	// m2 := make(map[key2]struct{})
	// for i := len(res_temp) - 1; i >= 0; i-- {
	// 	t := res_temp[i]
	// 	k1 := key1{t.ClassId, t.DurationId, t.FrameId}
	// 	k2 := key2{t.DurationId, t.FrameId, [2]int{t.TeacherIds[0], -1}}
	// 	if len(t.TeacherIds) > 1 {
	// 		k2.t[1] = t.TeacherIds[1]
	// 	}
	// 	if k2.t[0] < k2.t[1] {
	// 		k2.t[0], k2.t[1] = k2.t[1], k2.t[0]
	// 	}
	// 	_, exist1 := m1[k1]
	// 	_, exist2 := m2[k2]
	// 	if exist1 || exist2 {
	// 		continue
	// 	}
	// 	res = append(res, t)
	// 	m1[k1] = struct{}{}
	// 	m2[k2] = struct{}{}
	// }
	return res, nil
}

// GetDeletedNormalTimetable:
// sql, deleted_normal_timetable に入っている内容を返す
// これは normal_timetable の時間割の一部が変更された時に使う database
// class_ids は空を許す、その時は条件として入れない
// teacher_id は -1 を許す、その時は条件として入れない
func (*DatabaseTimetable) GetDeletedNormalTimetable(duration_id int, class_ids []int, teacher_id int) ([]usecase.DeletedNormalTimetable, error) {
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

	var res []usecase.DeletedNormalTimetable
	err := db.Select(&res, strings.Join([]string{
		"SELECT d.id, d.normal_id, d.day FROM deleted_normal_timetable AS d",
		"LEFT JOIN normal_timetable AS t ON d.normal_id=t.id",
		"WHERE",
		strings.Join(conditions, " AND "),
	}, " "))
	if err != nil {
		return []usecase.DeletedNormalTimetable{}, errors.ErrorWrap(err)
	}
	for i, u := range res {
		res[i].Day = correntDayTime(u.Day)
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
	start_day = correntDayTime(start_day)
	var res []usecase.Timetable

	// 曜日で分割
	var normal [7][]usecase.NormalTimetable // id sorted

	{ // 通常の時間割取得
		tt, err := dc.GetNomalTimetable(duration_id, class_ids, teacher_id)
		if err != nil {
			return res, errors.ErrorWrap(err)
		}
		sort.Slice(tt, func(i int, j int) bool {
			return tt[i].FrameId < tt[j].FrameId
		})
		j := 0
		for _, v := range tt {
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
	}

	var additional []usecase.Timetable
	// 追加時間割を取得
	// sql 条件生成
	{
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
		defer rows.Close()
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
			if err != nil {
				return res, errors.ErrorWrap(err)
			}
			t.TeacherIds = []int{teach_id}
			t.TeacherNames = []string{teach_name}
			if teach2_id.Valid {
				if !teach2_name.Valid {
					return nil, errors.NewError("addtional teacher name is null")
				}
				t.TeacherIds = append(t.TeacherIds, int(teach2_id.Int64))
				t.TeacherNames = append(t.TeacherNames, teach2_name.String)
			}
			t.Id = parseAdditional2id(t.Id, t.Day)
			t.Day = correntDayTime(t.Day)
			additional = append(additional, t)
		}

		// 重複削除
		// duration, frame, day 同じかつ (クラス同じか教師同じ) のとき重複
		// day と frame の曜日が異なるものは削除
		// TODO: 整合性チェックに入れる

		// type key1 struct {
		// 	c, d, f, day int
		// }
		// type key2 struct {
		// 	d, f, day int
		// 	t         [2]int
		// }
		// day_hash := func(t time.Time) int {
		// 	return t.Year()*400 + t.YearDay()
		// }
		// m1 := make(map[key1]struct{})
		// m2 := make(map[key2]struct{})
		// for i := len(additional_timetable_temp) - 1; i >= 0; i-- {
		// 	t := additional_timetable_temp[i]
		// 	dayweek := t.FrameId/usecase.PERIOD + 1
		// 	if int(t.Day.Weekday()) != dayweek {
		// 		continue
		// 	}
		// 	k1 := key1{t.ClassId, t.DurationId, t.FrameId, day_hash(t.Day)}
		// 	k2 := key2{t.DurationId, t.FrameId, day_hash(t.Day), [2]int{t.TeacherIds[0], -1}}
		// 	if len(t.TeacherIds) > 1 {
		// 		k2.t[1] = t.TeacherIds[1]
		// 	}
		// 	if k2.t[0] < k2.t[1] {
		// 		k2.t[0], k2.t[1] = k2.t[1], k2.t[0]
		// 	}
		// 	_, exist1 := m1[k1]
		// 	_, exist2 := m2[k2]
		// 	if exist1 || exist2 {
		// 		continue
		// 	}
		// 	additional_timetable = append(additional_timetable, t)
		// 	m1[k1] = struct{}{}
		// 	m2[k2] = struct{}{}
		// }
	}

	// コマで sort
	sort.Slice(additional, func(i int, j int) bool {
		t1 := additional[i]
		t2 := additional[j]
		return t1.Day.Before(t2.Day) || (t1.Day.Equal(t2.Day) && t1.FrameId < t2.FrameId)
	})

	// 通常時間割の一部削除を取得

	type delNormalTableKey struct {
		id  int
		day time.Time
	}
	del_normal := make(map[delNormalTableKey]struct{}) // index: day, value: id list
	{
		dels, err := dc.GetDeletedNormalTimetable(duration_id, class_ids, teacher_id)
		if err != nil {
			return nil, errors.ErrorWrap(err)
		}
		for _, d := range dels {
			del_normal[delNormalTableKey{
				id:  d.NormalId,
				day: d.Day,
			}] = struct{}{}
		}
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
				_, del := del_normal[delNormalTableKey{
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
				t.Id = parseNormal2id(t.Id, t.Day)
				res = append(res, *t)
			}
			if y < len(normal[day_week]) && normal[day_week][y].FrameId == id {
				y++
			}
		}
		for ; x < len(additional); x++ {
			if d.Format("2006-01-02") != additional[x].Day.Format("2006-01-02") {
				break
			}
			add_normal(additional[x].FrameId)
			res = append(res, additional[x])
		}
		add_normal(math.MaxInt)
	}
	return res, nil
}

func (dc *DatabaseTimetable) MoveTimetable(
	move []usecase.TimetableMove,
) error {
	// database と move の相違がないか調べる
	equal_teacher := func(t []int, t1 int, t2a sql.NullInt32) bool {
		if len(t) == 0 || len(t) > 2 {
			return false
		}
		if len(t) == 1 {
			return t[0] == t1 && !t2a.Valid
		}
		if !t2a.Valid {
			return false
		}
		t2 := int(t2a.Int32)
		if t1 == t2 {
			return false
		}
		return (t1 == t[0] && t2 == t[1]) || (t1 == t[1] && t2 == t[0])
	}
	type NormalTimetable struct {
		Id           int           `json:"id" db:"id"`
		DurationId   int           `json:"duration_id" db:"duration_id"`
		ClassId      int           `json:"class_id" db:"class_id"`
		TeacherId    int           `json:"teacher_id" db:"teacher_id"`
		SubjectId    int           `json:"subject_id" db:"subject_id"`
		FrameId      int           `json:"frame_id" db:"frame_id"`
		PlaceId      int           `json:"place_id" db:"place_id"`
		AddTeacherId sql.NullInt32 `json:"add_teacher_id" db:"add_teacher_id"`
	}
	type AdditionalTimetable struct {
		NormalTimetable
		Day time.Time `json:"day" db:"day"`
	}

	var nor, add []usecase.Timetable
	for _, v := range move {
		if isNormalId(v.Unit.Id) {
			nor = append(nor, v.Unit)
		} else {
			add = append(add, v.Unit)
		}
	}
	if len(nor) > 0 {
		var idstrs []string
		id2idx := make(map[int]int, len(nor))
		for i, t := range nor {
			idstrs = append(idstrs, strconv.Itoa(t.Id))
			id2idx[t.Id] = i
		}
		var rows []NormalTimetable
		err := db.Select(rows, strings.Join([]string{
			"SELECT id, duration_id, class_id, teacher_id, subject_id, frame_id,",
			"place_id, add_teacher_id FROM normal_timetable",
			"WHERE id in (" + strings.Join(idstrs, ",") + ")",
			"ORDERED BY id",
		}, " "))
		if err != nil {
			return errors.ErrorWrap(err)
		}
		cnt := 0
		for _, t1 := range rows {
			if err != nil {
				return errors.ErrorWrap(err)
			}
			idx, ok := id2idx[t1.Id]
			if !ok {
				return errors.NewError("assert error")
			}
			t2 := nor[idx]
			if !(t1.Id == parse2DatabaseId(t2.Id) && t1.DurationId == t2.DurationId && t1.ClassId == t2.ClassId && t1.SubjectId == t2.SubjectId &&
				t1.PlaceId == t2.PlaceId && equal_teacher(t2.TeacherIds, t1.TeacherId, t1.AddTeacherId) && t1.FrameId/usecase.PERIOD == int(t2.Day.Weekday())-1) {
				return errors.NewError(http.StatusBadRequest, "timetable is difference")
			}
			cnt++
		}
		if len(nor) != cnt {
			return errors.NewError(http.StatusBadRequest, "timetable is difference")
		}
	}
	if len(add) > 0 {
		var idstrs []string
		id2idx := make(map[int]int, len(nor))
		for i, t := range add {
			idstrs = append(idstrs, strconv.Itoa(t.Id))
			id2idx[t.Id] = i
		}
		var rows []AdditionalTimetable
		err := db.Select(rows, strings.Join([]string{
			"SELECT id, duration_id, class_id, teacher_id, subject_id, frame_id,",
			"place_id, add_teacher_id, day FROM additional_timetable",
			"WHERE id in (" + strings.Join(idstrs, ",") + ")",
			"ORDERED BY id",
		}, " "))
		if err != nil {
			return errors.ErrorWrap(err)
		}
		cnt := 0
		for _, t1 := range rows {
			if err != nil {
				return errors.ErrorWrap(err)
			}
			idx, ok := id2idx[t1.Id]
			if !ok {
				return errors.NewError("assert error")
			}
			t2 := add[idx]
			if !(t1.Id == parse2DatabaseId(t2.Id) && t1.DurationId == t2.DurationId && t1.ClassId == t2.ClassId && t1.SubjectId == t2.SubjectId &&
				t1.PlaceId == t2.PlaceId && t1.Day.Equal(t2.Day) && equal_teacher(t2.TeacherIds, t1.TeacherId, t1.AddTeacherId)) {
				return errors.NewError(http.StatusBadRequest, "timetable is difference")
			}
			cnt++
		}
		if len(add) != cnt {
			return errors.NewError(http.StatusBadRequest, "timetable is difference")
		}
	}

	// database から削除
	recoverNormal := func() {}
	if len(nor) > 0 {
		type DelTimetable struct {
			Id       int       `json:"id" db:"id"`
			NormalId int       `json:"normal_id" db:"normal_id"`
			Day      time.Time `json:"day" db:"day"`
		}
		var dels []DelTimetable
		for _, t := range nor {
			dels = append(dels, DelTimetable{
				Id:       t.Id,
				NormalId: parse2DatabaseId(t.Id),
				Day:      t.Day,
			})
		}
		res, err := db.Exec("INSERT INTO deleted_normal_timetable(id,normal_id,day) VALUES (:id,:normal_id,:day)", dels)
		if err != nil {
			return errors.ErrorWrap(err)
		}
		recoverNormal = func() {
			if len(nor) == 0 {
				return
			}
			var ids []string
			for _, t := range nor {
				ids = append(ids, strconv.Itoa(t.Id))
			}
			_, err = db.Exec("DELETE FROM deleted_normal_timetable WHERE id IN (" + strings.Join(ids, ",") + ")")
			if err != nil {
				log.Printf("Critical Error: %v", err)
			}
		}
		rnum, err := res.RowsAffected()
		if err != nil {
			recoverNormal()
			return errors.ErrorWrap(err)
		}
		if int(rnum) != len(nor) {
			recoverNormal()
			return errors.ErrorWrap(err)
		}
	}
	recoverAdditional := func() {}
	if len(add) > 0 {
		var ids []string
		for _, t := range add {
			ids = append(ids, strconv.Itoa(parse2DatabaseId(t.Id)))
		}
		res, err := db.Exec("DELETE FROM additional_timetable WHERE id in (" + strings.Join(ids, ",") + ")")
		if err != nil {
			recoverNormal()
			return errors.ErrorWrap(err)
		}
		var dels []AdditionalTimetable
		for _, t := range add {
			tim := AdditionalTimetable{
				NormalTimetable: NormalTimetable{
					Id:           parse2DatabaseId(t.Id),
					DurationId:   t.DurationId,
					ClassId:      t.ClassId,
					TeacherId:    t.TeacherIds[0],
					SubjectId:    t.SubjectId,
					FrameId:      t.FrameId,
					PlaceId:      t.PlaceId,
					AddTeacherId: sql.NullInt32{Int32: 0, Valid: false},
				},
				Day: t.Day,
			}
			if len(t.TeacherIds) == 2 {
				tim.AddTeacherId = sql.NullInt32{
					Int32: int32(t.TeacherIds[1]), Valid: true,
				}
			}
			dels = append(dels, tim)
		}
		recoverAdditional = func() {
			_, err = db.Exec("INSERT INTO additional_timetable"+
				"(id,duration_id,class_id,teacher_id,subject_id,frame_id,day,add_teacher_id) VALUES "+
				"(:id,:duration_id,:class_id,:teacher_id,:subject_id,:frame_id,:day,:add_teacher_id)",
				dels,
			)
			if err != nil {
				log.Printf("Critical Error: %v", err)
			}
		}
		recover := func() {
			recoverNormal()
			recoverAdditional()
		}
		rnum, err := res.RowsAffected()
		if err != nil {
			recover()
			return errors.ErrorWrap(err)
		}
		if int(rnum) != len(add) {
			recover()
			return errors.ErrorWrap(err)
		}
	}

	// database に追加
	recover := func() {
		recoverNormal()
		recoverAdditional()
	}
	{
		var inserts []AdditionalTimetable
		for _, mv := range move {
			t := mv.Unit
			tim := AdditionalTimetable{
				NormalTimetable: NormalTimetable{
					DurationId:   t.DurationId,
					ClassId:      t.ClassId,
					TeacherId:    t.TeacherIds[0],
					SubjectId:    t.SubjectId,
					FrameId:      mv.FrameId,
					PlaceId:      t.PlaceId,
					AddTeacherId: sql.NullInt32{Int32: 0, Valid: false},
				},
				Day: mv.Day,
			}
			if len(t.TeacherIds) == 2 {
				tim.AddTeacherId = sql.NullInt32{
					Int32: int32(t.TeacherIds[1]), Valid: true,
				}
			}
			inserts = append(inserts, tim)
		}
		_, err := db.Exec("INSERT INTO additional_timetable"+
			"(duration_id,class_id,teacher_id,subject_id,frame_id,day,add_teacher_id) VALUES "+
			"(:duration_id,:class_id,:teacher_id,:subject_id,:frame_id,:day,:add_teacher_id)",
			inserts,
		)
		if err != nil {
			recover()
			return errors.ErrorWrap(err)
		}
	}

	return nil
}
