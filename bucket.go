package ana

import (
	"bufio"
	"log"
	"os"
	"sync"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func ReadDict(fname string) (out_chan chan string) {
	out_chan = make(chan string, 10)
	go readDict(fname, out_chan)
	return
}
func readDict(fname string, outChan chan string) {
	// Extract messages from the log and put them
	// onto the lm.messageChan
	//fmt.Println("Reading from log:", fname)
	f, err := os.Open(fname)
	check(err)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		outChan <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		//fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	close(outChan)
}
func dummyDict(arr []string, outChan chan string) {
	for _, st := range arr {
		outChan <- st
	}
	close(outChan)
}
func dummyBlockDict(arr []string, outChan chan []string, cnt int) {

	for len(arr) > 0 {
		//log.Println("Array:", arr)
		if cnt > len(arr) {
			cnt = len(arr)
		}
		srA := arr[:cnt]
		arr = arr[cnt:]
		//log.Println("Sending:", srA)
		outChan <- srA
	}
	close(outChan)
}

type arrayBucket []int

func NewArrayBucket(in string) arrayBucket {
	retArray := make([]int, 26)

	for _, v := range in {
		tmp := rnInt(v)
		retArray[tmp] += 1
	}

	return retArray
}
func (b arrayBucket) got(v int) bool {
	if b[v] > 0 {
		b[v] -= 1
		return true
	} else {
		return false
	}
}
func (b arrayBucket) tM(candidate string) bool {
	for _, v := range candidate {
		tmp := rnInt(v)
		if !b.got(tmp) {
			return false
		}
	}
	return true
}
func (b arrayBucket) testMatch(candidate string) bool {
	tmpBucket := b.Copy()
	return tmpBucket.tM(candidate)
}
func (b arrayBucket) Copy() arrayBucket {
	retArray := make([]int, 26)
	copy(retArray, b)
	return retArray
}
func (b arrayBucket) Worker(inChan <-chan string, outChan chan<- string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for cand := range inChan {
			if b.testMatch(cand) {
				outChan <- cand
			}
		}
		wg.Done()
	}()
}
func (b arrayBucket) BlockWorker(inChan <-chan []string, outChan chan<- string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for candA := range inChan {
			for _, cand := range candA {
				//log.Println("tetsing", cand)
				if b.testMatch(cand) {
					outChan <- cand
				}
			}
		}
		wg.Done()
	}()
}

type mapBucket struct {
	bukMap map[int]int
}

func NewMapBucket(in string) *mapBucket {
	itm := new(mapBucket)
	itm.bukMap = make(map[int]int)

	for _, v := range in {
		tmp := rnInt(v)
		v, ok := itm.bukMap[tmp]
		if !ok {
			v = 1
		} else {
			v += 1
		}
		itm.bukMap[tmp] = v
	}
	return itm
}

func (b mapBucket) testMatch(candidate string) bool {
	tmpBucket := b.Copy()
	return tmpBucket.tM(candidate)
}
func (b mapBucket) Copy() arrayBucket {
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
		var tmp arrayBucket
		tmp = make(arrayBucket, 26)
		return &tmp
	},
}

type cacheBucket arrayBucket

func NewCacheBucket(in string) *cacheBucket {
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
func (b cacheBucket) Copy() *arrayBucket {
	retArray := arrayPool.Get().(*arrayBucket)
	cnt := copy(*retArray, b)
	if cnt != 26 {
		log.Fatal("Copy failed from", *retArray, b)
	}
	return retArray
}
