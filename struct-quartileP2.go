// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"sort"
)

type quartileP2 struct {
	QCount        int
	MarkCount     int
	_Count        int
	_QuartileList []float64
	_MarkersX     []float64
	_MarkersY     []float64
	_P2_n         []int
	_ResultValue  []float64
}

func (q *quartileP2) Quantile(p float64) float64 {
	//		assert(_Count > 0);
	if q._Count < q.MarkCount {
		return q._MarkersY[int(p*float64(q._Count-1))]
	} else if p <= 0.0 {
		return q._MarkersY[0]
	} else if p >= 1 {
		return q._MarkersY[q.MarkCount-1]
	} else if index := binarySearch(q._MarkersX, p); index > -1 {
		return q._MarkersY[index]
	} else {
		left, right := -index-2, -index-1
		pl, pr := q._MarkersX[left], q._MarkersX[right]
		return (q._MarkersY[left]*(pr-p) + q._MarkersY[right]*(p-pl)) / (pr - pl)
	}
}

func (q *quartileP2) Count() int { return q._Count }
func (q *quartileP2) Init(quartileList []float64) {
	q.QCount = len(quartileList)
	q._QuartileList = make([]float64, q.QCount)
	for i := 0; i < q.QCount; i++ {
		q._QuartileList[i] = quartileList[i]
	}
	sort.Float64s(q._QuartileList)
	q.MarkCount = q.QCount*2 + 3
	q._MarkersX = make([]float64, q.MarkCount)
	q._MarkersY = make([]float64, q.MarkCount)
	q._ResultValue = make([]float64, q.MarkCount)
	q._P2_n = make([]int, q.MarkCount)
	for i := 0; i < q.MarkCount; i++ {
		q._MarkersY[i], q._ResultValue[i], q._P2_n[i] = 0.0, 0.0, i
	}
	q._Count = 0
	q._MarkersX[0] = 0.0
	for i := 0; i < q.QCount; i++ {
		var marker float64
		marker, q._MarkersX[i*2+2] = q._QuartileList[i], q._QuartileList[i]
		q._MarkersX[i*2+1] = (marker + q._MarkersX[i*2]) / 2
	}
	q._MarkersX[q.MarkCount-2] = (1 + q._QuartileList[q.QCount-1]) / 2
	q._MarkersX[q.MarkCount-1] = 1.0
}

func (q *quartileP2) Add(data float64) {
	index := q._Count
	q._Count++
	if index < q.MarkCount {
		q._MarkersY[index] = data
		if q._Count == q.MarkCount {
			sort.Float64s(q._MarkersY)
		}
	} else {
		k := binarySearch(q._MarkersY, data)
		if k < 0 {
			k = -(k + 1)
		}
		if k == 0 {
			q._MarkersY[0], k = data, 1
		} else if k == q.MarkCount {
			k, q._MarkersY[q.MarkCount-1] = q.MarkCount-1, data
		}
		for i := k; i < q.MarkCount; i++ {
			q._P2_n[i] += 1
		}

		quadPred := func(d int, i int) float64 {
			qi, qip1, qim1 := q._MarkersY[i], q._MarkersY[i+1], q._MarkersY[i-1]
			ni, nip1, nim1 := q._P2_n[i], q._P2_n[i+1], q._P2_n[i-1]
			a := float64(ni-nim1+d) * (qip1 - qi) / float64(nip1-ni)
			b := float64(nip1-ni-d) * (qi - qim1) / float64(ni-nim1)
			return qi + (float64(d)*(a+b))/float64(nip1-nim1)
		}
		linPred := func(d int, i int) float64 {
			qi, qipd := q._MarkersY[i], q._MarkersY[i+d]
			ni, nipd := q._P2_n[i], q._P2_n[i+d]
			return qi + float64(d)*(qipd-qi)/float64(nipd-ni)
		}
		for i := 1; i < q.MarkCount-1; i++ {
			n_ := q._MarkersX[i] * float64(index)
			di := n_ - float64(q._P2_n[i])
			//			if index == 28 && i == 8 {
			//				fmt.Print("_MarkersX[", i, "]=", q._MarkersX[i], ",n_=", n_, ",_P2_n[", i, "]=", q._P2_n[i], ",_P2_n[", i+1, "]=", q._P2_n[i+1], ",di=", di, "\n")
			//			}
			if (di-1.0 >= 0.000001 && q._P2_n[i+1]-q._P2_n[i] > 1) || (di+1.0 <= 0.000001 && q._P2_n[i-1]-q._P2_n[i] < -1) {
				d := 1
				if di < 0 {
					d = -1
				}
				qi_ := quadPred(d, i)
				if qi_ < q._MarkersY[i-1] || qi_ > q._MarkersY[i+1] {
					qi_ = linPred(d, i)
				}
				q._MarkersY[i] = qi_
				q._P2_n[i] += d
			}
		}
	}
}

func (q *quartileP2) Markers() []float64 {
	if q.MarkCount <= q._Count {
		return q._MarkersY
	}
	result := make([]float64, q.MarkCount)
	pw_q_copy := make([]float64, q.MarkCount)
	for i := 0; i < q.MarkCount; i++ {
		pw_q_copy[i] = q._MarkersY[i]
	}
	sort.Float64s(pw_q_copy)
	for i, j := q.MarkCount-q._Count, 0; i < q.MarkCount; i, j = i+1, j+1 {
		result[j] = pw_q_copy[i]
	}
	for i := 0; i < q.MarkCount; i++ {
		q._ResultValue[i] = result[round(float64((q._Count-1)*i*1.0)/float64(q.MarkCount-1))]
	}
	return q._ResultValue
}
