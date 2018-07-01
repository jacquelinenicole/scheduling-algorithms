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
    if len(os.Args) > 3 { {
    	fmt.Println("Invalid argument list. Correct usage: \ngo run schedulingAlgorithms.go [file name]")
    	os.Exit(-1)
    }

    fileName := os.Args[1]
    runTime, algorithm, quantum, processes := parse(fileName)

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
		checkArrival(time, processes)
		numProcFinished = checkFinished(time, numProcFinished, numProcFinished, processes)

		if numProcFinished == len(processes) {
			time = finishIdle(time, runTime, numProcFinished, processes)
		} else {
			checkSelectedFCFS(time, numProcFinished, processes)
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

	printTimes(runTime, processes)
}

func checkSelectedFCFS(time int, numProcFinished int, processes []Process) {
	if processes[numProcFinished].selected == time {
		fmt.Printf("TIME %3d : %3s selected (burst %3d)\n", time, processes[numProcFinished].name, processes[numProcFinished].burst)
	}

	if processes[numProcFinished].selected <= time {
		processes[numProcFinished].timeBursted++
	} else {
		fmt.Printf("TIME %3d : Idle\n", time)
	}
}

func sjf(runTime int, processes []Process) {
	fmt.Printf("%3d processes\nUsing preemptive Shortest Job First\n", len(processes))

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
		checkArrival(time, processes)
		
		numProcFinished = checkFinished(time, mostRecent, numProcFinished, processes)
		
		
		if numProcFinished == len(processes) {
			time = finishIdle(time, runTime, numProcFinished, processes)
		}	else {
			mostRecent = checkSelectedSJF(runTime, time, processes, mostRecent)
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

	printTimes(runTime, processes)
}

// returns process index if running, -1 if idle
func checkSelectedSJF(runTime int, time int, processes []Process, mostRecent int) int{
	i := shortestProcess(runTime, time, processes)

	if i == -1 {
		fmt.Printf("TIME %3d : Idle\n", time)
		return -1
	}

	// grab inital selection time
	if processes[i].timeBursted == 0 {
		processes[i].selected = time
	}

	// only print it's been selected the first of each mini-burst
	if mostRecent != i {
		// burst := calcBurst(processes[i].burst - processes[i].timeBursted, processes, time, i)
		
		fmt.Printf("TIME %3d : %3s selected (burst %3d)\n", time, processes[i].name, processes[i].burst - processes[i].timeBursted)
	}
	processes[i].timeBursted++
	return i
}

// if a shorter job will arrive before the current process is done, burst time is reduced
func calcBurst(burst int, processes []Process, time int, i int) int{
	for j := 0 ; j < len(processes) ; j++ {
		if j == i {
			continue
		}

		// another process still needs to run
		if processes[j].timeBursted < processes[j].burst {

			// process arrives before current job is finished
			if processes[j].arrival < processes[i].burst - processes[i].timeBursted + time {

				// process has a shorter burst time than current process's burst time at its arrival
				if processes[j].burst - processes[j].timeBursted < processes[i].burst - processes[i].timeBursted - (processes[j].arrival - time) {

					// also shorter burst than any other process meeting prior critera
					if processes[j].burst - processes[j].timeBursted < burst {
						return processes[j].arrival - time
					}
				}
			}
		}
	}

	return burst
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

func rr(runTime int, quantum int, processes []Process) {
	fmt.Printf("%3d processes\nUsing Round Robin\nQuantum %3d\n", len(processes), quantum)

	numProcFinished := 0
	queue := []int{}
	mostRecent := -1

	for time := 0 ; time < runTime ; time++ {
		queue = checkArrivalRR(time, processes, queue)

		numProcFinished = checkFinished(time, mostRecent, numProcFinished, processes)

		if numProcFinished == len(processes) {
			time = finishIdle(time, runTime, numProcFinished, processes)
		}	else {
			mostRecent, queue = checkSelectedRR(runTime, time, processes, mostRecent, quantum, queue)	
		}
	}

	printTimes(runTime, processes)
}

// returns process index if running, -1 if idle
func checkSelectedRR(runTime int, time int, processes []Process, mostRecent int, quantum int, queue []int) (int, []int) {
	// process still needs to run
	if mostRecent != -1 && processes[mostRecent].timeBursted < processes[mostRecent].burst {
		
		// check if last process finished quantum
		if processes[mostRecent].timeBursted % quantum != 0 {
			processes[mostRecent].timeBursted++

			return mostRecent, queue
		} else if processes[mostRecent].timeBursted < processes[mostRecent].burst {
			
			// quantum is completed and still needs to burst so add back into queue
			queue = append(queue, mostRecent)
		}
	}

	// idle
	if len(queue) == 0 {
		fmt.Printf("TIME %3d : Idle\n", time)
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
		fmt.Printf("TIME %3d : %3s selected (burst %3d)\n", time, processes[i].name, processes[i].burst - processes[i].timeBursted)
		processes[i].timeBursted++
	}

	return mostRecent, queue
}

func checkArrivalRR(time int, processes []Process, queue []int) []int{
	for i := 0 ; i < len(processes) ; i++ {
		if processes[i].arrival == time {
			queue = append(queue, i)
			fmt.Printf("TIME %3d : %3s arrived \n", time, processes[i].name)
		}
	}

	return queue
}

func checkFinished(time int, curr int, numProcFinished int, processes []Process) int {
	if curr != -1 && processes[curr].timeBursted == processes[curr].burst {
		fmt.Printf("TIME %3d : %3s finished\n", time, processes[curr].name)
		processes[curr].finished = time
		numProcFinished++
	}

	return numProcFinished
}

func checkArrival(time int, processes []Process) {
	for i := 0 ; i < len(processes) ; i++ {
		if processes[i].arrival == time {
			fmt.Printf("TIME %3d : %3s arrived \n", time, processes[i].name)
		}
	}
}

func finishIdle(time int, runTime int, numProcFinished int, processes []Process) int {	
	for time < runTime {
		
		time++
	}

	return time
}

// print runtime, wait time, and turnaround time
func printTimes(runTime int, processes []Process) {
	fmt.Printf("Finished at time %3d\n\n", runTime)

	//TODO: sort alphabetically by process name
    for _, process := range processes {
		fmt.Printf("%s wait %3d turnaround %3d\n", process.name, process.finished - process.arrival - process.burst, process.finished - process.arrival)
	}
}