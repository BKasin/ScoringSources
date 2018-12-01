package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	pop3 "github.com/simia-tech/go-pop3"
)

var (
	inittime = time.Now()
	password = flag.String("password", "password123", "Password to use")
	ip       = flag.String("ip", "123.123.123.123", "IP address connect to")
	user     = flag.String("user", "root", "User to login as")
	attempts = flag.Int("attempts", 3, "Amount of times to attempt login")
	timer    = flag.Duration("timer", 300*time.Millisecond, "Timeout between attempts")
)

type resp struct {
	Error error
	mu    sync.Mutex
}

func pop3dialer() *resp {
	exitcode := &resp{}
	pop, err := pop3.Dial(*ip, pop3.UseTimeout(*timer))

	if err = pop.Auth(*user, *password); err == nil {
		end := time.Now()
		d := end.Sub(inittime)
		duration := d.Seconds()
		fmt.Fprintf(color.Output, "\n%s", color.YellowString("###########################"))
		fmt.Fprintf(color.Output, "\n%s", color.GreenString("Successful connection"))
		fmt.Fprintf(color.Output, "\n%s", color.YellowString("###########################"))
		fmt.Printf("\nCompleted in %v seconds\n", strconv.FormatFloat(duration, 'g', -1, 64))
		defer pop.Quit()
	} else {
		fmt.Println(err)
		defer pop.Quit()
	}

	exitcode.Error = err
	return exitcode
}

func main() {
	flag.Parse()

	for attempt := *attempts; attempt != 0; attempt-- {
		go func() {
			resp := pop3dialer()
			resp.mu.Lock()
			if resp.Error == nil {
				resp.mu.Unlock()
				os.Exit(0)
			}
		}()

		fmt.Println("Attempt: ", attempt)
		time.Sleep(*timer)
	}
}
