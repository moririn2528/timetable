package communicate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"timetable/errors"
	"timetable/library/logging"
	"timetable/usecase"
)

var (
	logger *logging.Logger = logging.NewLogger()
)

func ResponseJson(w http.ResponseWriter, v interface{}) error {
	res, err := json.Marshal(v)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, res, "", "  ")
	if err != nil {
		return errors.ErrorWrap(err)
	}
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	_, err = io.WriteString(w, buf.String())
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

type ClassNode struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ComposeClassNode(node usecase.ClassNode) ClassNode {
	return ClassNode{
		Id:   node.Id,
		Name: node.Name,
	}
}

func (node *ClassNode) parse() usecase.ClassNode {
	return usecase.ClassNode{
		Id:   node.Id,
		Name: node.Name,
	}
}

type ClassEdge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func ComposeClassEdge(edge usecase.ClassEdge) ClassEdge {
	return ClassEdge{
		From: edge.FromId,
		To:   edge.ToId,
	}
}

func (edge *ClassEdge) parse() usecase.ClassEdge {
	return usecase.ClassEdge{
		FromId: edge.From,
		ToId:   edge.To,
	}
}

type ClassRoom struct {
	Node []ClassNode `json:"nodes"`
	Edge []ClassEdge `json:"edges"`
}

func ComposeClassRoom(nodes []usecase.ClassNode, edges []usecase.ClassEdge) ClassRoom {
	ns := make([]ClassNode, 0, len(nodes))
	es := make([]ClassEdge, 0, len(edges))
	for _, n := range nodes {
		ns = append(ns, ComposeClassNode(n))
	}
	for _, e := range edges {
		es = append(es, ComposeClassEdge(e))
	}
	return ClassRoom{
		Node: ns,
		Edge: es,
	}
}

func get_class_structure(w http.ResponseWriter, req *http.Request) error {
	class_nodes, class_edges, err := usecase.GetGraph()
	if err != nil {
		return errors.ErrorWrap(err)
	}
	class := ComposeClassRoom(class_nodes, class_edges)
	err = ResponseJson(w, class)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func post_class_structure(w http.ResponseWriter, req *http.Request) error {
	var class ClassRoom
	err := json.NewDecoder(req.Body).Decode(&class)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	nodes := make([]usecase.ClassNode, 0, len(class.Node))
	edges := make([]usecase.ClassEdge, 0, len(class.Edge))
	for _, n := range class.Node {
		nodes = append(nodes, n.parse())
	}
	for _, e := range class.Edge {
		edges = append(edges, e.parse())
	}
	err = usecase.PostGraph(nodes, edges)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func Class_structure(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "GET":
		err = get_class_structure(w, req)
	case "POST":
		err = post_class_structure(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == nil {
		return
	}
	my_err, ok := err.(*errors.MyError)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("wrap error")
		return
	}
	w.WriteHeader(my_err.GetCode())
	logger.Error(my_err.Error())
}
