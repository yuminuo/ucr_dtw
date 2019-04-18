package UcrDtw

import "math"

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

		for j = max_int(0, i-r); j <= min_int(m-1, i+r); j, k = j+1, k+1 {
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
