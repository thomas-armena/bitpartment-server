package clockcycle

import(
	"time"
)

//ClockCycle sents ticks to a channel 
type ClockCycle struct {
	StartTime time.Time
	Interval time.Duration
	Frequency int
	Update chan ClockTime
}

//ClockTime is a struct that stores the time of a cycle
type ClockTime struct {
	Interval	int
	Cycle		int
}

//Start is a function that will begin the clock cycle and return updates to the Update channel
func (clockcycle *ClockCycle) Start() {
	for {
		startToNow := time.Since(clockcycle.StartTime)
		intervalsPassed := int(startToNow.Nanoseconds() / clockcycle.Interval.Nanoseconds())
		nextInterval := clockcycle.StartTime.Add(time.Duration(int64(intervalsPassed+1) * clockcycle.Interval.Nanoseconds()))
		waitTime := time.Until(nextInterval)
		time.Sleep(waitTime)

		interval := intervalsPassed % clockcycle.Frequency
		cycle := int(intervalsPassed / clockcycle.Frequency)
		clockcycle.Update <- ClockTime{ interval, cycle }
	}
}

/*
func main() {
	fmt.Println("clock cycle test")
	update := make(chan int)
	location, _ := time.LoadLocation("EST")
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 15, 3, 0, location)
	fmt.Println("startTime:", startTime)
	fmt.Println("now:      ",time.Now())
	clockCycle := ClockCycle{ StartTime: startTime, Interval: time.Duration(10*time.Second), Frequency: 12, Update: update }
	fmt.Println(time.Now())
	go clockCycle.Start()
	for {
		fmt.Println(<-clockCycle.Update)
	}
}
*/


