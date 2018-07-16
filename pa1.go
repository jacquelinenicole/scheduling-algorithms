/*
Jacqueline van der Meulen | 07/15/2018 | COP 4600
Reads in a file containing the following data: # processes | run length | scheduling algorithm | processes' names | processes' arrival time | processes' burst time
Scheduling algorithms supported: first-come first-served | preemptive shortest job first | round-robin
Simulates chosen scheduling algorithm and prints wait time and turnaround time for each process
*/
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
}

func main() {
    if len(os.Args) > 3 {
    	fmt.Println("Invalid argument list. Correct usage: \ngo run schedulingAlgorithms.go [file name]")
    	os.Exit(-1)
    }

    fileName := os.Args[1]
    runTime, algorithm, quantum, processes := parse(fileName)

    // output file
    o, _ := os.Create(os.Args[2])

    defer o.Close()

    if (algorithm == "fcfs") {
    	fcfs(o, runTime, processes)
    } else if (algorithm == "sjf") {
    	sjf(o, runTime, processes)
    } else if (algorithm == "rr") {
    	rr(o, runTime, quantum, processes)
    }  else {
    	fmt.Println("Invalid algorithm name give. Accepted algorithms: \nfcfs, sjf, rr")
    	os.Exit(-1)
    }

    o.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// gets all needed data from file; initializes processes array of struct Process
func parse(fileName string) (int, string, int, []Process) {
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

    	processes = append(processes, Process{name, arrival, burst, 0, 0, 0})
    }

    file.Close()

	return runTime, algorithm, quantum, processes
}

// finds keyword given then returns the word/value coming after it
func getValue(s *bufio.Scanner, word string) string {
	for string(s.Bytes()) != word {
		s.Scan()
	}
	s.Scan()

	return string(s.Bytes())
}


/***** First-Come First-Served *****/
func fcfs(o *os.File, runTime int, processes []Process) {
	fmt.Fprintf(o, "%3d processes\nUsing First-Come First-Served\n", len(processes))

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
		checkArrival(o, time, processes)
		numProcFinished = checkFinished(o, time, numProcFinished, numProcFinished, processes)

		if numProcFinished == len(processes) {
			time = finishIdle(o, time, runTime, numProcFinished, processes)
		} else {
			checkSelectedFCFS(o, time, numProcFinished, processes)
		}
	}

	// bubble sort alphabetically
	for i := 0 ; i < len(processes) ; i++ {
		for j := 0 ; j < len(processes) - i - 1 ; j++ {
			if processes[j].name > processes[j+1].name {
				processes[j], processes[j+1] = processes[j+1], processes[j]
			}
		}
	}

	printTimes(o, runTime, processes)
}

func checkSelectedFCFS(o *os.File, time int, numProcFinished int, processes []Process) {
	if processes[numProcFinished].selected == time {
		fmt.Fprintf(o, "Time %3d : %s selected (burst %3d)\n", time, processes[numProcFinished].name, processes[numProcFinished].burst)
	}

	if processes[numProcFinished].selected <= time {
		processes[numProcFinished].timeBursted++
	} else {
		fmt.Fprintf(o, "Time %3d : Idle\n", time)
	}
}


/***** Preemptive Shortest Job First *****/
func sjf(o *os.File, runTime int, processes []Process) {
	fmt.Fprintf(o, "%3d processes\nUsing preemptive Shortest Job First\n", len(processes))

	// bubble sort by arrival time
	for i := 0 ; i < len(processes) ; i++ {
		for j := 0 ; j < len(processes) - i - 1 ; j++ {
			if processes[j].arrival > processes[j+1].arrival {
				processes[j], processes[j+1] = processes[j+1], processes[j]
			}
		}
	}

	numProcFinished := 0
	mostRecent := -1
	for time := 0 ; time < runTime ; time++ {
		checkArrival(o, time, processes)
		
		numProcFinished = checkFinished(o, time, mostRecent, numProcFinished, processes)
		
		if numProcFinished == len(processes) {
			time = finishIdle(o, time, runTime, numProcFinished, processes)
		}	else {
			mostRecent = checkSelectedSJF(o, runTime, time, processes, mostRecent)
		}
	}

	// bubble sort alphabetically
	for i := 0 ; i < len(processes) ; i++ {
		for j := 0 ; j < len(processes) - i - 1 ; j++ {
			if processes[j].name > processes[j+1].name {
				processes[j], processes[j+1] = processes[j+1], processes[j]
			}
		}
	}

	printTimes(o, runTime, processes)
}

// returns process index if running, -1 if idle
func checkSelectedSJF(o *os.File, runTime int, time int, processes []Process, mostRecent int) int {
	i := shortestProcess(runTime, time, processes)

	// idle
	if i == -1 {
		fmt.Fprintf(o, "Time %3d : Idle\n", time)
		return -1
	}

	// grab inital selection time
	if processes[i].timeBursted == 0 {
		processes[i].selected = time
	}

	// only print it's been selected the first of each mini-burst
	if mostRecent != i {
		
		fmt.Fprintf(o, "Time %3d : %s selected (burst %3d)\n", time, processes[i].name, processes[i].burst - processes[i].timeBursted)
	}
	processes[i].timeBursted++
	return i
}

func shortestProcess(runTime int, time int, processes []Process) int {
	shortestProcessIndex := -1
	shortestProcess := runTime + 1

	for i := 0 ; i < len(processes) ; i++ {
		
		// check if process has arrived
		if processes[i].arrival <= time {

			// check if process still needs to be run
			if processes[i].timeBursted < processes[i].burst {
				if processes[i].burst - processes[i].timeBursted < shortestProcess {
					shortestProcessIndex = i
					shortestProcess = processes[i].burst - processes[i].timeBursted
				}
			}
		}
	}

	return shortestProcessIndex
}


/***** Round Robin *****/
func rr(o *os.File, runTime int, quantum int, processes []Process) {
	fmt.Fprintf(o, "%3d processes\nUsing Round-Robin\nQuantum %3d\n\n", len(processes), quantum)

	queue := []int{}
	numProcFinished := 0
	mostRecent := -1

	for time := 0 ; time < runTime ; time++ {
		queue = checkArrivalRR(o, time, processes, queue)

		numProcFinished = checkFinished(o, time, mostRecent, numProcFinished, processes)

		if numProcFinished == len(processes) {
			time = finishIdle(o, time, runTime, numProcFinished, processes)
		}	else {
			mostRecent, queue = checkSelectedRR(o, runTime, time, processes, mostRecent, quantum, queue)	
		}
	}

	printTimes(o, runTime, processes)
}

// returns process index if running, -1 if idle
func checkSelectedRR(o *os.File, runTime int, time int, processes []Process, mostRecent int, quantum int, queue []int) (int, []int) {
	
	// process last bursted still needs to run
	if mostRecent != -1 && processes[mostRecent].timeBursted < processes[mostRecent].burst {
		
		// if quantum incomplete, continue bursting
		if processes[mostRecent].timeBursted % quantum != 0 {
			processes[mostRecent].timeBursted++
			return mostRecent, queue
		} else {
			// quantum is completed and still needs to burst so add back into queue
			queue = append(queue, mostRecent)
		}
	}

	// idle
	if len(queue) == 0 {
		fmt.Fprintf(o, "Time %3d : Idle\n", time)
		return -1, queue	
	}

	// pop off queue
	i := queue[0]
	if len(queue) == 1 {
		queue = queue[:0]
	} else {
		queue = queue[1:]
	}

	mostRecent = i
	if processes[i].timeBursted < processes[i].burst {
		fmt.Fprintf(o, "Time %3d : %s selected (burst %3d)\n", time, processes[i].name, processes[i].burst - processes[i].timeBursted)
		processes[i].timeBursted++
	}

	return mostRecent, queue
}

// if a process has arrived, add to queue
func checkArrivalRR(o *os.File, time int, processes []Process, queue []int) []int {
	for i := 0 ; i < len(processes) ; i++ {
		if processes[i].arrival == time {
			queue = append(queue, i)
			fmt.Fprintf(o, "Time %3d : %s arrived\n", time, processes[i].name)
		}
	}

	return queue
}

// increase numProcFinished if most recently used process has finished
func checkFinished(o *os.File, time int, curr int, numProcFinished int, processes []Process) int {
	if curr != -1 && processes[curr].timeBursted == processes[curr].burst {
		fmt.Fprintf(o, "Time %3d : %s finished\n", time, processes[curr].name)
		processes[curr].finished = time
		numProcFinished++
	}

	return numProcFinished
}

// go through each process and see if it has arrived at current time slot
func checkArrival(o *os.File, time int, processes []Process) {
	for i := 0 ; i < len(processes) ; i++ {
		if processes[i].arrival == time {
			fmt.Fprintf(o, "Time %3d : %s arrived\n", time, processes[i].name)
		}
	}
}

// all processes have finished so print "idle" for remaining time slots
func finishIdle(o *os.File, time int, runTime int, numProcFinished int, processes []Process) int {	
	for time < runTime {
		fmt.Fprintf(o, "Time %3d : Idle\n", time)
		time++
	}

	return time
}

// print runtime, wait time, and turnaround time
func printTimes(o *os.File, runTime int, processes []Process) {
	fmt.Fprintf(o, "Finished at time %3d\n\n", runTime)

    for _, process := range processes {
		fmt.Fprintf(o, "%s wait %3d turnaround %3d\n", process.name, process.finished - process.arrival - process.burst, process.finished - process.arrival)
	}
}