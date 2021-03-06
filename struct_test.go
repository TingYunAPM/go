// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"fmt"
	"strings"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(md5sum("1234567890"))
}
func TestJsonDecode(t *testing.T) {
	x := jsonDecodeArray("[1,2,3,4,\"axxxx\"]")
	var a interface{} = x
	fmt.Println(strings.Replace(fmt.Sprintf("address = %p\n", a), "0x", "", -1))
	fmt.Println(x)
}
func TestFileEnum(t *testing.T) {
	dir := "c:\\"
	count, err := elementCount(dir)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("count of dir \"%s\" = %d\n", dir, count)
	}
}
func TestBinarySearch(t *testing.T) {
	a := []float64{1, 1.4, 3, 100, 120, 120.1, 120.5}
	b := []float64{1, 1.4, 3, 100, 120, 120.1, 0.5, 1.1, 1.6, 3.1, 101, 120.05, 121}
	fmt.Println(a)
	for i := range b {
		fmt.Printf("search %f %d\n", b[i], binarySearch(a, b[i]))
	}
	for i := range a {
		fmt.Println(round(a[i]))
	}
}

func TestQuartileP2(t *testing.T) {
	quartileList := []float64{1.0 / 7 * 2, 1.0 / 7 * 3, 1.0 / 7 * 4, 1.0 / 7 * 5}

	testdata1 := []float64{3110, 3770, 3990, 3990, 3000, 3880, 3330, 3440, 3440, 3000, 3550, 3440, 3220, 3330, 3220, 3110, 3770, 3110, 3880, 3330, 3440, 3440, 3000, 3110, 3220, 3440, 3000, 3550, 3880, 3990, 3990, 3550, 3110, 3330, 3330, 3000, 3880, 3660, 3110, 3990, 3110, 3990, 3440, 3330, 3660, 3660, 3000, 3220, 3220, 3330, 3660, 3990, 3440, 3220, 3220, 3440, 3440, 3880, 3990, 3550, 3440, 3110, 3330, 3440, 3110, 3770, 3330, 3220, 3440, 3440, 3440, 3220, 3660, 3000, 3440, 3550, 3220, 3220, 3110, 3660, 3110, 3880, 3550, 3770, 3220, 3440, 3220, 3110, 3220, 3220, 3110, 3880, 3550, 3770, 3770, 3550, 3990, 3220, 3880, 3660, 3110, 3550, 3990, 3110, 3550, 3880, 3110, 3770, 3770, 3110, 3440, 3770, 3000, 3110, 3990, 3110, 3550, 3770, 3990, 3000, 3990, 3110, 3550, 3880, 3110, 3880, 3550, 3220, 3550, 3550, 3770, 3990, 3550, 3220, 3330, 3990, 3770, 3770, 3550, 3770, 3880, 3660, 3330, 3990, 3990, 3000, 3110, 3440, 3220, 3000, 3550, 3770, 3110, 3550, 3330, 3990, 3110, 3330, 3550, 3440, 3220, 3550, 3220, 3330, 3000, 3660, 3110, 3770, 3110, 3660, 3550, 3440, 3990, 3330, 3990, 3550, 3990, 3550, 3330, 3770, 3550, 3550, 3000, 3440, 3000, 3330, 3440, 3110, 3880, 3110, 3550, 3660, 3990, 3220, 3330, 3330, 3000, 3660, 3660, 3220}
	testdata2 := []float64{3028.0, 2211.0}
	testdata3 := []float64{2, 1, 6, 6, 4, 9, 5, 6, 2, 7, 2, 4, 9, 7, 8, 4, 1, 8, 8, 8, 3, 5, 4, 1, 9, 5, 5, 6, 2, 0, 8, 5, 3, 6, 1, 4, 8, 0, 0, 1, 1, 3, 9, 9, 6, 7, 1, 5, 8, 7, 6, 9, 3, 1, 2, 4, 8, 2, 3, 4, 7, 2, 9, 2, 3, 9, 8, 7, 6, 1, 8, 5, 0, 8, 3, 2, 8, 0, 1, 7, 4, 9, 9, 3, 6, 2, 3, 0, 1, 4, 9, 3, 0, 9, 6, 1, 9, 7, 1, 9, 3, 4, 4, 0, 9, 6, 2, 6, 1, 4, 9, 9, 2, 2, 3, 0, 2, 0, 4, 4, 1, 9, 2, 7, 8, 9, 6, 5, 6, 5, 1, 4, 3, 6, 4, 7, 6, 6, 5, 6, 2, 2, 5, 4, 8, 6, 3, 4, 7, 2, 4, 8, 7, 0, 0, 0, 9, 6, 5, 2, 8, 3, 2, 1, 9, 2, 4, 0, 3, 2, 3, 6, 6, 6, 4, 8, 1, 0, 7, 7, 2, 8, 5, 1, 3, 0, 5, 3, 3, 3, 3, 8, 8, 7, 9, 1, 3, 3, 1, 1, 0, 5, 2, 2, 4, 9, 3, 3, 5, 7, 4, 0, 7, 4, 2, 6, 3, 2, 5, 4, 9, 0, 8, 8, 0, 6, 7, 0, 2, 3, 3, 4, 7, 9, 9, 7, 8, 5, 1, 4, 5, 0, 8, 5, 8, 7, 0, 7, 3, 9, 5, 0, 7, 1, 2, 6, 8, 3, 3, 6, 0, 6, 0, 0, 4, 5, 6, 3, 6, 8, 6, 3, 2, 8, 9, 1, 9, 3, 8, 6, 3, 5, 9, 0, 3, 6, 2, 9, 1, 1, 0, 6, 4, 1, 0, 9, 3, 2, 9, 5, 6, 3, 7, 8, 3, 4, 1, 0, 8, 1, 3, 0, 3, 3, 9, 9, 7, 2, 1, 3, 5, 6, 6, 9, 5, 1, 9, 8, 8, 7, 0, 7, 3, 9, 3, 1, 6, 1, 7, 3, 3, 3, 9, 9, 8, 4, 3, 8, 1, 2, 0, 1, 9, 6, 3, 2, 2, 5, 5, 7, 3, 4, 2, 2, 7, 5, 4, 7, 0, 6, 4, 3, 6, 4, 2, 9, 3, 4, 8, 7, 8, 2, 1, 6, 6, 7, 0, 7, 8, 4, 8, 0, 1, 6, 9, 9, 5, 5, 7, 5, 6, 7, 4, 8, 7, 6, 7, 1, 1, 1, 1, 9, 2, 2, 0, 3, 4, 0, 6, 9, 9, 1, 8, 5, 0, 5, 5, 4, 7, 6, 5, 6, 2, 1, 2, 5, 0, 6, 3, 7, 3, 6, 1, 3, 1, 0, 5, 3, 9, 2, 9, 9, 8, 0, 9, 5, 3, 9, 6, 8, 0, 0, 2, 5, 3, 1, 3, 4, 2, 9, 4, 3, 0, 8, 1, 7, 0, 5, 3, 9, 5, 3, 4, 4, 1, 7, 5, 5, 8, 9, 8, 0, 1, 0, 1, 9, 7, 8, 2, 3, 4, 7, 5, 3, 8, 7, 8, 4, 1, 3, 6, 6, 0, 8, 8, 1, 3, 5, 2, 6, 0, 1, 2, 1, 5, 3, 5, 0, 7, 0, 2, 3, 9, 8, 2, 5, 8, 4, 8, 9, 8, 7, 2, 7, 7, 1, 2, 3, 7, 9, 7, 4, 5, 2, 6, 2, 3, 8, 8, 8, 0, 8, 4, 7, 9, 2, 6, 7, 5, 1, 3, 0, 4, 1, 3, 8, 2, 8, 1, 0, 5, 6, 3, 5, 6, 7, 9, 2, 4, 5, 9, 2, 9, 5, 0, 6, 1, 1, 0, 3, 9, 2, 8, 6, 6, 8, 6, 3, 9, 0, 0, 7, 1, 9, 9, 6, 4, 7, 3, 0, 1, 8, 1, 6, 5, 7, 3, 9, 7, 0, 3, 7, 3, 3, 6, 0, 6, 3, 3, 4, 7, 4, 1, 3, 9, 3, 2, 2, 5, 0, 5, 2, 5, 1, 2, 3, 9, 0, 9, 8, 7, 9, 2, 0, 9, 8, 9, 0, 5, 4, 4, 4, 4, 3, 2, 9, 0, 2, 5, 8, 0, 9, 4, 6, 5, 0, 2, 1, 8, 1, 4, 8, 4, 2, 0, 2, 9, 7, 7, 7, 1, 2, 1, 3, 3, 5, 1, 9, 2, 3, 2, 7, 6, 5, 6, 9, 5, 0, 7, 9, 8, 3, 4, 5, 1, 6, 4, 6, 9, 4, 5, 0, 0, 0, 6, 8, 3, 4, 0, 7, 8, 6, 9, 8, 8, 9, 8, 7, 8, 0, 6, 5, 8, 5, 4, 6, 3, 5, 4, 1, 0, 9, 7, 2, 4, 1, 7, 9, 3, 4, 1, 5, 7, 9, 5, 8, 9, 4, 6, 1, 3, 5, 8, 4, 0, 4, 1, 4, 2, 4, 3, 0, 0, 8, 9, 5, 2, 7, 8, 2, 6, 1, 9, 8, 9, 9, 0, 7, 9, 6, 8, 8, 8, 7, 4, 7, 4, 3, 3, 2, 5, 9, 5, 3, 3, 2, 3, 1, 0, 5, 7, 7, 2, 5, 9, 3, 7, 4, 1, 3, 1, 5, 4, 0, 5, 6, 9, 1, 5, 2, 4, 5, 8, 5, 3, 0, 1, 3, 3, 4, 0, 1, 0, 8, 8, 0, 7, 9, 1, 1, 0, 5, 2, 1, 2, 8, 2, 6, 6, 8, 2, 7, 9, 5, 5, 1, 2, 0, 0, 5, 9, 8, 8, 1, 5, 7, 2, 1, 2, 3, 5, 5, 3, 9, 4, 0, 5, 2, 2, 5, 0, 0, 2, 5, 3, 3, 5, 2, 8, 9, 8, 8, 6, 7, 7, 9, 6, 7, 0, 6, 0, 6, 4, 8, 7, 8, 7, 4, 4, 1, 0, 4, 1, 9, 3, 2, 1, 7, 3, 1, 1, 2, 3, 2, 9, 5, 9, 1, 8, 4, 6, 6, 1, 0, 9, 9, 3, 4, 9, 9, 0, 5, 7, 0, 5, 8, 2, 1, 3, 4, 7, 6, 0, 3, 0, 7, 3, 0, 1, 3, 0, 1, 0, 5, 2, 3, 3, 4, 3, 2}
	testdata4 := []float64{0.0, 3009.0, 3046.0, 3070.0, 3102.0, 3119.0, 3139.0, 3150.0, 3163.0, 3179.0, 3228.0}
	testdatas := [][]float64{testdata1, testdata2, testdata3, testdata4}

	print := func(t *quartileP2) {
		p := t.Markers()

		for i := 0; i < t.MarkCount; i++ {
			fmt.Print(p[i], ",")
		}
		fmt.Print("\n")
	}
	for i := 0; i < 4; i++ {
		qp := &quartileP2{}
		qp.Init(quartileList)
		testdata := testdatas[i]
		for j := 0; j < len(testdata); j++ {
			qp.Add(testdata[j])
		}
		print(qp)
	}
	//	fmt.Println("-----------------------------------------------------")
	//	qp := &quartileP2{}
	//	qp.Init(quartileList)
	//	testdata := []float64{3110, 3770, 3990, 3990, 3000, 3880, 3330, 3440, 3440, 3000, 3550, 3440, 3220, 3330, 3220, 3110, 3770, 3110, 3880, 3330, 3440, 3440, 3000, 3110, 3220, 3440, 3000, 3550, 3880, 3990}
	//	printHead := func(x int) {
	//		if x < 10 {
	//			fmt.Print("0")
	//		}
	//		if x < 100 {
	//			fmt.Print("0")
	//		}
	//		fmt.Print(x, ": ")
	//	}
	//	for j := 0; j < len(testdata); j++ {
	//		qp.Add(testdata[j])
	//		printHead(j)
	//		for x := 0; x < qp.MarkCount; x++ {
	//			fmt.Print(qp._P2_n[x], ",")
	//		}
	//		fmt.Print("\n")
	//		printHead(j)
	//		for x := 0; x < qp.MarkCount; x++ {
	//			fmt.Print(qp._MarkersY[x], ",")
	//		}
	//		fmt.Print("\n")
	//		//		fmt.Print(j, ":")
	//		//		print(qp)
	//	}
}
