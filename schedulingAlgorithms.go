package main
import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

func main() {
    if len(os.Args) != 2 {
    	fmt.Println("Invalid argument list. Correct usage: \ngo run schedulingAlgorithms.go [file name]")
    	os.Exit(-1)
    }

    fileName := os.Args[1]
    s, runTime, algorithm, quantum, processesNames, arrivals, bursts := getInfo(fileName)

    fmt.Println(runTime)
    fmt.Println(algorithm)
    fmt.Println(quantum)
    fmt.Println(processesNames)
	fmt.Println(arrivals)
	fmt.Println(bursts)

    if (algorithm == "fcfs") {
    	fcfs(s, runTime, processesNames, arrivals, bursts)
    } else if (algorithm == "sjf") {
    	sjf(s, runTime, processesNames, arrivals, bursts)
    } else if (algorithm == "rr") {
    	rr(s, runTime, quantum, processesNames, arrivals, bursts)
    }  else {
    	fmt.Println("Invalid algorithm name give. Accepted algorithms: \nfcfs, sjf, rr")
    	os.Exit(-1)
    }

	/*
	for s.Scan() {
		word := s.Bytes()
		fmt.Println(string(word))
	}
	*/
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}


func getInfo(fileName string) (*bufio.Scanner, int, string, int, []string, []int, []int) {
    file, err := os.Open(fileName)
    check(err)

    s := bufio.NewScanner(file)
	s.Split(bufio.ScanWords)
	
	numProcesses, err := strconv.Atoi(getValue(s, "processcount"))
	check(err)

	runTime, err := strconv.Atoi(getValue(s, "runfor"))
	check(err)

	algorithm := getValue(s, "use")

	quantum := -1

	if algorithm == "rr" {
		quantum, err := strconv.Atoi(getValue(s, "quantum"))
		check(err)
		fmt.Println(quantum)
	}

    processesNames := make([]string, numProcesses)
    arrivals := make([]int, numProcesses)
    bursts := make([]int, numProcesses)

    for i := 0 ; i < numProcesses ; i++ {
    	processesNames[i] = getValue(s, "name")
    	arrivals[i], err = strconv.Atoi(getValue(s, "arrival"))
    	check(err)
    	bursts[i], err = strconv.Atoi(getValue(s, "burst"))
    	check(err)
    }

	return s, runTime, algorithm, quantum, processesNames, arrivals, bursts
}

func getValue(s *bufio.Scanner, word string) string {
	for string(s.Bytes()) != word {
		s.Scan()
	}

	s.Scan()

	return string(s.Bytes())
}

func fcfs(s *bufio.Scanner, runTime int, processesNames []string, arrivals []int, bursts []int) {
	fmt.Println("fcfs")
}

func sjf(s *bufio.Scanner, runTime int, processesNames []string, arrivals []int, bursts []int) {
	fmt.Println("sjf")
}

func rr(s *bufio.Scanner, runTime int, quantum int, processesNames []string, arrivals []int, bursts []int) {
	fmt.Println("rr")
}