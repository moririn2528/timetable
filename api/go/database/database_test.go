package database

import "testing"

/// データベース上のデータが想定される仮定を満たすかチェックする

func TestClassGraph(t *testing.T) {
	nodes, err := getClassroom()
	if err != nil {
		t.Fatal(err)
	}
	edges_id, err := getClassEdge()
	if err != nil {
		t.Fatal(err)
	}
	nodes_id_map := make(map[int]int)
	for i, c := range nodes {
		nodes_id_map[c.Id] = i
	}
	var edges []ClassEdge
	for _, e := range edges_id {
		f, ok := nodes_id_map[e.From]
		if !ok {
			t.Fatal("edge id error")
		}
		to, ok := nodes_id_map[e.To]
		if !ok {
			t.Fatal("edge id error")
		}
		edges = append(edges, ClassEdge{
			From: f,
			To:   to,
		})
	}
	periods, err := getFrames()
	if err != nil {
		t.Fatal(err)
	}
	n := len(nodes)

	path := make([][]int, n)
	in_deg := make([]int, n)
	for _, e := range edges {
		path[e.From] = append(path[e.From], e.To)
		in_deg[e.To]++
	}

	// DAG チェック
	{
		trail := make([]bool, n)
		used := make([]bool, n)
		// ループがなければ True
		var dfs func(int) bool
		dfs = func(x int) bool {
			if used[x] {
				return true
			}
			used[x] = true
			if trail[x] {
				return false
			}
			trail[x] = true
			for _, to := range path[x] {
				if !dfs(to) {
					return false
				}
			}
			trail[x] = false
			return true
		}
		for i := 0; i < n; i++ {
			if in_deg[i] > 0 {
				continue
			}
			if !dfs(i) {
				t.Errorf("class graph dag error, start index: %v", i)
				break
			}
		}
	}

	// 推移できる辺がないかチェック
	{
		for i, e := range edges {
			f := e.From
			l := len(path[f])
			if l == 0 {
				t.Fatalf("internal error, edge: %v", edges_id[i])
			}
			// 1 辺削除
			for i, v := range path[f] {
				if v != e.To {
					continue
				}
				path[f][i] = path[f][l-1]
				path[f] = path[f][:l-1]
			}
			if len(path[f])+1 != l {
				t.Fatalf("internal path error, edge: %v", edges_id[i])
			}
			used := make([]bool, n)
			vec := []int{f}
			for len(vec) > 0 {
				l := len(vec)
				x := vec[l-1]
				vec = vec[:l-1]
				if used[x] {
					continue
				}
				used[x] = true
				if x == e.To {
					t.Errorf("too many edges, unnecessary edge: %v", edges_id[i])
					break
				}
				vec = append(vec, path[x]...)
			}
		}
	}

	// available チェック
	{
		for i := 0; i < n; i++ {
			if len(nodes[i].Available) != len(periods) {
				t.Errorf("class available len error, len: %v, len(peroid): %v, class idx: %v", len(nodes[i].Available), len(periods), i)
			}
		}
		for k, e := range edges {
			a1 := nodes[e.From].Available
			a2 := nodes[e.To].Available
			if len(a1) != len(a2) {
				t.Error()
				continue
			}
			for i := 0; i < len(a1); i++ {
				if a1[i] == '1' && a2[i] == '0' {
					t.Errorf("class available error, edge: %v, index: %v", edges_id[k], i)
				}
			}
		}
	}
}

func TestTeacher(t *testing.T) {
	var db DatabaseAny
	teachers, err := db.GetTeacher()
	if err != nil {
		t.Fatal(err)
	}
	periods, err := getFrames()
	if err != nil {
		t.Fatal(err)
	}
	for _, edu := range teachers {
		if len(edu.Avoid) != len(periods) {
			t.Errorf("teacher avoid length error, len(t.Avoid): %v, len(periods): %v", len(edu.Avoid), len(periods))
		}
	}
}

func TestTimetable(t *testing.T) {
	classes, err := getClassroom()
	if err != nil {
		t.Fatal(err)
	}
	class_map := make(map[int]Classroom)
	for _, c := range classes {
		class_map[c.Id] = c
	}
	var dt DatabaseTimetable
	normal_table, err := dt.GetNomalTimetable(-1, []int{}, -1)
	if err != nil {
		t.Fatal(err)
	}
	for _, tab := range normal_table {
		c, ok := class_map[tab.ClassId]
		if !ok {
			t.Errorf("class id not found, class id: %v, normal timetable id: %v", c.Id, tab.Id)
			continue
		}
		f := tab.FrameId
		if c.Available[f] == '0' {
			t.Errorf("class available error,class: %v, timetable: %v", c, tab)
		}
	}
}
