package UcrDtw

import (
	"golang.org/x/exp/errors/fmt"
	"math"
)

func dist(x float64, y float64) float64 {
	return (x - y) * (x - y)
}

func min_float(x float64, y float64) float64 {
	if x < y {
		return x
	}

	return y
}

func max_float(x float64, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func max_int(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
func min_int(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

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

func dtw(A []float64, B []float64, cb []float64, m int64, r int64, bsf float64) float64 {

	var cost []float64
	var cost_prev []float64
	var cost_tmp []float64
	var i, j, k int64
	var x, y, z, min_cost float64

	cost = make([]float64, 2*r+1, 2*r+1)
	cost_prev = make([]float64, 2*r+1, 2*r+1)
	for k = 0; k < 2*r+1; k++ {
		cost[k] = math.Inf(1)
		cost_prev[k] = math.Inf(1)

	}

	for i = 0; i < m; i++ {
		k = max_int(0, r-i)
		min_cost = math.Inf(1)

		for j = max_int(0, i-r); j < min_int(m-1, i+r); j, k = j+1, k+1 {
			// Initialize all row and column
			if i == 0 && j == 0 {
				cost[k] = dist(A[0], B[0])
				min_cost = cost[k]
				continue
			}

			if (j-1 < 0) || (k-1 < 0) {
				y = math.Inf(1)
			} else {
				y = cost[k-1]
			}

			if (i-1 < 0) || (k+1 > 2*r) {
				x = math.Inf(1)
			} else {
				x = cost_prev[k+1]
			}

			if (i-1 < 0) || (j-1 < 0) {
				z = math.Inf(1)
			} else {
				z = cost_prev[k]
			}

			// Classic DTW calculation
			cost[k] = min_float(min_float(x, y), z) + dist(A[i], B[j])

			// Find minimum cost in row for early abandoning (possibly to use column instead of row).
			if cost[k] < min_cost {
				min_cost = cost[k]
			}

		}

		// We can abandon early if the current cummulative distace with lower bound together are larger than bsf
		if i+r < m-1 && min_cost+cb[i+r+1] >= bsf {
			return min_cost + cb[i+r+1]
		}

		// Move current array to previous array.
		cost_tmp = cost
		cost = cost_prev
		cost_prev = cost_tmp
	}
	k--

	final_dtw := cost_prev[k]

	return final_dtw

}

func prepareQuery(queryArr []float64, queryLength int64, wrappingWindow float64, EPOCH int) *Query {
	var d float64
	var order []int64
	var q, u, l, qo, uo, lo []float64

	var i int64
	var ex, ex2, mean, std float64
	var m int64 = -1
	var r int64 = -1

	var Q_tmp []IndexValue

	m = queryLength

	if wrappingWindow >= 0 {
		R := wrappingWindow
		if R <= 1 {
			r = int64(math.Floor(R * float64(m)))
		} else {
			r = int64(math.Floor(R))
		}

	}

	q = make([]float64, m, m)
	qo = make([]float64, m, m)
	uo = make([]float64, m, m)
	lo = make([]float64, m, m)
	order = make([]int64, m, m)
	Q_tmp = make([]IndexValue, m, m)

	u = make([]float64, m, m)
	l = make([]float64, m, m)

	i = 0
	ex, ex2 = 0, 0

	for itt := 0; i < m && itt < len(queryArr); itt++ {
		d = queryArr[itt]
		ex += d
		ex2 += d * d
		q[i] = d
		i++
	}

	mean = ex / float64(m)
	std = ex2 / float64(m)
	std = float64(math.Sqrt(std - mean*mean))
	for i = 0; i < m; i++ {
		q[i] = (q[i] - mean) / std
	}

	lower_upper_lemire(q, m, r, l, u)

	for i = 0; i < m; i++ {
		Q_tmp[i] = IndexValue{Index: i, Value: q[i]}
	}

	/*for {
		sort.SliceStable(Q_tmp, func(i, j int) bool {

			a := Q_tmp[i]
			b := Q_tmp[j]
			val := int16(math.Abs(b.Value)-math.Abs(a.Value))

			return  val>0
		})

		isSorted := sort.SliceIsSorted(Q_tmp, func(i, j int) bool {

			a := Q_tmp[i]
			b := Q_tmp[j]
			val := int16(math.Abs(b.Value)-math.Abs(a.Value))

			return  val>0
		})

		if isSorted{
			break
		}
	}*/

	/*for {
		sort.Sort(sort.Reverse(IndexArray(Q_tmp)))
		if sort.IsSorted(sort.Reverse(IndexArray(Q_tmp))) {
			break
		}
	}*/

	//Q_tmp = QuickSort(Q_tmp)

	fmt.Println()

	for i = 0; i < m; i++ {
		o := Q_tmp[i].Index
		order[i] = o
		qo[i] = q[o]
		uo[i] = u[o]
		lo[i] = l[o]
	}

	var query = Query{order, q, u, l, qo, uo, lo, m, r}
	return &query

}

func similiarity_finder(data Queue, query *Query, min_bsf float64, EPOCH int) ([]LocationDtw, int64, float64) {
	min_bsf = min_bsf * min_bsf
	var d float64
	var t []float64
	var i, j int64 = 0, 0
	var ex, ex2, mean, std float64
	var loc int64
	var bsf float64
	var kim, keogh, keogh2 int64
	var dist, lb_kim, lb_k, lb_k2 float64
	var buffer, u_buff, l_buff, tz, cb, cb1, cb2 []float64
	//var u_d, l_d [] float64

	cb = make([]float64, query.m, query.m)
	cb1 = make([]float64, query.m, query.m)
	cb2 = make([]float64, query.m, query.m)
	//u_d = make([]float64, query.m, query.m)
	//l_d = make([]float64, query.m, query.m)

	t = make([]float64, query.m*2, query.m*2)
	tz = make([]float64, query.m, query.m)
	buffer = make([]float64, EPOCH, EPOCH)
	u_buff = make([]float64, EPOCH, EPOCH)
	l_buff = make([]float64, EPOCH, EPOCH)

	done := false

	var it, ep, k, I int64

	bsf = math.Inf(1)

	for i = 0; i < query.m; i++ {
		cb[i] = 0
		cb1[i] = 0
		cb2[i] = 0
	}

	var results []LocationDtw

	for !done {
		ep = 0
		if it == 0 {
			for k = 0; k < query.m-1; k++ {
				if data.Count() > 0 {
					d = data.Dequeue()
					buffer[k] = d
				}
			}
		} else {
			for k = 0; k < query.m-1; k++ {
				buffer[k] = buffer[int64(EPOCH)-query.m+1+k]
			}
		}

		ep = query.m - 1
		for ep < int64(EPOCH) {
			if data.Count() == 0 {
				break
			}
			d = data.Dequeue()
			buffer[ep] = d
			ep++
		}

		// Data are read in chunk of size EPOCH.
		// When there is nothing to read, the loop is end.
		if ep <= query.m-1 {
			done = true
		} else {
			lower_upper_lemire(buffer, ep, query.r, l_buff, u_buff)
			ex = 0
			ex2 = 0

			for i = 0; i < ep; i++ {
				d = buffer[i]

				// Calcualte sum and sum square
				ex += d
				ex2 += d * d

				t[i%query.m] = d
				t[(i%query.m)+query.m] = d

				if i >= query.m-1 {
					mean = ex / float64(query.m)
					std = ex2 / float64(query.m)
					std = float64(math.Sqrt(std - mean*mean))

					// compute the start location of the data in the current circular array, t
					j = (i + 1) % query.m
					// the start location of the data in the current chunk
					I = i - (query.m - 1)

					// Use a constant lower bound to prune the obvious subsequence
					lb_kim = lb_kim_hierarchy(t, query.q, j, query.m, mean, std, bsf)

					if lb_kim < min_bsf {
						// Use a linear time lower bound to prune; z_normalization of t will be computed on the fly.
						// uo, lo are envelop of the query.
						lb_k = lb_keogh_cumulative(query.order, t, query.uo, query.lo, cb1, j, query.m, mean, std, bsf)

						if lb_k < min_bsf {
							// Take another linear time to compute z_normalization of t.
							// Note that for better optimization, this can merge to the previous function.
							for k = 0; k < query.m; k++ {
								tz[k] = (t[(k+j)] - mean) / std
							}

							// Use another lb_keogh to prune
							// qo is the sorted query. tz is unsorted z_normalized data.
							// l_buff, u_buff are big envelop for all data in this chunk

							//ArraySegment<double> l_buff_partial = new ArraySegment<double>(l_buff, 0, I);
							//ArraySegment<double> u_buff_partial = new ArraySegment<double>(u_buff, 0, I);

							lb_k2 = lb_keogh_data_cumulative(query.order, tz, query.qo, cb2, l_buff,
								u_buff, I, query.m, mean,
								std, bsf)

							if lb_k2 < min_bsf {
								// Choose better lower bound between lb_keogh and lb_keogh2 to be used in early abandoning DTW
								// Note that cb and cb2 will be cumulative summed here.
								if lb_k > lb_k2 {
									cb[query.m-1] = cb1[query.m-1]
									for k = query.m - 2; k >= 0; k-- {
										cb[k] = cb[k+1] + cb1[k]
									}

								} else {
									cb[query.m-1] = cb2[query.m-1]
									for k = query.m - 2; k >= 0; k-- {
										cb[k] = cb[k+1] + cb2[k]
									}

								}

								dist = dtw(tz, query.q, cb, query.m, query.r, bsf) //TODO maybe min_bsf

								cur_loc := it*(int64(EPOCH)-query.m+1) + i - query.m + 1
								if cur_loc == int64(756562) {
									fmt.Printf("In location distance is %v\n", math.Sqrt(dist))
								}

								if dist < min_bsf {

									results = append(results, LocationDtw{cur_loc, dist})

								}

								if dist < bsf {
									// Update bsf
									// loc is the real starting location of the nearest neighbor in the file
									bsf = dist
									loc = it*(int64(EPOCH)-query.m+1) + i - query.m + 1

								}

							} else {
								keogh2++
							}

						} else {
							keogh++
						}

					} else {
						kim++
					}

					ex -= t[j]
					ex2 -= t[j] * t[j]

				}

			}

			if ep < int64(EPOCH) {
				done = true
			} else {
				it++
			}

		}
	}

	if len(results) < 1 {
		return results, loc, bsf
	}

	var selected []LocationDtw

	preLocation := results[0]
	for ri := 1; ri < len(results); ri++ {
		var curLoc = results[ri]
		if curLoc.index > preLocation.index+20 {
			selected = append(selected, preLocation)
			preLocation = curLoc
			continue
		}

		if curLoc.value < preLocation.value {
			preLocation = curLoc
		}
	}

	selected = append(selected, preLocation)

	return selected, loc, bsf
}
