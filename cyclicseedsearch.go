package main

import (
	"crypto/rand"
	"fmt"
	"runtime"
)

type Seed struct {
	a, b, c, d uint32
	overlap    int
	at         int
}

var globalBests = make([]*Seed, 9)
var bestFinish = 5

func main() {
	th := runtime.NumCPU()
	runtime.GOMAXPROCS(th)
	ch := make(chan Seed)
	for i := 0; i < th; i++ {
		go DoSearch(ch)
	}
	for {
		seed := <-ch
		if globalBests[seed.overlap] == nil || globalBests[seed.overlap].at < seed.at {
			globalBests[seed.overlap] = &seed
			fmt.Println(seed)
			if globalBests[seed.overlap].at == CheckLen {
				bestFinish = seed.overlap
			}
		}
	}
}

func xorShift(xs1, xs2, xs3, xs4 uint32) (uint32, uint32, uint32, uint32) {
	t, xs1, xs2, xs3 := xs1^(xs1<<11), xs2, xs3, xs4
	xs4 = xs4 ^ (xs4 >> 19) ^ t ^ (t >> 8)
	return xs1, xs2, xs3, xs4
}

func ru32() uint32 {
	b := make([]byte, 4)
	rand.Read(b)
	return uint32(b[0]) + 256*uint32(b[1]) + 256*256*uint32(b[2]) + 256*256*256*uint32(b[3])
}

// large enough for about 50M
const CheckLen = 400000

func DoSearch(out chan<- Seed) {

	// index[ position ][ value ] == list of indexes that had that position in that value
	index := make([][][]int, 10)
	for i := 0; i < 10; i++ {
		index[i] = make([][]int, 256)
		for j := 0; j < 256; j++ {
			index[i][j] = make([]int, 0)
		}
	}

	for {
		xs1i, xs2i, xs3i, xs4i := ru32(), ru32(), ru32(), ru32()
		xs1, xs2, xs3, xs4 := xs1i, xs2i, xs3i, xs4i
		overlaps := 0
		finished := true

		//reset index
		for i := 0; i < 10; i++ {
			for j := 0; j < 256; j++ {
				index[i][j] = make([]int, 0)
			}
		}
	top:
		for count := 0; count < CheckLen; count++ {
			//accumulator index
			accIdx := make(map[int]int)
			for i := 0; i < 10; i++ {
				xs1, xs2, xs3, xs4 = xorShift(xs1, xs2, xs3, xs4)
				v := int(xs4 & 255)
				for _, idx := range index[i][v] {
					if _, ok := accIdx[idx]; ok {
						accIdx[idx]++
					} else {
						accIdx[idx] = 1
					}
				}
				index[i][v] = append(index[i][v], count)
			}
			for _, c := range accIdx {
				if c > overlaps {
					if globalBests[overlaps] == nil || globalBests[overlaps].at < count {
						out <- Seed{xs1i, xs2i, xs3i, xs4i, overlaps, count}
					}
					if c >= bestFinish {
						finished = false
						break top
					}
					overlaps = c
				}
			}

		}
		if finished && (globalBests[overlaps] == nil || globalBests[overlaps].at < CheckLen) {
			out <- Seed{xs1i, xs2i, xs3i, xs4i, overlaps, CheckLen}
		}
	}
}
