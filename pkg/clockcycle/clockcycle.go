package clockcycle

import (
	"time"
)

//ClockCycle sents ticks to a channel
type ClockCycle struct {
	StartTime        time.Time
	IntervalDuration time.Duration
	Frequency        int
	Update           chan ClockTime
	Interval         int
	Cycle            int
}

//ClockTime is a struct that stores the time of a cycle
type ClockTime struct {
	Interval int
	Cycle    int
}

//Start is a function that will begin the clock cycle and return updates to the Update channel
func (clockcycle *ClockCycle) Start() {
	for {
		startToNow := time.Since(clockcycle.StartTime)
		intervalsPassed := int(startToNow.Nanoseconds() / clockcycle.IntervalDuration.Nanoseconds())
		nextInterval := clockcycle.StartTime.Add(time.Duration(int64(intervalsPassed+1) * clockcycle.IntervalDuration.Nanoseconds()))
		waitTime := time.Until(nextInterval)
		time.Sleep(waitTime)

		interval := intervalsPassed % clockcycle.Frequency
		cycle := int(intervalsPassed / clockcycle.Frequency)

		clockcycle.Interval = interval
		clockcycle.Cycle = cycle

		clockcycle.Update <- ClockTime{interval, cycle}
	}
}
