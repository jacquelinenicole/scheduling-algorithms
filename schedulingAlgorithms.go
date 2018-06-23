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
    s, numProcesses, runTime, algorithm, quantum := getInfo(fileName)



    fmt.Println(numProcesses)
    fmt.Println(runTime)
    fmt.Println(algorithm)
    fmt.Println(quantum)

    if (algorithm == "fcfs") {
    	fcfs(s, runTime)
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


func getInfo(fileName string) (*bufio.Scanner, int, int, string, int) {

   
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

    //processesNames := make([]int, numProcesses)
    //arrivals := make([]int, numProcesses)
    //bursts := make([]int, numProcesses)

	return s, numProcesses, runTime, algorithm, quantum
}

func getValue(s *bufio.Scanner, word string) string {
	for string(s.Bytes()) != word {
		s.Scan()
	}

	s.Scan()

	return string(s.Bytes())

}

func fcfs(s *bufio.Scanner, runTime int) {
	fmt.Println("fcfs")
}