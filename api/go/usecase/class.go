package usecase

import (
	"fmt"

	"timetable/errors"
	"timetable/library/bitset"
)

// 子のクラスの和集合が親クラスとなっているか判定するための要素
type ClassFill struct {
	Board int
	Piece int
}

type ClassNode struct {
	Id         int
	Name       string
	Available  []bool
	Child      []int
	Parent     []int
	ParentFill []ClassFill
}

type ClassGraph struct {
	// DAG
	Id2index    map[int]int
	Nest        []int
	Nodes       []ClassNode
	descendants []bitset.Bitset
}

type ClassEdge struct {
	FromId int
	ToId   int
	Board  int // -1 を許す。
	Piece  int
}

func (graph *ClassGraph) Valid() bool {
	// DAG checker
	graph.Nest = []int{}
	n := len(graph.Nodes)
	in_edge := make([]int, n)
	queue := make([]int, 0, n)
	for i := 0; i < n; i++ {
		in_edge[i] = len(graph.Nodes[i].Parent)
	}
	for i := 0; i < n; i++ {
		if in_edge[i] == 0 {
			queue = append(queue, i)
			graph.Nest = append(graph.Nest, i)
		}
	}
	for i := 0; i < n; i++ {
		if len(queue) == 0 {
			return false
		}
		id := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		for _, c := range graph.Nodes[id].Child {
			in_edge[c]--
			if in_edge[c] == 0 {
				queue = append(queue, c)
			}
		}
	}
	return true
}

func (graph *ClassGraph) initDescendants() {
	logger.Debug("test")
	n := len(graph.Nodes)
	vec := []int{}
	deg := make([]int, n)
	graph.descendants = make([]bitset.Bitset, n)
	for i := range graph.Nodes {
		deg[i] = len(graph.Nodes[i].Child)
		if deg[i] == 0 {
			vec = append(vec, i)
		}
	}
	for len(vec) > 0 {
		x := vec[len(vec)-1]
		vec = vec[:len(vec)-1]
		deg[x]--
		if deg[x] > 0 {
			continue
		}
		des := graph.descendants
		des[x] = bitset.NewBitset(len(graph.Nodes))
		for _, c := range graph.Nodes[x].Child {
			des[x] = des[x].Or(des[c])
		}
		des[x].Set(x, true)
		vec = append(vec, graph.Nodes[x].Parent...)
	}
}

// graph に辺を加える
// 頂点は存在している必要
// database には追加しない
func (graph *ClassGraph) AddEdge(from_id int, to_id int, board_id int, peice_id int) error {
	from, ok1 := graph.Id2index[from_id]
	to, ok2 := graph.Id2index[to_id]
	n := len(graph.Nodes)
	if !ok1 || !ok2 {
		return errors.NewError("sql, edge error, ", from, to)
	}
	if from < 0 || n <= from || to < 0 || n <= to || from == to {
		return errors.NewError(fmt.Sprintf(
			"sql, edge error, from: %d, to: %d",
			from, to,
		))
	}
	graph.Nodes[from].Child = append(graph.Nodes[from].Child, to)
	graph.Nodes[to].Parent = append(graph.Nodes[to].Parent, from)
	graph.Nodes[to].ParentFill = append(graph.Nodes[to].ParentFill, ClassFill{
		Board: board_id,
		Piece: peice_id,
	})
	return nil
}

// database から グラフを取ってくる
func NewClassGraph() (*ClassGraph, error) {
	g, err := Db_class.GetClassGraph()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	if !g.Valid() {
		return nil, errors.NewError("database graph is not valid")
	}
	g.initDescendants()
	return g, nil
}

// database, graph 両方に追加する
func (graph *ClassGraph) Append(nodes []ClassNode, edges []ClassEdge) error {
	for _, node := range nodes {
		_, ok := graph.Id2index[node.Id]
		if ok {
			return errors.NewError("duplicate id", node)
		}
		graph.Id2index[node.Id] = len(graph.Nodes)
		graph.Nodes = append(graph.Nodes, node)
	}
	for _, e := range edges {
		if err := graph.AddEdge(e.FromId, e.ToId, e.Board, e.Piece); err != nil {
			return errors.ErrorWrap(err)
		}
	}
	if !graph.Valid() {
		return errors.NewError("graph is not valid")
	}
	err := Db_class.InsertClassroom(nodes)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	err = Db_class.InsertClassEdge(edges)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	graph.initDescendants()
	return nil
}

func (graph *ClassGraph) GetEdgeArray() []ClassEdge {
	var edges []ClassEdge
	for _, node := range graph.Nodes {
		for _, idx := range node.Child {
			edges = append(edges, ClassEdge{
				FromId: node.Id,
				ToId:   graph.Nodes[idx].Id,
			})
		}
	}
	return edges
}

func (graph *ClassGraph) Id2IndexArray(node_ids []int) ([]int, error) {
	var res []int
	for _, id := range node_ids {
		idx, ok := graph.Id2index[id]
		if !ok {
			return nil, errors.NewError("no id")
		}
		res = append(res, idx)
	}
	return res, nil
}

func (graph *ClassGraph) Index2IdArray(node_idx []int) ([]int, error) {
	var res []int
	for _, idx := range node_idx {
		if idx < 0 || len(graph.Nodes) <= idx {
			return nil, errors.NewError("index error")
		}
		res = append(res, idx)
	}
	return res, nil
}

func (graph *ClassGraph) GetDescendants(node_idxs []int) ([]int, error) {
	b := bitset.NewBitset(len(graph.Nodes))
	for _, idx := range node_idxs {
		b = b.Or(graph.descendants[idx])
	}
	var res []int
	for i := range graph.Nodes {
		if b.Test(i) {
			res = append(res, i)
		}
	}
	return res, nil
}
func (graph *ClassGraph) GetAncestors(node_idxs []int) []int {
	used := make([]bool, len(graph.Nodes))
	vec := make([]int, len(node_idxs))
	var res []int
	copy(vec, node_idxs)
	for len(vec) > 0 {
		idx := vec[len(vec)-1]
		vec = vec[:len(vec)-1]
		if used[idx] {
			continue
		}
		used[idx] = true
		res = append(res, idx)
		vec = append(vec, graph.Nodes[idx].Parent...)
	}
	return res
}

func (graph *ClassGraph) NodeIn(a int, b int) bool {
	// a,b: index, b が a の子孫かどうか
	return graph.descendants[a].Test(b)
}

// func (graph *ClassGraph) NodeCross(a int, b int) bool {
// 	// a,b: index, a の子孫と b の子孫が被っているかどうか
// 	bt := graph.descendants[a].And(graph.descendants[b])
// 	return bt.TestAll()
// }

func GetGraph() ([]ClassNode, []ClassEdge, error) {
	var nodes []ClassNode
	var edges []ClassEdge
	graph, err := NewClassGraph()
	if err != nil {
		return nodes, edges, errors.ErrorWrap(err)
	}
	nodes = graph.Nodes
	edges = graph.GetEdgeArray()
	return nodes, edges, nil
}

func PostGraph(nodes []ClassNode, edges []ClassEdge) error {
	graph, err := NewClassGraph()
	if err != nil {
		return errors.ErrorWrap(err)
	}
	err = graph.Append(nodes, edges)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}
