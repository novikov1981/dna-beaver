package measurements

import (
	dna_beaver "github.com/novikov1981/dna-beaver"
	"github.com/novikov1981/dna-beaver/utils"
	"sync"
	"sync/atomic"
)

type Statistics struct {
	LinksCount                 map[string]int
	LinkMutex                  sync.Mutex
	Links, Oligs, WrongSymbols int32
}

func Measure(oo []string) (synthesisStatistic Statistics) {
	stat := Statistics{Oligs: int32(len(oo)),
		LinksCount: make(map[string]int, len(dna_beaver.ValidNotations))}
	// initialise map with right symbols
	for _, r := range dna_beaver.ValidNotations {
		stat.LinksCount[string(r)] = 0
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
				stat.LinkMutex.Lock()
				i, ok := stat.LinksCount[string(l)]
				if ok {
					stat.LinksCount[string(l)] = i + 1
				} else {
					stat.WrongSymbols += 1
				}
				stat.LinkMutex.Unlock()
			}
			wg.Done()
		}(o)
	}
	wg.Wait()
	return stat
}

//
//func (r *Statistic) MeasureOne(o string) (statisticOne Statistic) {
//	dna := strings.ToUpper(utils.ExtractDna(o))
//	count := 0
//	for _, o := range ValidNotations {
//		c := strings.Count(dna, string(o))
//		if c > 0 {
//			switch string(o) {
//			case "A":
//				statisticOne.a += c
//			case "C":
//				statisticOne.c += c
//			case "G":
//				statisticOne.g += c
//			case "T":
//				statisticOne.t += c
//			case "R":
//				statisticOne.r += c
//				//					dA += cf / 2
//				//					dG += cf / 2
//			case "Y":
//				statisticOne.y += c
//				//					dC += cf / 2
//				//					dT += cf / 2
//			case "K":
//				statisticOne.k += c
//				//					dG += cf / 2
//				//					dT += cf / 2
//			case "M":
//				statisticOne.m += c
//				//					dA += cf / 2
//				//					dC += cf / 2
//			case "S":
//				statisticOne.s += c
//				//					dG += cf / 2
//				//					dC += cf / 2
//			case "W":
//				statisticOne.w += c
//				//					dA += cf / 2
//				//					dT += cf / 2
//			case "B":
//				statisticOne.b += c
//				//					dC += cf / 3
//				//					dG += cf / 3
//				//					dT += cf / 3
//			case "D":
//				statisticOne.d += c
//				//					dA += cf / 3
//				//					dG += cf / 3
//				//					dT += cf / 3
//			case "H":
//				statisticOne.h += c
//				//					dA += cf / 3
//				//					dC += cf / 3
//				//					dG += cf / 3
//			case "V":
//				statisticOne.v += c
//				//					dA += cf / 3
//				//					dC += cf / 3
//				//					dG += cf / 3
//			case "N":
//				statisticOne.n += c
//				//					dA += cf / 4
//				//					dC += cf / 4
//				//					dG += cf / 4
//				//					dT += cf / 4
//			}
//		}
//		count += c
//	}
//
//	statisticOne.wrongSymbols += len(dna) - count
//	statisticOne.allLinks += len(dna)
//
//	return statisticOne
//}
//
