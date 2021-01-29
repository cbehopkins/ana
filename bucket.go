package ana

import (
	"fmt"
	"log"
	"sync"
)

func rnInt(rn rune) (int, bool) {
	tmp := int(rn) - 65
	if (tmp < 0) || (tmp > 57) {
		log.Println("Invalid Input", rn)
		return -1, false
	}
	// Convert to lowercase
	if tmp > 31 {
		tmp -= 32
	}
	if tmp > 25 {
		log.Fatal("You gave me an invalid character", rn)
	}
	return tmp, true
}

// ArrayBucket is the count of
// the number of times a letter appears in a ref string
type ArrayBucket []int

// NewArrayBucket creates a new ab from a string
func NewArrayBucket(in string) ArrayBucket {
	retArray := make([]int, 26)

	for _, v := range in {
		tmp, ok := rnInt(v)
		if ok {
			retArray[tmp]++
		}
	}
	return retArray
}
func (b ArrayBucket) got(v int) bool {
	if (v >= len(b)) || (v < 0) {
		fmt.Println("got has been asked for:", v)
		return false
	}
	if b[v] > 0 {
		b[v]--
		return true
	}
	return false

}
func (b ArrayBucket) tM(candidate string) bool {
	for _, v := range candidate {
		tmp, ok := rnInt(v)
		if ok && !b.got(tmp) {
			// if we can't successfully take a token from the bucket, then fail
			return false
		}
	}
	return true
}
func (b ArrayBucket) testMatch(candidate string) bool {
	// We make a copy because each test subtracts one from the total
	tmpBucket := b.copy()
	return tmpBucket.tM(candidate)
}
func (b ArrayBucket) copy() ArrayBucket {
	retArray := make([]int, 26)
	copy(retArray, b)
	return retArray
}

// Worker is what you call on the bucket
// to test each input possible word
// against those in the created bucket
func (b ArrayBucket) Worker(inChan <-chan string, outChan chan<- string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for cand := range inChan {
			//Test against the bucket if
			// it the cand can be made from available tokens
			if b.testMatch(cand) {
				outChan <- cand
			}
		}
		wg.Done()
	}()
}

// BlockWorker is what you call on the bucket
// to test each input possible word
// against those in the created bucket
// except it works on blocks of strings at a time
// to reduce GC and channel congestion
func (b ArrayBucket) BlockWorker(inChan <-chan []string, outChan chan<- string, wg *sync.WaitGroup) {

	for candA := range inChan {
		for _, cand := range candA {
			if b.testMatch(cand) {
				outChan <- cand
			}
		}
	}
	wg.Done()
}

// A map bucket is a map from a letter to the numnber of
// times we see that letter (for the input string)
type mapBucket struct {
	bukMap map[int]int
}

func newMapBucket(in string) *mapBucket {
	itm := new(mapBucket)
	itm.bukMap = make(map[int]int)

	for _, v := range in {
		tmp, ok0 := rnInt(v)
		if ok0 {
			v, ok := itm.bukMap[tmp]
			if !ok {
				v = 1
			} else {
				v++
			}
			itm.bukMap[tmp] = v
		}
	}
	return itm
}

func (b mapBucket) testMatch(candidate string) bool {
	tmpBucket := b.copy()
	return tmpBucket.tM(candidate)
}
func (b mapBucket) copy() ArrayBucket {
	retArray := make([]int, 26)
	for key, val := range b.bukMap {
		retArray[key] = val
	}
	return retArray
}

var arrayPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		var tmp ArrayBucket
		tmp = make(ArrayBucket, 26)
		return &tmp
	},
}

type cacheBucket ArrayBucket

func newCacheBucket(in string) *cacheBucket {
	itm := make(cacheBucket, 26)
	tmp := NewArrayBucket(in)
	cnt := copy(itm, tmp)
	if cnt != 26 {
		log.Fatal("Init Copy failed from", itm, tmp)
	}
	return &itm
}

func (b cacheBucket) testMatch(candidate string) bool {
	tmpBucket := b.Copy()
	result := tmpBucket.tM(candidate)
	arrayPool.Put(tmpBucket)
	return result
}
func (b cacheBucket) Copy() *ArrayBucket {
	retArray := arrayPool.Get().(*ArrayBucket)
	cnt := copy(*retArray, b)
	if cnt != 26 {
		log.Fatal("Copy failed from", *retArray, b)
	}
	return retArray
}
