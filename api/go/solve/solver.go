package solve

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"time"

	"timetable/errors"
	"timetable/library/bitset"
	"timetable/usecase"
)

const (
	INF  int = 1e9 + 7
	D, P int = usecase.COUNT_DAY, usecase.PERIOD
)

var (
	ERR_CANT_SOLVE error = errors.NewError("internal error, dp all INF, can't solve")
)

type Heap [][4]int

func (h Heap) Len() int {
	return len(h)
}
func (h Heap) Less(i, j int) bool {
	return h[i][0] < h[j][0]
}
func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *Heap) Push(x interface{}) {
	y := x.([4]int)
	*h = append(*h, y)
}
func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func isListCross(a []int, b []int) bool {
	for _, i := range a {
		for _, j := range b {
			if i == j {
				return true
			}
		}
	}
	return false
}

// index が timetable に対応した teacher の avoid を返す
func NewTimetableAvoid(teacher []usecase.Teacher, table []usecase.Timetable, start_day time.Time) ([][][]int, error) {
	// avoid[timetable index][day][period]
	id2index := make(map[int]int, len(teacher))
	for i, t := range teacher {
		id2index[t.Id] = i
	}
	avoid := make([][][]int, len(table))
	for i := 0; i < len(table); i++ {
		avoid[i] = make([][]int, D)
		for j := 0; j < D; j++ {
			avoid[i][j] = make([]int, P)
		}
		for _, tid := range table[i].TeacherIds {
			idx, ok := id2index[tid]
			if !ok {
				return nil, errors.NewError("internal error, table index not found")
			}
			for j := 0; j < D; j++ {
				d := (start_day.Day() + j) % 7
				if d == 0 {
					continue
				}
				for k := 0; k < P; k++ {
					wi := (d-1)*P + k
					if wi < len(teacher[idx].Avoid) {
						avoid[i][j][k] += teacher[idx].Avoid[wi]
					}
				}
			}
		}
	}
	return avoid, nil
}

// NewPlaceIndexes: index が timetable に対応した place の index の配列を返す
func NewPlaceIndexes(place []usecase.Place, table []usecase.Timetable) ([]int, error) {
	res := make([]int, len(table))
	// m1: place id -> index
	m1 := make(map[int]int)
	for i, p := range place {
		m1[p.Id] = i
	}
	for i := 0; i < len(table); i++ {
		pi, ok := m1[table[i].PlaceId]
		if !ok {
			return nil, errors.NewError("place id is not found")
		}
		res[i] = pi
	}
	return res, nil
}

// timetableChangeSolver
func timetableChangeSolver(cost [][][][]int, start [2]int, units *[D][P][]int) ([][2]int, error) {
	// response: s[0]->s[1], s[1]->s[2], ... s[-1]->s[0]
	const cost_move int = 100   // 1 コマの移動に対するコスト、まとめた数かける
	const cost_move_day int = 1 // 移動の日数に対するコスト
	const K int = 10            // 入れ替え回数の上限+1
	const IDX int = 3           // index サイズ

	abs := func(x int) int {
		if x >= 0 {
			return x
		} else {
			return -x
		}
	}

	n := len(cost)
	if n <= 0 {
		return nil, errors.NewError("input error, length")
	}
	m := len(cost[0])
	if m <= 0 {
		return nil, errors.NewError("input error, length")
	}
	// dp[i][j][k]: i,j 枠、k 回入れ替えの min cost
	dp := make([][][]int, n)
	bef := make([][][][IDX]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([][]int, m)
		bef[i] = make([][][IDX]int, m)
		for j := 0; j < m; j++ {
			dp[i][j] = make([]int, K)
			bef[i][j] = make([][IDX]int, K)
			for k := 0; k < K; k++ {
				dp[i][j][k] = INF
			}
		}
	}
	sx := start[0]
	sy := start[1]
	if sx < 0 || n <= sx || sy < 0 || m <= sy {
		return nil, errors.NewError("input error, start")
	}
	dp[sx][sy][0] = 0
	h := &Heap{
		[IDX + 1]int{0, sx, sy, 0},
	}
	log.Println("start index", sx, sy)
	heap.Init(h)
	s := INF
	var ex, ey, ez int
	for h.Len() > 0 {
		mc := heap.Pop(h).([IDX + 1]int)
		c, x, y, z := mc[0], mc[1], mc[2], mc[3]
		log.Println("in dp", x, y, y, z)
		if s <= c {
			break
		}
		if dp[x][y][z] < c {
			continue
		}
		if c < s && sx == x && sy == y && z > 0 {
			s = c
			ex, ey, ez = x, y, z
		}
		if sx == x && sy == y && z > 0 {
			continue
		}
		if K <= z+1 {
			continue
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if i == x && j == y {
					continue
				}
				tc := c + cost[x][y][i][j] + cost_move*len(units[i][j]) + cost_move_day*abs(x-i)
				if dp[i][j][z+1] <= tc {
					continue
				}
				dp[i][j][z+1] = tc
				bef[i][j][z+1] = [IDX]int{x, y, z}
				heap.Push(h, [IDX + 1]int{tc, i, j, z + 1})
			}
		}
	}
	log.Println("end index, score", ex, ey, ez, s)
	if ez == 0 {
		return nil, ERR_CANT_SOLVE
	}
	//log.Println("dp", dp)
	//log.Println("bef", bef)
	var res [][2]int
	for i := 0; i < K && ez > 0; i++ {
		b := bef[ex][ey][ez]
		x, y, z := b[0], b[1], b[2]
		res = append(res, [2]int{x, y})
		ex, ey, ez = x, y, z
		log.Println(x, y, z)
	}
	if ex != sx || ey != sy || ez != 0 {
		return nil, errors.NewError("internal error")
	}
	rn := len(res)
	for i := 0; i < rn/2; i++ {
		res[i], res[rn-i-1] = res[rn-i-1], res[i]
	}
	return res, nil
}

// 時間割がまとめられるかどうか
// return able: []bitset able[i][j*P+k]: クラス i (index) は j,k コマでまとめられるかどうか
func ableCompressTimetable(
	tt_all []usecase.Timetable,
	graph usecase.ClassGraph,
	start_day time.Time,
) ([]bitset.Bitset, error) {
	n := len(graph.Nodes)
	classes := make([][]int, D*P)
	for _, t := range tt_all {
		cls, ok := graph.Id2index[t.ClassId]
		if !ok {
			return nil, errors.NewError("class id error", t.ClassId)
		}
		d := int(t.Day.Sub(start_day).Hours())
		if d < 0 || D <= d/24 {
			continue
		}
		d /= 24
		idx := d*P + t.FramePeriod
		classes[idx] = append(classes[idx], cls)
	}
	ful := make([][]int, n)
	// これが最適化されないため遅いかもしれない
	for i := 0; i < n; i++ {
		for j, f := range graph.Nodes[i].Parent {
			pf := graph.Nodes[i].ParentFill[j]
			if pf.Board == -1 {
				continue
			}
			nes := pf.Board - len(ful[f]) + 1
			if nes > 0 {
				ful[f] = append(ful[f], make([]int, nes)...)
			}
			ful[f][pf.Board] |= 1 << pf.Piece
		}
	}
	res := make([]bitset.Bitset, n)
	for i := 0; i < n; i++ {
		res[i] = bitset.NewBitset(D * P)
	}
	for i := 0; i < D*P; i++ {
		// class index, board, piece
		// board -2 のときは full とする、-1 は捨て
		var vec [][3]int
		for j := 0; j < len(classes[i]); j++ {
			vec = append(vec, [3]int{
				classes[i][j],
				-2, -2,
			})
		}
		board := make([]int, n)
		want := make([]int, n)
		for j := 0; j < n; j++ {
			board[j] = -1
		}
		for len(vec) > 0 {
			las := vec[len(vec)-1]
			vec = vec[:len(vec)-1]
			c, b, p := las[0], las[1], las[2]
			if b == -1 || board[c] == -2 {
				continue
			}
			if b != -2 {
				if board[c] == -1 {
					board[c] = b
					want[c] = ful[c][b]
				}
				if board[c] != b {
					board[c] = -2
					continue
				}
				want[c] &= ^(1 << p)
				if want[c] > 0 {
					continue
				}
			}
			board[c] = -2
			res[c].Set(i, true)
			for j, par := range graph.Nodes[c].Parent {
				pf := graph.Nodes[c].ParentFill[j]
				vec = append(vec, [3]int{
					par,
					pf.Board, pf.Piece,
				})
			}
		}

	}
	return res, nil
}

// 時間割圧縮
// return (units, others, idxs, error)
// units: 圧縮された時間割、tt_all の index が入っている
// others: units 以外のもの
// idxs: units に入っているコマの index
func compressTimetable(
	tt_all []usecase.Timetable,
	graph usecase.ClassGraph,
	class_idx int,
	start_day time.Time,
	class_avail []bool,
	holidays []time.Time,
	can_compress bitset.Bitset,
) (*[D][P][]int, *[D][P][]int, [][2]int, error) {
	n := len(graph.Nodes)
	in_class := make([]bool, n)
	for i := 0; i < n; i++ {
		if graph.NodeIn(class_idx, i) {
			in_class[i] = true
			continue
		}
	}
	var idxs [][2]int
	var units [D][P][]int
	var other [D][P][]int
	var flag [D][P]bool // まとめられないコマが true
	for i := 0; i < D; i++ {
		for j := 0; j < P; j++ {
			flag[i][j] = !can_compress.Test(i*P + j)
		}
	}
	for i, t := range tt_all {
		d := int(t.Day.Sub(start_day).Hours())
		if d < 0 || D <= d/24 {
			continue
		}
		d /= 24
		p := t.FramePeriod
		c := graph.Id2index[t.ClassId]
		if !flag[d][p] && in_class[c] {
			units[d][p] = append(units[d][p], i)
		} else {
			other[d][p] = append(other[d][p], i)
		}

	}
	// 祝日などのときも flag true にする
	for _, t := range holidays {
		d := int(t.Sub(start_day).Hours())
		if d < 0 {
			continue
		}
		d /= 24
		if D <= d {
			break
		}
		for i := 0; i < P; i++ {
			flag[d][i] = true
		}
	}
	for i := 0; i < D; i++ {
		day := start_day.AddDate(0, 0, i)
		for j := 0; j < P; j++ {
			if day.Weekday() == 0 {
				flag[i][j] = true
				continue
			}
			k := (int(day.Weekday())-1)*P + j
			if len(class_avail) <= k || !class_avail[k] {
				flag[i][j] = true
			}
		}
	}
	for i := 0; i < D; i++ {
		for j := 0; j < P; j++ {
			if flag[i][j] {
				other[i][j] = append(other[i][j], units[i][j]...)
				units[i][j] = []int{}
			} else {
				idxs = append(idxs, [2]int{i, j})
			}
		}
	}
	return &units, &other, idxs, nil
}

type calcCost struct {
}

func (*calcCost) initCost() [][][][]int {
	// init cost
	cost := make([][][][]int, D)
	for i := 0; i < D; i++ {
		cost[i] = make([][][]int, P)
		for j := 0; j < P; j++ {
			cost[i][j] = make([][]int, D)
			for k := 0; k < D; k++ {
				cost[i][j][k] = make([]int, P)
				for l := 0; l < P; l++ {
					cost[i][j][k][l] = INF
				}
			}
		}
	}
	return cost
}

func (*calcCost) getPlaceCount(
	other_units *[D][P][]int,
	place_indexes []int,
	places []usecase.Place,
) *[D][P][]int {
	var place_count [D][P][]int // other を埋めたときの count の残り
	for i := 0; i < D; i++ {
		for j := 0; j < P; j++ {
			place_count[i][j] = make([]int, len(places))
			for k, p := range places {
				place_count[i][j][k] = p.Count
			}
			for _, u := range other_units[i][j] {
				pi := place_indexes[u]
				place_count[i][j][pi]--
			}
		}
	}
	return &place_count
}

// teacher について、コマが重複するかどうか
func (*calcCost) getTeacherInval(
	other_units *[D][P][]int,
	tt_all []usecase.Timetable,
	start_day time.Time,
) ([][]int, []bitset.Bitset) {
	tids := make([][]int, len(tt_all)) // timetable index に対する teacher index の配列
	teach_id2idx := make(map[int]int)
	var tea_inval []bitset.Bitset // teacher is invalid or not
	for i, t := range tt_all {
		for _, id := range t.TeacherIds {
			idx, ok := teach_id2idx[id]
			if !ok {
				idx = len(teach_id2idx)
				teach_id2idx[id] = idx
				tea_inval = append(tea_inval, bitset.NewBitset(D*P))
			}
			tids[i] = append(tids[i], idx)
		}
	}
	for i := 0; i < D; i++ {
		for j := 0; j < P; j++ {
			for _, u := range other_units[i][j] {
				t := tt_all[u]
				d := int((t.Day.Sub(start_day)).Hours())
				if d < 0 || D <= d/24 {
					continue
				}
				d /= 24
				p := t.FramePeriod
				for _, tid := range tids[u] {
					tea_inval[tid].Set(d*P+p, true)
				}
			}
		}
	}
	return tids, tea_inval
}

// 同じ先生のコマが戻らないように
func (*calcCost) BanReturn(
	cost [][][][]int,
	start_idx [2]int,
	units *[D][P][]int,
	tt_all []usecase.Timetable,
	start_day time.Time,
	change_unit *usecase.Timetable,
) {
	si, sj := start_idx[0], start_idx[1]
	for i := 0; i < D; i++ {
		for j := 0; j < P; j++ {
			for _, u := range units[i][j] {
				t := tt_all[u]
				if !isListCross(t.TeacherIds, change_unit.TeacherIds) {
					continue
				}
				d := int((t.Day.Sub(start_day)).Hours())
				if d < 0 || D <= d/24 {
					continue
				}
				d /= 24
				p := t.FramePeriod
				cost[d][p][si][sj] = INF
			}
		}
	}
}

func (*calcCost) getAvoidCost(avoid int) int {
	if avoid > 5 {
		return INF
	}
	return avoid * 10
}

//idxs: 使えるコマの index
func getCost(
	units *[D][P][]int,
	other_units *[D][P][]int,
	idxs [][2]int,
	avoid [][][]int,
	place_indexes []int,
	places []usecase.Place,
	start_idx [2]int,
	tt_all []usecase.Timetable,
	start_day time.Time,
	change_unit *usecase.Timetable,
) [][][][]int {
	var cal calcCost
	cost := cal.initCost()
	place_count := cal.getPlaceCount(other_units, place_indexes, places)
	tids, tea_inval := cal.getTeacherInval(other_units, tt_all, start_day)
	//define cost
	for _, ij := range idxs {
		i, j := ij[0], ij[1]
		for _, kl := range idxs {
			k, l := kl[0], kl[1]
			cost[i][j][k][l] = 0
			for _, u := range units[i][j] {
				cost[i][j][k][l] += cal.getAvoidCost(avoid[u][k][l])
				pi := place_indexes[u]
				place_count[k][l][pi]--
				if place_count[k][l][pi] < 0 {
					cost[i][j][k][l] = INF
				}

				//先生がそのコマに入れるかどうか
				for _, tea_id := range tids[u] { // teacher id
					if tea_inval[tea_id].Test(k*P + l) {
						cost[i][j][k][l] = INF
					}
				}
			}
			for _, u := range units[i][j] {
				pi := place_indexes[u]
				place_count[k][l][pi]++
			}
		}
	}

	cal.BanReturn(cost, start_idx, units, tt_all, start_day, change_unit)
	return cost
}

func getFinalCost(
	mv [][2]int,
	first_cost [][][][]int, //getCost で得られるもの

) int {
	// mv: indexes
	var cost int = 0
	for i, v := range mv {
		to := mv[(i+1)%len(mv)]
		cost += first_cost[v[0]][v[1]][to[0]][to[1]]
	}
	cost += len(mv) * 1000
	return cost
}

// コマに対して同じ先生が重複してもよいとする。コストには反映
// cost_teach_inval: 先生が重複することに対するコスト
func getCostRelaxTeacher(
	units *[D][P][]int,
	other_units *[D][P][]int,
	idxs [][2]int,
	avoid [][][]int,
	place_indexes []int,
	places []usecase.Place,
	start_idx [2]int,
	tt_all []usecase.Timetable,
	start_day time.Time,
	cost_teach_inval int,
	change_unit *usecase.Timetable,
) [][][][]int {
	var cal calcCost
	cost := cal.initCost()
	place_count := cal.getPlaceCount(other_units, place_indexes, places)
	tids, tea_inval := cal.getTeacherInval(other_units, tt_all, start_day)
	//define cost
	for _, ij := range idxs {
		i, j := ij[0], ij[1]
		for _, kl := range idxs {
			k, l := kl[0], kl[1]
			cost[i][j][k][l] = 0
			for _, u := range units[i][j] {
				cost[i][j][k][l] += cal.getAvoidCost(avoid[u][k][l])
				pi := place_indexes[u]
				place_count[k][l][pi]--
				if place_count[k][l][pi] < 0 {
					cost[i][j][k][l] = INF
				}

				//先生がそのコマに入れるかどうか
				for _, tea_id := range tids[u] {
					if tea_inval[tea_id].Test(k*P + l) {
						cost[i][j][k][l] += cost_teach_inval
					}
				}
			}
			for _, u := range units[i][j] {
				pi := place_indexes[u]
				place_count[k][l][pi]++
			}
		}
	}

	cal.BanReturn(cost, start_idx, units, tt_all, start_day, change_unit)
	return cost
}

// 	others で、動かしたとき教師が被っているを列挙する
func getTeacherInval(units *[D][P][]int, others *[D][P][]int, tt_all []usecase.Timetable, move [][2]int) []usecase.Timetable {
	var res []usecase.Timetable
	for i := 0; i < len(move); i++ {
		vi, vj := move[i][0], move[i][1]
		ni := (i + 1) % len(move)
		ti, tj := move[ni][0], move[ni][1]
		for _, u := range units[vi][vj] {
			for _, bs := range others[ti][tj] {
				if isListCross(tt_all[u].TeacherIds, tt_all[bs].TeacherIds) {
					res = append(res, tt_all[bs])
				}
			}
		}
	}
	return res
}

func ApplyChange(tt_all []usecase.Timetable, move []usecase.TimetableMove) {
	m1 := make(map[int]int)
	for i, mv := range move {
		m1[mv.Unit.Id] = i
	}
	for i, t := range tt_all {
		idx, ok := m1[t.Id]
		if !ok {
			continue
		}
		mv := move[idx]
		tt_all[i].Day = mv.Day
		tt_all[i].FrameId = mv.FrameId
		tt_all[i].FrameDayWeek = mv.FrameId / P
		tt_all[i].FramePeriod = mv.FrameId % P
	}
}
func CancelChange(tt_all []usecase.Timetable, move []usecase.TimetableMove) {
	m1 := make(map[int]int)
	for i, mv := range move {
		m1[mv.Unit.Id] = i
	}
	for i, t := range tt_all {
		idx, ok := m1[t.Id]
		if !ok {
			continue
		}
		mv := move[idx].Unit
		tt_all[i].Day = mv.Day
		tt_all[i].FrameId = mv.FrameId
		tt_all[i].FrameDayWeek = mv.FrameDayWeek
		tt_all[i].FramePeriod = mv.FramePeriod
	}
}
func AppendMove(mv []usecase.TimetableMove, plus []usecase.TimetableMove) error {
	for _, v := range plus {
		flag := false
		for _, t := range mv {
			if v.Unit.Id == t.Unit.Id {
				if t.Day != v.Unit.Day || t.FrameId != v.Unit.FrameId {
					return errors.NewError(fmt.Sprintf("append error, from: %v,to: %v", t, v))
				}
				t.Day = v.Day
				t.FrameId = v.FrameId
				flag = true
				break
			}
			if !flag {
				mv = append(mv, v)
			}
		}
	}
	return nil
}

type SolverClass struct {
}

// 時間割変更
// args:
// teacher_relax: -1: 緩和しない。0 以上で cost に対応
// return move, cost, error
func (*SolverClass) TimetableChange(
	tt_all []usecase.Timetable,
	graph usecase.ClassGraph,
	change_unit *usecase.Timetable,
	places []usecase.Place,
	teachers []usecase.Teacher,
	start_day time.Time,
	holidays []time.Time,
	teacher_relax int,
) ([]usecase.TimetableMove, int, error) {
	cost_inf := math.MaxInt

	class_queue := [][2]int{
		{
			graph.Id2index[change_unit.ClassId], 0,
		},
	}
	class_used := make([]bool, len(graph.Nodes))
	var final_move []usecase.TimetableMove
	var final_cost int = INF
	avoid, err := NewTimetableAvoid(teachers, tt_all, start_day)
	if err != nil {
		return nil, 0, errors.ErrorWrap(err)
	}
	place_indexes, err := NewPlaceIndexes(places, tt_all)
	if err != nil {
		return nil, 0, errors.ErrorWrap(err)
	}
	can_compress, err := ableCompressTimetable(tt_all, graph, start_day)
	if err != nil {
		return nil, 0, errors.ErrorWrap(err)
	}
	createMove := func(units *[D][P][]int, mv [][2]int) []usecase.TimetableMove {
		var res []usecase.TimetableMove
		for i, v := range mv {
			to := mv[(i+1)%len(mv)]
			vi, vj := v[0], v[1]
			ti, tj := to[0], to[1]
			day := start_day.AddDate(0, 0, ti)
			for _, idx := range units[vi][vj] {
				t := tt_all[idx]
				res = append(res, usecase.TimetableMove{
					Unit:    t,
					Day:     day,
					FrameId: (int(day.Weekday())-1)*P + tj,
				})
			}
		}
		return res
	}

	change_unit_idx := [2]int{
		int(change_unit.Day.Sub(start_day).Hours()) / 24, change_unit.FramePeriod,
	}
	for len(class_queue) > 0 {
		class_idx := class_queue[0][0]
		class_dis := class_queue[0][1]
		class_queue = class_queue[1:]
		for _, p := range graph.Nodes[class_idx].Parent {
			if class_used[p] {
				continue
			}
			class_used[p] = true
			class_queue = append(class_queue, [2]int{
				p, class_dis + 1,
			})
		}

		units, others, idxs, err := compressTimetable(
			tt_all, graph, class_idx, start_day, graph.Nodes[class_idx].Available, holidays, can_compress[class_idx],
		)
		if err != nil {
			return nil, cost_inf, errors.ErrorWrap(err)
		}
		// log.Println("class id", graph.Nodes[class_idx].Id)
		// log.Println("compress able", can_compress[class_idx])
		// log.Println("idxs", idxs)
		// for _, l := range idxs {
		// 	i, j := l[0], l[1]
		// 	log.Println("units", i, j, units[i][j])
		// }

		var cost [][][][]int
		if teacher_relax == -1 {
			cost = getCost(units, others, idxs, avoid, place_indexes, places, change_unit_idx, tt_all, start_day, change_unit)
		} else {
			cost = getCostRelaxTeacher(units, others, idxs, avoid, place_indexes, places, change_unit_idx, tt_all, start_day, teacher_relax, change_unit)
		}

		mv, err := timetableChangeSolver(cost, change_unit_idx, units)
		if err == ERR_CANT_SOLVE {
			continue
		}
		if err != nil {
			return nil, cost_inf, errors.ErrorWrap(err)
		}
		if teacher_relax == -1 {
			if c := getFinalCost(mv, cost); c < final_cost {
				final_cost = c
				final_move = createMove(units, mv)
			}
			continue
		}
		bef_unit := *change_unit
		co := getFinalCost(mv, cost)
		next_units := getTeacherInval(units, others, tt_all, mv)
		move := createMove(units, mv)
		var sc SolverClass
		log.Print("movebef", co, move)
		ApplyChange(tt_all, move)
		for _, u := range next_units {
			flag := false // unit を動かす必要があるのかどうか
			for _, t := range tt_all {
				if u.Id == t.Id {
					flag = (u.Day.Equal(t.Day) && u.FrameId == t.FrameId)
					break
				}
			}
			if !flag {
				continue
			}
			tm, c, err := sc.TimetableChange(tt_all, graph, &u, places, teachers, start_day, holidays, -1)
			if err != nil || c >= INF {
				co = INF
				break
			}
			co += c
			err = AppendMove(move, tm)
			if err != nil {
				return nil, 0, errors.ErrorWrap(err)
			}
			ApplyChange(tt_all, tm)
		}
		for _, t := range tt_all {
			if isListCross(t.TeacherIds, bef_unit.TeacherIds) && t.Day.Equal(bef_unit.Day) && t.FrameId == bef_unit.FrameId {
				// 元々のコマに動かす必要のある先生のコマ存在
				tm, c, err := sc.TimetableChange(tt_all, graph, &t, places, teachers, start_day, holidays, -1)
				if err != nil || c >= INF {
					co = INF
				}
				co += c
				err = AppendMove(move, tm)
				if err != nil {
					return nil, 0, errors.ErrorWrap(err)
				}
				ApplyChange(tt_all, tm)
				break
			}
		}
		log.Print("move", co, move)
		if co < final_cost {
			final_cost = co
			final_move = move
		}
		CancelChange(tt_all, move)
	}
	if final_cost == INF {
		return nil, INF, errors.ErrorWrap(ERR_CANT_SOLVE)
	}
	return final_move, final_cost, nil
}
