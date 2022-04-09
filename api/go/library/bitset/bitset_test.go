package bitset

import (
	"math/rand"
	"testing"
)

func TestSetnTest(t *testing.T) {
	ns := []int{1, 5, 10, 20, 100, 200, 10000}
	is := []int{0, 2, 5, 10, 20, 100, 200}
	for _, n := range ns {
		bt := NewBitset(n)
		for _, i := range is {
			if n <= i {
				break
			}
			bt.Set(i, true)
			if !bt.Test(i) {
				t.Errorf("bt.Test(%v): %v", i, bt.Test(i))
			}
			j := rand.Intn(n)
			if i != j && bt.Test(j) {
				t.Errorf("bt.Test(%v): %v, n: %v, i: %v", j, bt.Test(j), n, i)
			}
			if !bt.TestAll() {
				t.Error("bt.TestAll(): false")
			}
			bt.Set(i, false)
			if bt.Test(i) {
				t.Errorf("bt.Test(%v): %v", i, bt.Test(i))
			}
			if bt.TestAll() {
				t.Error("bt.TestAll(): true")
			}
		}
	}
}

func TestAndOr(t *testing.T) {
	create := func(n int) (Bitset, []bool) {
		a := make([]bool, n)
		b := NewBitset(n)
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				b.Set(i, true)
				a[i] = true
			}
		}
		return b, a
	}
	min := func(a int, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a int, b int) int {
		if a < b {
			return b
		}
		return a
	}
	for i := 0; i < 5; i++ {
		n := rand.Intn(1000)
		m := rand.Intn(1000)
		b1, a1 := create(n)
		b2, a2 := create(m)
		b := b1.And(b2)
		if b.size != min(n, m) {
			t.Errorf("b.size: %v, n: %v, m: %v", b.size, n, m)
		}
		for j := 0; j < min(n, m); j++ {
			if b.Test(j) != (a1[j] && a2[j]) {
				t.Errorf("b.Test(%v): %v, a1[j]: %v, a2[j]: %v", j, b.Test(j), a1[j], a2[j])
			}
		}
		b = b1.Or(b2)
		if b.size != max(n, m) {
			t.Errorf("b.size: %v, n: %v, m: %v", b.size, n, m)
		}
		for len(a1) < max(n, m) {
			a1 = append(a1, false)
		}
		for len(a2) < max(n, m) {
			a2 = append(a2, false)
		}
		for j := 0; j < max(n, m); j++ {
			if b.Test(j) != (a1[j] || a2[j]) {
				t.Errorf("b.Test(%v): %v, a1[j]: %v, a2[j]: %v, n: %v, m: %v", j, b.Test(j), a1[j], a2[j], n, m)
			}
		}
	}
}
