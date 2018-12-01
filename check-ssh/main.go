package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

var (
	inittime = time.Now()
	password = flag.String("password", "password123", "Password to use")
	ip       = flag.String("ip", "123.123.123.123", "IP address connect to")
	port     = flag.Int("port", 22, "Port of server")
	user     = flag.String("user", "root", "User to login as")
	attempts = flag.Int("attempts", 3, "Amount of times to attempt login")
	timer    = flag.Duration("timer", 300*time.Millisecond, "Timeout between attempts")
)

type resp struct {
	Error error
	mu    sync.Mutex
}

func sshdialer() *resp {
	exitcode := &resp{}
	config := &ssh.ClientConfig{

		User:            *user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(*password)},
		Timeout:         *timer,
	}

	_, err := ssh.Dial("tcp", *ip+":"+strconv.Itoa(*port), config)
	if err != nil {
		fmt.Printf("\nFailed connection")
	} else {
		end := time.Now()
		d := end.Sub(inittime)
		duration := d.Seconds()
		fmt.Fprintf(color.Output, "\n%s", color.YellowString("###########################"))
		fmt.Fprintf(color.Output, "\n%s", color.GreenString("Successful login"))
		fmt.Fprintf(color.Output, "\n%s", color.YellowString("###########################"))
		fmt.Printf("\nCompleted in %v seconds\n", strconv.FormatFloat(duration, 'g', -1, 64))
	}

	exitcode.Error = err
	return exitcode
}

func main() {
	flag.Parse()

	for attempt := *attempts; attempt != 0; attempt-- {
		go func() {
			resp := sshdialer()
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