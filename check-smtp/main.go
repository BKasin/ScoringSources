package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

var (
	inittime = time.Now()
	user     = flag.String("user", "root", "User to login as")
	password = flag.String("password", "password123", "Password to use")
	ip       = flag.String("ip", "123.123.123.123", "IP address connect to")
	port     = flag.Int("port", 25, "Port of server")
	attempts = flag.Int("attempts", 3, "Amount of times to attempt login")
	timer    = flag.Duration("timer", 300*time.Millisecond, "Timeout between attempts")
)

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func smtpdialer() (err error) {
	smtpServer := SmtpServer{host: *ip, port: strconv.Itoa(*port)}
	auth := smtp.PlainAuth("", *user, *password, smtpServer.host)

	client, err := smtp.Dial(smtpServer.host)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Auth(auth); err == nil {
		end := time.Now()
		d := end.Sub(inittime)
		duration := d.Seconds()
		fmt.Fprintf(color.Output, "\n%s", color.YellowString("###########################"))
		fmt.Fprintf(color.Output, "\n%s", color.GreenString("Successful connection"))
		fmt.Fprintf(color.Output, "\n%s", color.YellowString("###########################"))
		fmt.Printf("\nCompleted in %v seconds\n", strconv.FormatFloat(duration, 'g', -1, 64))
		defer client.Quit()
	} else {
		fmt.Println(err)
		defer client.Quit()
	}

	return err
}

func main() {
	flag.Parse()

	for attempt := *attempts; attempt != 0; attempt-- {
		go func() {
			resp := smtpdialer()
			if resp == nil {
				os.Exit(0)
			}
		}()

		fmt.Println("Attempt: ", attempt)
		time.Sleep(*timer)
	}
}
