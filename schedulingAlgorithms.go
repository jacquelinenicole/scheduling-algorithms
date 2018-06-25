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
	timeBursted int
	selected int
	finished int
	previouslySelected bool
}

func main() {
    if len(os.Args) != 2 {
    	fmt.Println("Invalid argument list. Correct usage: \ngo run schedulingAlgorithms.go [file name]")
    	os.Exit(-1)
    }

    fileName := os.Args[1]
    runTime, algorithm, quantum, processes := getInfo(fileName)

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

    	processes = append(processes, Process{name, arrival, burst, 0, 0, 0, false})

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
		if (checkFinished(time, numProcFinished, processes)) {
			numProcFinished++
		}

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

func checkFinished(time int, numProcFinished int, processes []Process) bool {
	if processes[numProcFinished].selected + processes[numProcFinished].burst == time {
		fmt.Printf("TIME %3d : %3s finished\n", time, processes[numProcFinished].name)
		processes[numProcFinished].finished = time

		return true
	}

	return false
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
	fmt.Printf("%3d processes\nUsing preemptive Shortest Job First\n", len(processes))

	// bubble sort by shortest job
	for i := 0 ; i < len(processes) ; i++ {
		for j := 0 ; j < len(processes) - i - 1 ; j++ {
			if processes[j].burst > processes[j+1].burst {
				processes[j], processes[j+1] = processes[j+1], processes[j]
			}
		}
	}

	numProcFinished := 0
	mostRecent := -1
	for time := 0 ; time < runTime ; time++ {
		checkArrival(time, numProcFinished, processes)
		
		if mostRecent != -1 {
			numProcFinished = checkFinishedSJF(time, numProcFinished, mostRecent, processes)
		
		}
		
		if numProcFinished != len(processes) {
			mostRecent = checkSelectedSJF(runTime, time, processes, mostRecent)
		}

		time = checkIdleSJF(time, runTime, numProcFinished, processes)
	}



	fmt.Printf("Finished at time %3d\n\n", runTime)

	//TODO: sort alphabetically
    // calculate wait and turnaround time
    for _, process := range processes {
		fmt.Printf("%s wait %3d turnaround %3d\n", process.name, process.selected - process.arrival, process.finished - process.arrival)
		fmt.Println(process)
	}

}

func checkSelectedSJF(runTime int, time int, processes []Process, mostRecent int) int{
	for i := 0 ; i < len(processes) ; i++ {

		// check if process has arrived
		if processes[i].arrival <= time {

			// check if process still needs to be run
			if processes[i].timeBursted < processes[i].burst {

				// grab inital selection time
				if !processes[i].previouslySelected {
					processes[i].previouslySelected = true
					processes[i].selected = time
				}

				// only print it's been selected the first of each mini-burst
				if mostRecent != i {
					mostRecent = i

					burst := runTime + 1

					// if a shorter job will arrive before the current process is done, burst time is reduced
					for j := 0 ; j < i ; j++ {
						if processes[j].timeBursted < processes[j].burst {
							if processes[j].arrival < processes[i].burst - processes[i].timeBursted + time {
								if processes[j].arrival < burst {
									burst = processes[j].arrival - time
								}
							}
						}
					}

					// current burst will not be interruped and process will finish
					if burst == runTime + 1 {
						burst = processes[mostRecent].burst - processes[mostRecent].timeBursted
					}

					fmt.Printf("TIME %3d : %3s selected (burst %3d)\n", time, processes[mostRecent].name, burst)
				}
				
				processes[i].timeBursted++
				return mostRecent
			}				
		}
	}
	return -1
}


func checkIdleSJF(time int, runTime int, numProcFinished int, processes []Process) int {
			
	if numProcFinished == len(processes) {
		for time < runTime {
			fmt.Printf("TIME %3d : Idle\n", time)
			time++
		}

		return time
	} 

	idle := true

	for i := 0 ; i < len(processes) ; i++ {
		if processes[i].arrival <= time {
			if processes[i].timeBursted <= processes[i].burst {
				idle = false
				break
			}
		}
	}
	if idle {
		fmt.Printf("TIME %3d : Idle\n", time)
	}

	return time
}

func checkFinishedSJF(time int, numProcFinished int, mostRecent int, processes []Process) int{
	if processes[mostRecent].timeBursted == processes[mostRecent].burst {
		fmt.Printf("TIME %3d : %3s finished\n", time, processes[mostRecent].name)
		processes[mostRecent].finished = time
		numProcFinished++
	}

	return numProcFinished
}

func rr(runTime int, quantum int, processes []Process) {
	fmt.Printf("%3d processes\nUsing Round Robin\n Quantum %3d\n", len(processes), quantum)
}
