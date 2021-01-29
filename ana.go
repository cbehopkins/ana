package ana

import (
	"log"
	"sort"
	"strings"
	"sync"
)

// AnagramWord return a list of results
// from a input word and data as a byte array
func AnagramWord(refString string, data []byte) Results {

	refString = strings.Replace(refString, "\n", "", -1)
	resultChan := HelperBa(data, refString, 4)
	results := make(Results, 0)
	for res := range resultChan {
		results = append(results, Result(res))
	}

	sort.Sort(sort.Reverse(results))
	return results
}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Helper is the function an external program is expected to use
// give it a file to read a refString to look for
// and the number of parallel routines it is allowed to use
// and it will return a channel of words that match
func Helper(filename string, refString string, parCnt int) (dstChan chan string) {
	dstChan = make(chan string)
	go func() {
		var wg sync.WaitGroup
		blkSize := 128
		srcChan := ReadBlockDict(filename, blkSize)
		catBucket := NewArrayBucket(refString)
		wg.Add(parCnt)
		for i := 0; i < parCnt; i++ {
			catBucket.BlockWorker(srcChan, dstChan, &wg)
		}
		wg.Wait()
		close(dstChan)
	}()
	return dstChan
}

// HelperBa is the function an external program is expected to use
// give it a byte array to read a refString to look for
// and the number of parallel routines it is allowed to use
// and it will return a channel of words that match
func HelperBa(ba []byte, refString string, parCnt int) (dstChan chan string) {
	dstChan = make(chan string)

	go func() {
		blkSize := 128
		var wg sync.WaitGroup
		srcChan := ReadBaDict(ba, blkSize)
		catBucket := NewArrayBucket(refString)
		wg.Add(parCnt)
		for i := 0; i < parCnt; i++ {
			catBucket.BlockWorker(srcChan, dstChan, &wg)
		}
		wg.Wait()
		close(dstChan)
	}()
	return dstChan
}
