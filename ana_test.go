package ana

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/pkg/profile"
)

func generateRandomString(leng int) string {
	tmp := ""
	for i := 0; i < leng; i++ {
		rn := rune(rand.Intn(26) + 65)
		tmp += string(rn)
	}
	return tmp
}

func rnInt(rn rune) int {
	tmp := int(rn) - 65
	if (tmp < 0) || (tmp > 57) {
		log.Println("Invalid Input")
	}
	// Convert to lowercase
	if tmp > 31 {
		tmp -= 32
	}
	if tmp > 25 {
		log.Fatal("You gave me an invalid character", rn)
	}
	return tmp
}

func TestCategorise(t *testing.T) {
	testString := "azAZ"
	// Categorise the input string into buckets
	catBucket := NewArrayBucket(testString)
	for i, v := range catBucket {
		if i == 0 || i == 25 {
			if v != 2 {
				log.Fatal("expected 2")
			}
		} else {
			if v != 0 {
				log.Fatal("expected 2")
			}
		}
	}
}

func TestMatchPass(t *testing.T) {
	// This is the string we are looking for words in
	// i.e. does this have the letters in passString inside it
	refString := "abcdefffst"
	passStrings := []string{
		"cat",
		"bat",
		"faff",
		"bead",
	}

	// categorise the reference
	catBucket := NewArrayBucket(refString)
	for _, passString := range passStrings {
		match := catBucket.testMatch(passString)
		if match {
			log.Println("Success:", passString)
		} else {
			log.Fatal("Fail:", passString)
		}
	}
}
func TestMatchFail(t *testing.T) {
	// This is the string we are looking for words in
	// i.e. does this have the letters in passString inside it
	refString := "abcdeffffst"
	failStrings := []string{
		"dead",
		"stack",
		"dad",
	}

	// categorise the reference
	catBucket := NewArrayBucket(refString)
	for _, failString := range failStrings {
		match := catBucket.testMatch(failString)
		if !match {
			log.Println("Success:", failString)
		} else {
			log.Fatal("Fail:", failString)
		}
	}
}
func TestMapMatchPass(t *testing.T) {
	// This is the string we are looking for words in
	// i.e. does this have the letters in passString inside it
	refString := "abcdefffst"
	passStrings := []string{
		"cat",
		"bat",
		"faff",
		"bead",
	}

	// categorise the reference
	catBucket := NewMapBucket(refString)
	for _, passString := range passStrings {
		match := catBucket.testMatch(passString)
		if match {
			log.Println("Success:", passString)
		} else {
			log.Fatal("Fail:", passString)
		}
	}
}
func TestMapMatchFail(t *testing.T) {
	// This is the string we are looking for words in
	// i.e. does this have the letters in passString inside it
	refString := "abcdeffffst"
	failStrings := []string{
		"dead",
		"stack",
		"dad",
	}

	// categorise the reference
	catBucket := NewMapBucket(refString)
	for _, failString := range failStrings {
		match := catBucket.testMatch(failString)
		if !match {
			log.Println("Success:", failString)
		} else {
			log.Fatal("Fail:", failString)
		}
	}
}
func TestCacheMatchPass(t *testing.T) {
	// This is the string we are looking for words in
	// i.e. does this have the letters in passString inside it
	refString := "abcdefffst"
	passStrings := []string{
		"cat",
		"bat",
		"faff",
		"bead",
	}

	// categorise the reference
	catBucket := NewCacheBucket(refString)
	for _, passString := range passStrings {
		match := catBucket.testMatch(passString)
		if match {
			log.Println("Success:", passString)
		} else {
			log.Fatal("Fail:", passString)
		}
	}
}
func TestCacheMatchFail(t *testing.T) {
	// This is the string we are looking for words in
	// i.e. does this have the letters in passString inside it
	refString := "abcdeffffst"
	failStrings := []string{
		"dead",
		"stack",
		"dad",
	}

	// categorise the reference
	catBucket := NewCacheBucket(refString)
	for _, failString := range failStrings {
		match := catBucket.testMatch(failString)
		if !match {
			log.Println("Success:", failString)
		} else {
			log.Fatal("Fail:", failString)
		}
	}
}
func initWordList(numItems int) []string {
	wordList := make([]string, 0, numItems)
	for i := 0; i < numItems; i++ {
		wrdLen := rand.Intn(3) + 5
		word := generateRandomString(wrdLen)
		wordList = append(wordList, word)

	}
	return wordList
}
func BenchmarkArrayBucket(b *testing.B) {
	defer profile.Start().Stop()
	//defer profile.Start(profile.MemProfile).Stop()
	b.StopTimer()
	numItems := b.N

	refString := generateRandomString(9)
	wordList := initWordList(numItems)
	b.StartTimer()
	catBucket := NewArrayBucket(refString)

	for _, passString := range wordList {
		catBucket.testMatch(passString)
	}
}
func BenchmarkMapBucket(b *testing.B) {
	b.StopTimer()
	numItems := b.N
	wordList := initWordList(numItems)

	refString := generateRandomString(9)
	b.StartTimer()
	catBucket := NewMapBucket(refString)
	for _, passString := range wordList {
		catBucket.testMatch(passString)
	}
}
func BenchmarkCacheBucket(b *testing.B) {
	b.StopTimer()
	numItems := b.N
	wordList := initWordList(numItems)

	refString := generateRandomString(9)
	b.StartTimer()
	catBucket := NewCacheBucket(refString)
	for _, passString := range wordList {
		catBucket.testMatch(passString)
	}
}
func parRun(parCnt int, b *testing.B) {
	numItems := b.N
	wordList := initWordList(numItems)
	refString := generateRandomString(9)
	srcChan := make(chan string)
	dstChan := make(chan string)
	var wg sync.WaitGroup
	go dummyDict(wordList, srcChan)
	go func() {
		for _ = range dstChan {
			//fmt.Printf("Word %v, matched against %v\n", v, refString)
		}
	}()
	b.ResetTimer()
	catBucket := NewArrayBucket(refString)
	for i := 0; i < parCnt; i++ {
		catBucket.Worker(srcChan, dstChan, &wg)
	}
	wg.Wait()
}
func BenchmarkParallelBucket_1(b *testing.B) {
	numWorkers := []int{1, 2, 4, 8, 16, 32, 64, 128}
	for _, workerCnt := range numWorkers {
		wcs := strconv.Itoa(workerCnt)

		b.Run(wcs, func(br *testing.B) { parRun(workerCnt, br) })
	}
}

func parBlockRun(parCnt int, blkLen int, b *testing.B) {
	numItems := b.N
	wordList := initWordList(numItems)
	refString := generateRandomString(9)
	srcChan := make(chan []string)
	dstChan := make(chan string)
	var wg sync.WaitGroup
	//log.Println("gnerating array with Items", numItems)
	go dummyBlockDict(wordList, srcChan, blkLen)
	go func() {
		for _ = range dstChan {
			//fmt.Printf("Word %v, matched against %v\n", v, refString)
		}
	}()
	b.ResetTimer()
	catBucket := NewArrayBucket(refString)
	for i := 0; i < parCnt; i++ {
		catBucket.BlockWorker(srcChan, dstChan, &wg)
	}
	wg.Wait()
}
func BenchmarkParallelBlockBucket(b *testing.B) {
	//numWorkers := []int{1}
	//sizeBlocks := []int{1}
	//numWorkers := []int{1, 2, 4, 8, 16, 32, 64, 128}
	//sizeBlocks := []int{1, 2, 4, 8, 16, 64, 128, 1024}
	numWorkers := []int{4, 8, 16, 32, 64, 128}
	sizeBlocks := []int{128}
	for _, workerCnt := range numWorkers {
		wcs := strconv.Itoa(workerCnt)
		for _, sb := range sizeBlocks {
			sbs := strconv.Itoa(sb)
			runStr := "W:" + wcs + "_" + "S:" + sbs
			b.Run(runStr, func(br *testing.B) {
				//log.Println("Starting block run with", br.N)
				parBlockRun(workerCnt, sb, br)
			})
		}
	}
}

//func ExampleParam() {
//  prm := Param{Name: "fred", Value: "steve"}
//  output, err := xml.MarshalIndent(prm, "", "  ")
//  check(err)
//  fmt.Println("", string(output))
//  //Output:
//  // <param>
//  //   <name>fred</name>
//  //   <value>steve</value>
//  // </param>
//}
