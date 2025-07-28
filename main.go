package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

// Windows power state flags
const (
	ES_CONTINUOUS       = 0x80000000
	ES_SYSTEM_REQUIRED  = 0x00000001
	ES_DISPLAY_REQUIRED = 0x00000002
	VK_F15              = 0x7E
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	setThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
)

func preventSleep() {
	// Prevent sleep and keep display on
	setThreadExecutionState.Call(
		ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED,
	)
}

func allowSleep() {
	// Restore default power state
	setThreadExecutionState.Call(ES_CONTINUOUS)
}

func moveMouse() {
	// Get initial mouse position
	initX, initY := robotgo.GetMousePos()

	// Random movement range (pixels)
	const moveRange = 50

	// Generate random movement
	deltaX := rand.Intn(moveRange*2) - moveRange
	deltaY := rand.Intn(moveRange*2) - moveRange

	// Move mouse to new position
	robotgo.MoveSmooth(initX+deltaX, initY+deltaY, 1.0, 2.0)

	// Small delay to make movement visible
	time.Sleep(100 * time.Millisecond)

	// Return to initial position
	robotgo.MoveSmooth(initX, initY, 1.0, 2.0)
}

func main() {
	move := flag.Bool("move", false, "a bool")
	interval := flag.Int("interval", 10, "an int")

	flag.Parse()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Prevent sleep
	preventSleep()

	// Make sure we restore default state when program exits
	defer allowSleep()

	fmt.Printf("Avake every %d seconds. Press Ctrl+C to stop.\n", *interval)

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	// Infinite loop until program is interrupted
	for {
		select {
		case <-signalChan:
			fmt.Println("\nRestoring normal power settings...")
			return
		case t := <-ticker.C:
			if *move {
				moveMouse()
			}
			robotgo.KeyPress("F15")
			fmt.Printf("Working... %s\n", t.Format("15:04:05"))
		}
	}
}
