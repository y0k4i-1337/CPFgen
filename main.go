// Generate one or more valid CPF numbers
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultThreads int = 20

var regionCodes = [10]string{
	"RS",
	"DF, GO, MS, MT, TO",
	"AC, AM, AP, PA, RO, RR",
	"CE, MA, PI",
	"AL, PB, PE, RN",
	"BA, SE",
	"MG",
	"ES, RJ",
	"SP",
	"PR, SC",
}

func main() {
	// command-line options
	var list = flag.Bool("l", false, "list regions and their codes")
	var output = flag.String("o", "", "output results to file")
	var heuristics = flag.Bool("e", false, "use heuristics when generating CPF")
	var format = flag.Int("f", 1, `output format (1: 11122233345,
  2: 111.222.333-45, 3: 111222333-45)`)
	var regions = flag.String("r", "0,1,2,3,4,5,6,7,8,9",
		"comma-separated list with region codes")
	var random = flag.Int("n", 0, "generate n random CPF numbers")
	var verbose = flag.Bool("v", false, "verbose")
	// channels
	jobs := make(chan [9]int)

	// Create a WaitGroup to synchronize the goroutines.
	var wg sync.WaitGroup

	flag.Parse()

	if *list {
		listRegions()
		return
	}

	// Remove files if they already exist
	if *output != "" {
		if _, err := os.Stat(*output); err == nil {
			os.Remove(*output)
		}
	}

	regionsSlice := strings.Split(*regions, ",")
	// TODO: check if regionsSlice contains only digits from 0 to 9
	var regionsInt = make([]int, len(regionsSlice))
	for i, c := range regionsSlice {
		regionsInt[i], _ = strconv.Atoi(c)
	}

	// Increment the WaitGroup for each goroutine we want to wait for.
	wg.Add(1)
	go produce(regionsInt, *heuristics, jobs, *random, &wg)

	go consume(0, jobs, *verbose, &wg, *format, *output)

	// Wait for all the goroutines to finish their work before proceeding.
	wg.Wait()
}

// print the regions related to the 9th digit of CPF
func listRegions() {
	fmt.Println("List of Regions:")
	for i, s := range regionCodes {
		fmt.Printf("\t%d\t%q\n", i, s)
	}
}

// produce combinations of numbers to construct CPF numbers later
func produce(regions []int, heuristics bool, jobs chan<- [9]int, random int, wg *sync.WaitGroup) {
	if random > 0 {
		randgen := rand.New(rand.NewSource(time.Now().UnixNano()))
		uniqueNumbers := make(map[string]bool)
		for len(uniqueNumbers) < random {
			firstEightDigits := randgen.Intn(100000000)
			lastDigit := regions[randgen.Intn(len(regions))]
			numberStr := fmt.Sprintf("%08d%d", firstEightDigits, lastDigit)

			uniqueNumbers[numberStr] = true
		}
		var number [9]int
		for numberStr := range uniqueNumbers {

			for i, digitRune := range numberStr {
				digit := int(digitRune - '0')
				number[i] = digit
			}
			jobs <- number
		}
	} else {

		var base [9]int
		var count [10]int

		for d1 := 0; d1 <= 9; d1++ {
			count[d1]++
			base[0] = d1
			for d2 := 0; d2 <= 9; d2++ {
				count[d2]++
				base[1] = d2
				for d3 := 0; d3 <= 9; d3++ {
					count[d3]++
					base[2] = d3
					for d4 := 0; d4 <= 9; d4++ {
						count[d4]++
						base[3] = d4
						for d5 := 0; d5 <= 9; d5++ {
							base[4] = d5
							if heuristics && count[d5] == 4 {
								continue
							}
							count[d5]++
							for d6 := 0; d6 <= 9; d6++ {
								base[5] = d6
								if heuristics && count[d6] == 4 {
									continue
								}
								count[d6]++
								for d7 := 0; d7 <= 9; d7++ {
									base[6] = d7
									if heuristics && count[d7] == 4 {
										continue
									}
									count[d7]++
									for d8 := 0; d8 <= 9; d8++ {
										base[7] = d8
										if heuristics && count[d8] == 4 {
											continue
										}
										count[d8]++
										for _, d9 := range regions {
											base[8] = d9
											if heuristics && count[d9] == 4 {
												continue
											}
											jobs <- base
										}
										count[d8]--
									}
									count[d7]--
								}
								count[d6]--
							}
							count[d5]--
						}
						count[d4]--
					}
					count[d3]--
				}
				count[d2]--
			}
			count[d1]--
		}
	}
	close(jobs)
}

// consume base numbers and append verification numbers
func consume(worker int, jobs <-chan [9]int,
	verbose bool, wg *sync.WaitGroup, format int, output string) {
	defer wg.Done()
	for base := range jobs {
		if verbose {
			log.Printf("Worker %d: received %v", worker, base)
		}
		cpf := base[:]
		cpf = append(cpf, verificationNumber(cpf))
		cpf = append(cpf, verificationNumber(cpf))
		var cpfArr [11]int
		for i, v := range cpf {
			cpfArr[i] = v
		}
		if verbose {
			log.Printf("Worker %d: generated CPF %v", worker, cpfArr)
		}
		if output == "" {
			write(os.Stdout, format, cpfArr)
		} else {
			file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal(err)
			}
			write(file, format, cpfArr)
		}
	}
}

// write results to file or stdout
func write(f io.Writer, format int, cpfArr [11]int) {
	output := bufio.NewWriter(f)
	resConv := make([]string, 11)
	for i, d := range cpfArr {
		resConv[i] = strconv.Itoa(d)
	}
	var cpf string
	switch format {
	case 1:
		cpf = strings.Join(resConv, "")
	case 2:
		cpf = strings.Join(resConv[:3], "") + "." +
			strings.Join(resConv[3:6], "") + "." + strings.Join(resConv[6:9], "") +
			"-" + strings.Join(resConv[9:11], "")
	case 3:
		cpf = strings.Join(resConv[:9], "") + "-" +
			strings.Join(resConv[9:11], "")
	}
	fmt.Fprintf(output, "%s\n", cpf)
	output.Flush() // Flush the buffer to ensure all data is written
}

// calculate verification number
func verificationNumber(base []int) int {
	i := 0
	var sum int
	for mult := 1 + len(base); mult > 1; mult-- {
		sum += mult * base[i]
		i++
	}
	rem := (sum * 10) % 11
	if rem == 10 {
		return 0
	} else {
		return rem
	}
}
