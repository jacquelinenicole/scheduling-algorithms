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
	selected int
}

func main() {
    if len(os.Args) != 2 {
    	fmt.Println("Invalid argument list. Correct usage: \ngo run schedulingAlgorithms.go [file name]")
    	os.Exit(-1)
    }

    fileName := os.Args[1]
    runTime, algorithm, quantum, processes := getInfo(fileName)
/*
    fmt.Println(runTime)
    fmt.Println(algorithm)
    fmt.Println(quantum)

    for _,process := range processes {
    	fmt.Printf("Name: %3s\nArrival: %3d\nBurst: %3d\n\n", process.name, process.arrival, process.burst)
    }
*/
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
	
	numProcesses, _ := strconv.Atoi(getValue(s, "processcount"))

	runTime, _ := strconv.Atoi(getValue(s, "runfor"))

	algorithm := getValue(s, "use")

	quantum := -1

	if algorithm == "rr" {
		quantum, _ = strconv.Atoi(getValue(s, "quantum"))
	}

    var processes []Process

    for i := 0 ; i < numProcesses ; i++ {
    	name := getValue(s, "name")
    	arrival, _ := strconv.Atoi(getValue(s, "arrival"))
    	burst, _ := strconv.Atoi(getValue(s, "burst"))

    	processes = append(processes, Process{name, arrival, burst, 0})

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
	fmt.Printf("%3d processes\nUsing First-Come First-Served\n", len(processes))

	numProcFinished := 0

	// bubble sort by arrival time
	for i := 0 ; i < len(processes) ; i++ {
		for j := 0 ; j < len(processes) - i - 1 ; j++ {
			if processes[j].arrival > processes[j+1].arrival {
				processes[j], processes[j+1] = processes[j+1], processes[j]
			}
		}
	}

	// find selected times
	for i := 0 ; i < len(processes) - 1  ; i++ {
		if processes[i].selected + processes[i].burst <= processes[i+1].arrival {
			processes[i+1].selected = processes[i+1].arrival
		} else {
			processes[i+1].selected = processes[i].selected + processes[i].burst
		}
	}

	for time := 0 ; time < runTime ; time++ {
		checkArrival(time, numProcFinished, processes)
		numProcFinished = checkFinished(time, numProcFinished, processes)
		if numProcFinished != len(processes) {
			checkSelected(time, numProcFinished, processes)
		}

		time = checkIdle(time, runTime, numProcFinished, processes)

	}

	fmt.Printf("Finished at time %3d\n\n", runTime)

	//TODO: sort alphabetically

    // calculate wait and turnaround time

    for _, process := range processes {
    	wait := process.selected - process.arrival
		fmt.Printf("%s wait %3d turnaround %3d\n", process.name, wait, process.burst + wait)
	}

}

func checkArrival(time int, numProcFinished int, processes []Process) {
	for i := numProcFinished ; i < len(processes) ; i++ {
		if processes[i].arrival == time {
			fmt.Printf("TIME %3d : %3s arrived \n", time, processes[i].name)

		}
	}
}

func checkSelected(time int, numProcFinished int, processes []Process) {
	if processes[numProcFinished].selected == time {
		fmt.Printf("TIME %3d : %3s selected (burst %3d)\n", time, processes[numProcFinished].name, processes[numProcFinished].burst)
	}
}

func checkFinished(time int, numProcFinished int, processes []Process) int {
	if processes[numProcFinished].selected + processes[numProcFinished].burst == time {
		fmt.Printf("TIME %3d : %3s finished\n", time, processes[numProcFinished].name)

		numProcFinished = numProcFinished + 1
	}

	return numProcFinished
}

func checkIdle(time int, runTime int, numProcFinished int, processes []Process) int {
			
	if numProcFinished == len(processes) {
		for time < runTime {
			fmt.Printf("TIME %3d : Idle\n", time)
			time++
		}
	} else if processes[numProcFinished].arrival > time {
		fmt.Printf("TIME %3d : Idle\n", time)
	}

	return time
}

func sjf(runTime int, processes []Process) {
	fmt.Println("sjf")
}

func rr(runTime int, quantum int, processes []Process) {
	fmt.Println("rr")
}
