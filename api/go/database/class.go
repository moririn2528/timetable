package database

import (
	"strings"

	"timetable/errors"
	"timetable/usecase"
)

type Classroom struct {
	Id        int
	Name      string
	Available string
}

type ClassEdge struct {
	From  int
	To    int
	Board int
	Piece int
}

type DatabaseClass struct {
}

func getClassroom() ([]Classroom, error) {
	var res []Classroom
	rows, err := db.Query("SELECT id, name, available FROM classroom")
	if err != nil {
		return res, errors.ErrorWrap(err)
	}
	for rows.Next() {
		var c Classroom
		err := rows.Scan(&c.Id, &c.Name, &c.Available)
		if err != nil {
			return res, errors.ErrorWrap(err)
		}
		res = append(res, c)
	}
	return res, nil
}

func getClassEdge() ([]ClassEdge, error) {
	var res []ClassEdge
	rows, err := db.Query("SELECT from_id, to_id, board_id, piece_id FROM class_struct_edge")
	if err != nil {
		return res, errors.ErrorWrap(err)
	}
	for rows.Next() {
		var c ClassEdge
		err := rows.Scan(&c.From, &c.To, &c.Board, &c.Piece)
		if err != nil {
			return res, errors.ErrorWrap(err)
		}
		res = append(res, c)
	}
	return res, nil
}

func (*DatabaseClass) InsertClassroom(classes []usecase.ClassNode) error {
	query_arg := make([]interface{}, 0, 2*len(classes))
	query_str := make([]string, 0, len(classes))
	for _, c := range classes {
		query_arg = append(query_arg, c.Id, c.Name)
		query_str = append(query_str, "(?,?)")
	}
	_, err := db.Exec("INSERT INTO classroom(id, name) VALUES "+
		strings.Join(query_str, ","), query_arg...)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func (*DatabaseClass) InsertClassEdge(edges []usecase.ClassEdge) error {
	query_arg := make([]interface{}, 0, 2*len(edges))
	query_str := make([]string, 0, len(edges))
	for _, e := range edges {
		query_arg = append(query_arg, e.FromId, e.ToId, e.Board, e.Piece)
		query_str = append(query_str, "(?,?,?,?)")
	}
	_, err := db.Exec("INSERT INTO classroom(from_id, to_id, board_id, peice_id) VALUES "+
		strings.Join(query_str, ","), query_arg...)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func (*DatabaseClass) GetClassGraph() (*usecase.ClassGraph, error) {
	sql_classroom, err := getClassroom()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	sql_class_edge, err := getClassEdge()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	n := len(sql_classroom)
	graph := &usecase.ClassGraph{
		Nodes:    make([]usecase.ClassNode, 0, n),
		Id2index: make(map[int]int),
	}
	for i, c := range sql_classroom {
		graph.Id2index[c.Id] = i
		node := &usecase.ClassNode{
			Id:        c.Id,
			Name:      c.Name,
			Available: make([]bool, len(c.Available)),
		}
		for i, ca := range c.Available {
			if ca == '1' {
				node.Available[i] = true
			}
		}
		graph.Nodes = append(graph.Nodes, *node)
	}
	for _, e := range sql_class_edge {
		if err := graph.AddEdge(e.From, e.To, e.Board, e.Piece); err != nil {
			return nil, errors.ErrorWrap(err)
		}
	}
	if !graph.Valid() {
		return nil, errors.NewError("graph is invalid")
	}
	return graph, nil
}
