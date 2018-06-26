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
	lastStopped int
	wait int
}

func main() {
    if len(os.Args) != 2 {
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

    	processes = append(processes, Process{name, arrival, burst, 0, 0, 0, false, 0, 0})
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
		checkArrival(time, numProcFinished, processes)
		numProcFinished = checkFinished(time, numProcFinished, numProcFinished, processes)

		if numProcFinished == len(processes) {
			time = finishIdle(time, runTime, numProcFinished, processes)
		} else {
			checkSelectedFCFS(time, numProcFinished, processes)
		}

		
	}

	printTimes(runTime, processes)
}

func checkSelectedFCFS(time int, numProcFinished int, processes []Process) {
	if processes[numProcFinished].selected == time {
		fmt.Printf("TIME %3d : %3s selected (burst %3d)\n", time, processes[numProcFinished].name, processes[numProcFinished].burst)
		processes[numProcFinished].wait = processes[numProcFinished].selected - processes[numProcFinished].arrival
	}

	if processes[numProcFinished].selected <= time {
		processes[numProcFinished].timeBursted++
	} else {
		fmt.Printf("TIME %3d : Idle\n", time)
	}
}

func sjf(runTime int, processes []Process) {
	fmt.Printf("%3d processes\nUsing preemptive Shortest Job First\n", len(processes))

	numProcFinished := 0
	mostRecent := -1

	// bubble sort by arrival time
	for i := 0 ; i < len(processes) ; i++ {
		for j := 0 ; j < len(processes) - i - 1 ; j++ {
			if processes[j].arrival > processes[j+1].arrival {
				processes[j], processes[j+1] = processes[j+1], processes[j]
			}
		}
	}


	for time := 0 ; time < runTime ; time++ {
		checkArrival(time, numProcFinished, processes)
		
		if mostRecent != -1 {
			numProcFinished = checkFinished(time, mostRecent, numProcFinished, processes)
		}
		
		if numProcFinished == len(processes) {
			time = finishIdle(time, runTime, numProcFinished, processes)
		}	else {
			mostRecent = checkSelectedSJF(runTime, time, processes, mostRecent)
		}
	}

	printTimes(runTime, processes)
}

// returns process index if running, -1 if idle
func checkSelectedSJF(runTime int, time int, processes []Process, mostRecent int) int{
	i := shortestProcess(runTime, time, processes)

	if i == runTime + 1 {
		fmt.Printf("TIME %3d : Idle\n", time)
		return -1
	}

	// grab inital selection time
	if !processes[i].previouslySelected {
		processes[i].previouslySelected = true
		processes[i].selected = time
		processes[i].wait = processes[i].selected - processes[i].arrival
	}


	// only print it's been selected the first of each mini-burst
	if mostRecent != i {
		burst := getBurst(processes[i].burst - processes[i].timeBursted, processes, time, i)

		if processes[i].timeBursted != 0 {
			processes[i].wait += time - processes[i].lastStopped
		}

		fmt.Printf("TIME %3d : %3s selected (burst %3d)\n", time, processes[i].name, burst)
	}
	processes[i].lastStopped = time + 1
	processes[i].timeBursted++
	return i
}

func getBurst(burst int, processes []Process, time int, i int) int{
	// if a shorter job will arrive before the current process is done, burst time is reduced

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
	fmt.Printf("%3d processes\nUsing Round Robin\n Quantum %3d\n", len(processes), quantum)
}

func checkArrival(time int, numProcFinished int, processes []Process) {
	for i := numProcFinished ; i < len(processes) ; i++ {
		if processes[i].arrival == time {
			fmt.Printf("TIME %3d : %3s arrived \n", time, processes[i].name)

		}
	}
}

func checkFinished(time int, curr int, numProcFinished int, processes []Process) int {
	if processes[curr].timeBursted == processes[curr].burst {
		fmt.Printf("TIME %3d : %3s finished\n", time, processes[curr].name)
		processes[curr].finished = time
		numProcFinished++
	}

	return numProcFinished
}

func finishIdle(time int, runTime int, numProcFinished int, processes []Process) int {	
		for time < runTime {
			fmt.Printf("TIME %3d : Idle\n", time)
			time++
		}

	return time
}

// print runtime, wait time, and turnaround time
func printTimes(runTime int, processes []Process) {
	fmt.Printf("Finished at time %3d\n\n", runTime)

	//TODO: sort alphabetically by process name
    for _, process := range processes {
		fmt.Printf("%s wait %3d turnaround %3d\n", process.name, process.wait, process.finished - process.arrival)
	}
}