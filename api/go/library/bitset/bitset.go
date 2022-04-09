package bitset

const unit int = 64

type Bitset struct {
	size int
	arr  []int64
}

func NewBitset(n int) Bitset {
	if n <= 0 {
		n = 1
	}
	return Bitset{
		size: n,
		arr:  make([]int64, (n-1)/unit+1),
	}
}

func (bt *Bitset) Test(pos int) bool {
	if pos < 0 || bt.size <= pos {
		return false
	}
	i := pos / unit
	j := pos % unit
	return (bt.arr[i] & (int64(1) << j)) != int64(0)
}

func (bt *Bitset) TestAll() bool {
	for _, v := range bt.arr {
		if v != int64(0) {
			return true
		}
	}
	return false
}

func (bt *Bitset) Set(pos int, val bool) {
	if pos < 0 || bt.size <= pos {
		return
	}
	i := pos / unit
	j := pos % unit
	if val {
		bt.arr[i] |= int64(1) << j
	} else {
		bt.arr[i] &= ^(int64(1) << j)
	}
}

func (b1 *Bitset) And(b2 Bitset) Bitset {
	n := b1.size
	if b2.size < n {
		n = b2.size
	}
	bt := NewBitset(n)
	for i := 0; i < len(b1.arr) && i < len(b2.arr); i++ {
		bt.arr[i] = b1.arr[i] & b2.arr[i]
	}
	return bt
}

func (b1 *Bitset) Or(b2 Bitset) Bitset {
	n := b1.size
	if n < b2.size {
		n = b2.size
	}
	bt := NewBitset(n)
	for i := 0; i < len(b1.arr) || i < len(b2.arr); i++ {
		if i < len(b1.arr) && i < len(b2.arr) {
			bt.arr[i] = b1.arr[i] | b2.arr[i]
		} else if i < len(b1.arr) {
			bt.arr[i] = b1.arr[i]
		} else {
			bt.arr[i] = b2.arr[i]
		}
	}
	return bt
}
