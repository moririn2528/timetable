package solve

import (
	"container/heap"
	"log"
	"math"
	"time"

	"timetable/errors"
	"timetable/library/bitset"
	"timetable/usecase"
)

const (
	INF  int = 1e6 + 7
	D, P int = usecase.COUNT_DAY, usecase.PERIOD
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

// index が timetable に対応した teacher の avoid を返す
func NewTimetableAvoid(teacher []usecase.Teacher, table []usecase.Timetable, start_day time.Time) ([][][]int, error) {
	// avoid[timetable index][day][period]
	id2index := make(map[int]int, len(teacher))
	for i, t := range teacher {
		id2index[t.Id] = i
	}
	avoid := make([][][]int, len(table))
	for i := 0; i < len(table); i++ {
		tid := table[i].TeacherId
		idx, ok := id2index[tid]
		if !ok {
			return nil, errors.NewError("internal error, table index not found")
		}
		avoid[i] = make([][]int, D)
		for j := 0; j < D; j++ {
			avoid[i][j] = make([]int, P)
			d := (start_day.Day() + j) % 7
			if d == 0 {
				continue
			}
			for k := 0; k < P; k++ {
				wi := (d-1)*P + k
				if wi < len(teacher[idx].Avoid) {
					avoid[i][j][k] = teacher[idx].Avoid[wi]
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
	const cost_move int = 100 // 1 コマの移動に対するコスト、まとめた数かける
	const K int = 10          // 入れ替え回数の上限+1
	const IDX int = 3         // index サイズ
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
				tc := c + cost[x][y][i][j] + cost_move*len(units[i][j])
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
		return nil, errors.NewError("internal error, dp all INF, can't solve")
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

func getCost(
	units *[D][P][]int,
	other_units *[D][P][]int,
	idxs [][2]int,
	avoid [][][]int,
	place_indexes []int,
	places []usecase.Place,
) [][][][]int {
	//idxs: 使えるコマの index
	// 係数
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

	//define cost
	for _, ij := range idxs {
		i, j := ij[0], ij[1]
		for _, kl := range idxs {
			k, l := kl[0], kl[1]
			cost[i][j][k][l] = 0
			for _, u := range units[i][j] {
				cost[i][j][k][l] += avoid[u][k][l]
				pi := place_indexes[u]
				place_count[k][l][pi]--
				if place_count[k][l][pi] < 0 {
					cost[i][j][k][l] = INF
				}
			}
			for _, u := range units[i][j] {
				pi := place_indexes[u]
				place_count[k][l][pi]++
			}
		}
	}
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
	return cost
}

type SolverClass struct {
}

func (*SolverClass) TimetableChange(
	tt_all []usecase.Timetable,
	graph usecase.ClassGraph,
	change_unit *usecase.Timetable,
	places []usecase.Place,
	teachers []usecase.Teacher,
	start_day time.Time,
	holidays []time.Time,
) ([]usecase.Timetable, int, error) {
	// 時間割変更
	// move, cost, error
	cost_inf := math.MaxInt

	class_queue := [][2]int{
		{
			graph.Id2index[change_unit.ClassId], 0,
		},
	}
	class_used := make([]bool, len(graph.Nodes))
	var final_move [][2]int
	var final_cost int = INF
	var final_units *[D][P][]int
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

	change_unit_idx := [2]int{
		int(change_unit.Day.Sub(start_day).Hours()) / 24, change_unit.FramePeriod,
	}
	for len(class_queue) > 0 {
		class_idx := class_queue[0][0]
		class_dis := class_queue[0][1]
		class_queue = class_queue[1:]
		units, others, idxs, err := compressTimetable(
			tt_all, graph, class_idx, start_day, graph.Nodes[class_idx].Available, holidays, can_compress[class_idx],
		)
		if err != nil {
			return nil, cost_inf, errors.ErrorWrap(err)
		}
		log.Println("class id", graph.Nodes[class_idx].Id)
		log.Println("compress able", can_compress[class_idx])
		log.Println("idxs", idxs)
		for _, l := range idxs {
			i, j := l[0], l[1]
			log.Println("units", i, j, units[i][j])
		}

		cost := getCost(units, others, idxs, avoid, place_indexes, places)

		mv, err := timetableChangeSolver(cost, change_unit_idx, units)
		if err != nil {
			return nil, cost_inf, errors.ErrorWrap(err)
		}
		if c := getFinalCost(mv, cost); c < final_cost {
			final_cost = c
			final_move = mv
			final_units = units
		}
		for _, p := range graph.Nodes[class_idx].Parent {
			if class_used[p] {
				continue
			}
			class_used[p] = true
			class_queue = append(class_queue, [2]int{
				p, class_dis + 1,
			})
		}
	}
	var res []usecase.Timetable
	for i, v := range final_move {
		to := final_move[(i+1)%len(final_move)]
		vi, vj := v[0], v[1]
		ti, tj := to[0], to[1]
		day := start_day.AddDate(0, 0, ti)
		for _, idx := range final_units[vi][vj] {
			t := tt_all[idx]
			t.Day = day
			t.FramePeriod = tj
			t.FrameDayWeek = int(day.Weekday()) - 1
			t.FrameId = t.FrameDayWeek*P + t.FramePeriod
			res = append(res, t)
		}
	}
	return res, final_cost, nil
}
