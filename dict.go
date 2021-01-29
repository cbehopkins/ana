package ana

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

// ReadDict is called to read in the dictionary
// and output the words line by line
func ReadDict(fname string) (outChan chan string) {
	outChan = make(chan string, 10)
	go readDict(fname, outChan)
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

// ReadBlockDict is called to read in the dictionary
// and output the words line by line
// It uses arrays of strings (blocks) to reduce GC and channel congestion
func ReadBlockDict(fname string, blkSize int) (outChan chan []string) {
	f, err := os.Open(fname)
	check(err)
	outChan = make(chan []string, 16)
	go readReaderDict(f, outChan, blkSize)
	return
}

// ReadBaDict is called to read in the dictionary
// and output the words line by line
// It uses arrays of strings (blocks) to reduce GC and channel congestion
func ReadBaDict(ba []byte, blkSize int) (outChan chan []string) {
	outChan = make(chan []string, 16)
	r := bytes.NewReader(ba)
	go readReaderDict(r, outChan, blkSize)
	return
}

func readReaderDict(r io.Reader, outChan chan []string, blkSize int) {
	scanner := bufio.NewScanner(r)
	var tmpArr []string
	for scanner.Scan() {
		tmpArr = append(tmpArr, scanner.Text())
		if len(tmpArr) >= blkSize {
			fmt.Println("Sending:", tmpArr)
			outChan <- tmpArr
			tmpArr = make([]string, 0)
		}
	}
	if len(tmpArr) > 0 {
		outChan <- tmpArr
	}
	if err := scanner.Err(); err != nil {
		//fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Println("Finished reading in dict, Closing")
	close(outChan)
}
