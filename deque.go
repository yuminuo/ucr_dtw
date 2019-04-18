package UcrDtw

type Deque struct {
	dq       []int64
	size     int64
	capacity int64
	f        int64
	r        int64
}

func (d *Deque) Empty() bool {
	if d.size == 0 {
		return true
	}
	return false
}

func (d *Deque) init(cap int64) {
	d.capacity = cap
	d.size = int64(0)
	d.dq = make([]int64, d.capacity, d.capacity)
	d.f = int64(0)
	d.r = d.capacity - 1

}

func (d *Deque) push_back(v int64) {
	d.dq[d.r] = v
	d.r--
	if d.r < 0 {
		d.r = d.capacity - 1
	}
	d.size++
}

func (d *Deque) pop_front() {
	d.f--
	if d.f < 0 {
		d.f = d.capacity - 1
	}
	d.size--
}

func (d *Deque) pop_back() {
	d.r = (d.r + 1) % d.capacity
	d.size--
}

func (d *Deque) front() int64 {
	aux := d.f - 1
	if aux < 0 {
		aux = d.capacity - 1
	}

	return d.dq[aux]

}

func (d *Deque) back() int64 {
	aux := (d.r + 1) % d.capacity
	return d.dq[aux]

}
