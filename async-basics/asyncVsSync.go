package main

import (
    "fmt"
    "sync"
    "time"
)

func wait(n int) int {
    // wait n seconds
    durr := time.Duration(n) * time.Second
    fmt.Println("waiting for:", durr)
    time.Sleep(durr)
    return n
}

func syncro(durs []int) {
    start := time.Now()

    // Loop over the calls and append results to slice
    results := make([]int, 0)
    for idx, val := range durs {
        fmt.Println(idx, fmt.Sprintf("seconds:%v", val))
        result := wait(val)
        results = append(results, result)
        fmt.Println("got result:", result)
    }

    fmt.Println("Ran in:", time.Since(start))
    fmt.Println("final results:", results)
}


func asyncWait(n int, c chan int, wg *sync.WaitGroup) {
    // Mark this branch as done after the function ends
    // (has added the result to the channel after waiting)
    defer wg.Done()

    // wait n seconds
    durr := time.Duration(n) * time.Second
    fmt.Println("waiting for:", durr)
    time.Sleep(durr)
    fmt.Println("Got result", n)
    c <- n
}

func async(durs []int) {
    c := make(chan int)
    var wg sync.WaitGroup

    start := time.Now()

    // Create branches
    for idx, val := range durs {
        fmt.Println(idx, fmt.Sprintf("seconds:%v", val))

        wg.Add(1)
        go asyncWait(val, c, &wg)
    }

    // Wait for completion
    go func() {
        fmt.Println("Waiting on group!!")
        wg.Wait()
        close(c)
    }()

    // As results stream in, add them to a slice
    results := make([]int, 0)
    for result := range c {
        fmt.Println("Got chan result:", result)
        results = append(results, result)
    }
    fmt.Println("Ran in:", time.Since(start))
    fmt.Println("final results", results)
}

func main() {
    durs := []int{3, 1, 2}

    fmt.Println("Running syncro...")
    syncro(durs)

    fmt.Println("\nRunning ayncro...")
    async(durs)
}
