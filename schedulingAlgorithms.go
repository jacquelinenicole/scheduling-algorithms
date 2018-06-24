package main
import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

type Process struct {
	name string
	arrival int
	burst int
}

func main() {
    if len(os.Args) != 2 {
    	fmt.Println("Invalid argument list. Correct usage: \ngo run schedulingAlgorithms.go [file name]")
    	os.Exit(-1)
    }

    fileName := os.Args[1]
    runTime, algorithm, quantum, processes := getInfo(fileName)

    fmt.Println(runTime)
    fmt.Println(algorithm)
    fmt.Println(quantum)

    for _,process := range processes {
    	fmt.Printf("Name: %3s\nArrival: %3d\nBurst: %3d\n\n", process.name, process.arrival, process.burst)
    }

    if (algorithm == "fcfs") {
    	fcfs(runTime, processes)
    } else if (algorithm == "sjf") {
    	sjf(runTime, processes)
    } else if (algorithm == "rr") {
    	rr(runTime, quantum, processes)
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


func getInfo(fileName string) (int, string, int, []Process) {
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

    var processes []Process

    for i := 0 ; i < numProcesses ; i++ {
    	name := getValue(s, "name")
    	arrival, err := strconv.Atoi(getValue(s, "arrival"))
    	check(err)
    	burst, err := strconv.Atoi(getValue(s, "burst"))
    	check(err)

    	processes = append(processes, Process{name, arrival, burst})

    }

	return runTime, algorithm, quantum, processes
}

func getValue(s *bufio.Scanner, word string) string {
	for string(s.Bytes()) != word {
		s.Scan()
	}

	s.Scan()

	return string(s.Bytes())
}

func fcfs(runTime int, processes []Process) {
	fmt.Println("fcfs")

	fmt.Printf("%3d processes\n", len(processes))

	for time := 0 ; time < runTime ; time++ {
		checkArrival(time, processes)

		fmt.Printf("TIME %3d : \n", time)
	}

	fmt.Printf("Finished at time %3d", runTime)

}

func checkArrival(time int, processes []Process) {
	for _,process := range processes {
		if process.arrival == time {
			fmt.Printf("TIME %3d : %3s arrived \n", time, process.name)
		}
	}
}

func sjf(runTime int, processes []Process) {
	fmt.Println("sjf")
}

func rr(runTime int, quantum int, processes []Process) {
	fmt.Println("rr")
}
