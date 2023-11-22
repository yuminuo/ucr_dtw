package UcrDtw

import (
	"math"
	"sort"
)

func PrepareQuery(queryArr []float64, queryLength int64, wrappingWindow float64) *Query {
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

	sort.Sort(sort.Reverse(IndexArray(Q_tmp)))

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

func FindSimilar(data Queue, query *Query, min_bsf float64, epoch int, step int64) ([]LocationDtw, int64, float64) {
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
	buffer = make([]float64, epoch, epoch)
	u_buff = make([]float64, epoch, epoch)
	l_buff = make([]float64, epoch, epoch)

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
				buffer[k] = buffer[int64(epoch)-query.m+1+k]
			}
		}

		ep = query.m - 1
		for ep < int64(epoch) {
			if data.Count() == 0 {
				break
			}
			d = data.Dequeue()
			buffer[ep] = d
			ep++
		}

		// Data are read in chunk of size epoch.
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

								dist = dtw(tz, query.q, cb, query.m, query.r, bsf)

								cur_loc := it*(int64(epoch)-query.m+1) + i - query.m + 1

								if dist < min_bsf {

									results = append(results, LocationDtw{cur_loc, dist})

								}

								if dist < bsf {
									// Update bsf
									// loc is the real starting location of the nearest neighbor in the file
									bsf = dist
									loc = cur_loc

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

			if ep < int64(epoch) {
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
		if curLoc.Index > preLocation.Index+step {
			selected = append(selected, preLocation)
			preLocation = curLoc
			continue
		}

		if curLoc.Value < preLocation.Value {
			preLocation = curLoc
		}
	}

	selected = append(selected, preLocation)

	return selected, loc, bsf
}
