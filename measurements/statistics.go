package measurements

import (
	. "github.com/novikov1981/dna-beaver"
	"github.com/novikov1981/dna-beaver/utils"
	"sync"
	"sync/atomic"
)

type Statistics struct {
	LinksCount                 map[string]int
	AmeditCount                map[string]float32
	mtx                        sync.Mutex
	Links, Oligs, WrongSymbols int32
}

func Measure(oo []string) (synthesisStatistic Statistics) {
	stat := Statistics{
		Oligs:       int32(len(oo)),
		LinksCount:  make(map[string]int, len(ValidNotations)),
		AmeditCount: make(map[string]float32, len(Amedits)),
	}
	// initialise link count map with right symbols
	for _, r := range ValidNotations {
		stat.LinksCount[string(r)] = 0
	}
	// initialise amedit count map with amedits symbols
	for _, s := range Amedits {
		stat.LinksCount[s] = 0
	}
	// prepare wait group
	wg := sync.WaitGroup{}
	// iterate over oligs and count symbols
	for _, o := range oo {
		wg.Add(1)
		go func(olig string) {
			dna := utils.ExtractDna(olig)
			atomic.AddInt32(&stat.Links, int32(len(dna)))
			for _, l := range dna {
				ls := string(l)
				stat.mtx.Lock()
				i, ok := stat.LinksCount[ls]
				if ok {
					stat.LinksCount[ls] = i + 1
					countAmedits(ls, &stat)
				} else {
					stat.WrongSymbols += 1
				}
				stat.mtx.Unlock()
			}
			wg.Done()
		}(o)
	}
	wg.Wait()
	return stat
}

func countAmedits(link string, stat *Statistics) {
	switch link {
	case "R":
		stat.AmeditCount[DA] += 0.5
		stat.AmeditCount[DG] += 0.5
	case "Y":
		stat.AmeditCount[DC] += 0.5
		stat.AmeditCount[DT] += 0.5
	case "K":
		stat.AmeditCount[DG] += 0.5
		stat.AmeditCount[DT] += 0.5
	case "M":
		stat.AmeditCount[DC] += 0.5
		stat.AmeditCount[DA] += 0.5
	case "S":
		stat.AmeditCount[DC] += 0.5
		stat.AmeditCount[DG] += 0.5

	case "W":
		stat.AmeditCount[DA] += 0.5
		stat.AmeditCount[DT] += 0.5

	case "B":
		stat.AmeditCount[DC] += 0.33
		stat.AmeditCount[DG] += 0.33
		stat.AmeditCount[DT] += 0.33

	case "D":
		stat.AmeditCount[DA] += 0.33
		stat.AmeditCount[DG] += 0.33
		stat.AmeditCount[DT] += 0.33

	case "H":
		stat.AmeditCount[DA] += 0.33
		stat.AmeditCount[DC] += 0.33
		stat.AmeditCount[DG] += 0.33

	case "V":
		stat.AmeditCount[DA] += 0.33
		stat.AmeditCount[DG] += 0.33
		stat.AmeditCount[DC] += 0.33

	case "N":
		stat.AmeditCount[DA] += 0.25
		stat.AmeditCount[DG] += 0.25
		stat.AmeditCount[DC] += 0.25
		stat.AmeditCount[DT] += 0.25

	default:
		stat.AmeditCount[link] += 1
	}

}
