package UcrDtw

func lower_upper_lemire(t []float64, len int64, r int64, l []float64, u []float64) {
	du := Deque{}
	du.init(2*r + 2)
	du.push_back(0)

	dl := Deque{}
	dl.init(2*r + 2)
	dl.push_back(0)

	var i int64 = 1
	for ; i < len; i++ {
		if i > r {
			u[i-r-1] = t[du.front()]
			l[i-r-1] = t[dl.front()]
		}

		if t[i] > t[i-1] {
			du.pop_back()
			for !du.Empty() && t[i] > t[du.back()] {
				du.pop_back()
			}
		} else {
			dl.pop_back()
			for !dl.Empty() && t[i] < t[dl.back()] {
				dl.pop_back()
			}

		}

		du.push_back(i)
		dl.push_back(i)

		if i == 2*r+1+du.front() {
			du.pop_front()
		} else {
			if i == 2*r+1+dl.front() {
				dl.pop_front()
			}
		}
	}

	for i = len; i < len+r+1; i++ {
		u[i-r-1] = t[du.front()]
		l[i-r-1] = t[dl.front()]

		if i-du.front() >= 2*r+1 {
			du.pop_front()
		}
		if i-dl.front() >= 2*r+1 {
			dl.pop_front()
		}
	}

}

func lb_kim_hierarchy(t []float64, q []float64, j int64, len int64, mean float64, std float64, bsf float64) float64 {
	// 1 point at front and back

	var d float64
	var lb float64
	x0 := (t[j] - mean) / std
	y0 := (t[(len-1+j)] - mean) / std
	lb = dist(x0, q[0]) + dist(y0, q[len-1])
	if lb > bsf {
		return lb
	}

	/// 2 points at front
	x1 := (t[(j+1)] - mean) / std
	d = min_float(dist(x1, q[0]), dist(x0, q[1]))
	d = min_float(d, dist(x1, q[1]))
	lb += d
	if lb >= bsf {
		return lb
	}

	/// 2 points at back
	y1 := (t[(len-2+j)] - mean) / std
	d = min_float(dist(y1, q[len-1]), dist(y0, q[len-2]))
	d = min_float(d, dist(y1, q[len-2]))
	lb += d
	if lb >= bsf {
		return lb
	}

	/// 3 points at front
	x2 := (t[(j+2)] - mean) / std
	d = min_float(dist(x0, q[2]), dist(x1, q[2]))
	d = min_float(d, dist(x2, q[2]))
	d = min_float(d, dist(x2, q[1]))
	d = min_float(d, dist(x2, q[0]))
	lb += d
	if lb >= bsf {
		return lb
	}

	/// 3 points at back
	y2 := (t[(len-3+j)] - mean) / std
	d = min_float(dist(y0, q[len-3]), dist(y1, q[len-3]))
	d = min_float(d, dist(y2, q[len-3]))
	d = min_float(d, dist(y2, q[len-2]))
	d = min_float(d, dist(y2, q[len-1]))
	lb += d

	return lb
}

func lb_keogh_cumulative(order []int64, t []float64, uo []float64, lo []float64, cb []float64,
	j int64, len int64, mean float64, std float64, best_so_far float64) float64 {
	var lb float64 = 0
	var x, d float64

	var i int64 = 0
	for ; i < len && lb < best_so_far; i++ {
		x = (t[(order[i]+j)] - mean) / std
		d = 0
		if x > uo[i] {
			d = dist(x, uo[i])
		} else {
			if x < lo[i] {
				d = dist(x, lo[i])
			}

		}
		lb += d
		cb[order[i]] = d
	}
	return lb
}

func lb_keogh_data_cumulative(order []int64, tz []float64, qo []float64, cb []float64, l []float64,
	u []float64, I int64, len int64, mean float64, std float64, best_so_far float64) float64 {
	var lb float64 = 0
	var uu, ll, d float64

	var i int64 = 0
	for ; i < len && lb < best_so_far; i++ {
		uu = (u[order[i]+I] - mean) / std
		ll = (l[order[i]+I] - mean) / std
		d = 0
		if qo[i] > uu {
			d = dist(qo[i], uu)
		} else {
			if qo[i] < ll {
				d = dist(qo[i], ll)
			}
		}
		lb += d
		cb[order[i]] = d
	}
	return lb
}
