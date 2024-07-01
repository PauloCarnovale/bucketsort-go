package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

func generateRandomSlice(size, interval int) []int {
	slice := make([]int, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(interval)
	}
	return slice
}

func main() {
	// Abrir arquivo para salvar resultados
	file, err := os.Create("resultados.txt")
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return
	}
	defer file.Close()

	for i := 0; i < 16; i++ {
		exec(1000000, 1000000000, i+1, file)
	}
}

func exec(numberCount, interval, nrBuckets int, file *os.File) {
	runtime.GOMAXPROCS(nrBuckets)
	var wg sync.WaitGroup
	wg.Add(nrBuckets)

	numbersToSort := generateRandomSlice(numberCount, interval)
	buckets := make([][]int, nrBuckets)
	finalArray := make([]int, 0, numberCount)

	fmt.Printf("Sorting %d numbers from 0 to %d using %d buckets\n", numberCount, interval-1, nrBuckets)

	startTime := time.Now()
	for i := 0; i < nrBuckets; i++ {
		go func(i int) {
			defer wg.Done()
			bucketStart := interval / nrBuckets * i
			bucketEnd := interval / nrBuckets * (i + 1)

			for _, number := range numbersToSort {
				if number >= bucketStart && number < bucketEnd {
					buckets[i] = append(buckets[i], number)
				}
			}
			sort.Ints(buckets[i])
		}(i)
	}

	wg.Wait()
	for _, bucket := range buckets {
		finalArray = append(finalArray, bucket...)
	}

	timeElapsed := time.Since(startTime)
	fmt.Printf("Time elapsed: %v\n", timeElapsed)

	// Salvar resultados no arquivo
	file.WriteString(fmt.Sprintf("Sorting %d numbers from 0 to %d using %d buckets\n", numberCount, interval-1, nrBuckets))
	file.WriteString(fmt.Sprintf("Time elapsed: %v\n", timeElapsed))
}

