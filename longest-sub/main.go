package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/golang-collections/collections/set"
)

// Global sets to track already checked substrings
var mu sync.Mutex
var knownValid = set.New()
var knownInvalid = set.New()

func subInAll(sub string, strs []string, validity chan string, doneGroup *sync.WaitGroup) bool {
	mu.Lock()
	defer mu.Unlock()
	defer doneGroup.Done() // mark substring as done after

	// Do preliminary check
	if (knownValid.Has(sub)) {
		return true;
	}
	if (knownInvalid.Has(sub)) {
		return false;
	}

	inAll := true
	for _, s := range strs {
		if (!strings.Contains(s, sub)) {
			inAll = false
			break
		}
	}

	if (inAll) {
		validity <- sub
		knownValid.Insert(sub)
	} else {
		knownInvalid.Insert(sub)
	}

	return inAll
}

// Iterates over all possible substrings and dumps to channel if valid
func generateCandidates(idx int, strs []string, validCh chan string, finishedWg *sync.WaitGroup) {
	var doneGroup sync.WaitGroup
	base := strs[idx]

	// Just attempt every possible substring
	for i := 0; i < len(base); i++ {
		for j := 0; j < len(base); j++ {
			if j <= i {
				continue
			}

			sub := base[i:j]
			doneGroup.Add(1)
			subInAll(sub, strs, validCh, &doneGroup)
		}
	}

	// Mark the string as completed when all the subs finish
	go func() {
		doneGroup.Wait()
		finishedWg.Done()
	}()
}

func longest(strs []string) string {
	// Build the candidate branches
	validCh := make(chan string)
	var finishedGeneration sync.WaitGroup
	for idx := 0; idx < len(strs); idx++ {
		finishedGeneration.Add(1)
		go generateCandidates(idx, strs, validCh, &finishedGeneration)
	}

	// Close the channels when all candidates built
	go func() {
		finishedGeneration.Wait()
		close(validCh)
	}()

	// Process incoming valid strings, track the longest
	valid := make([]string, 0)
	var longestString string
	for s := range validCh {
		valid = append(valid, s)
		if (len(longestString) < len(s)) {
			longestString = s
		}
	}

	return longestString
}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func main() {
	// Process input
	strs := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		scanned := scanner.Text()
		if len(strings.Trim(scanned, "*")) > 0 {
			strs = append(strs, scanned)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(strs) <= 1 {
		log.Fatal("Passed in text must be longer than one line")
	}

	// Find longest substring
	start := time.Now()
	result := longest(strs)
	runtime := time.Since(start)

	// Log the result
	if result != "" {
		fmt.Println("\nGot this final result:")
		fmt.Println(result)
	} else {
		fmt.Println("\nNo string matches found.")
	}
	fmt.Println("\nFinished in:", runtime)
}