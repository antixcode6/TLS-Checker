package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/scorredoira/email"
)

var (
	dialer = &net.Dialer{Timeout: 5 * time.Second}
	file = flag.String("f", "", "read server names from `file`")
	alertFlag = 0
)

func check(server string, width int) {
	conn, err := tls.DialWithDialer(dialer, "tcp", server+":443", nil)
	if err != nil {
		alertFlag++
		go handleError(server, err.Error())
		return
	}
	defer conn.Close()
	valid := conn.VerifyHostname(server)

	for _, c := range conn.ConnectionState().PeerCertificates {
		if valid == nil {
			fmt.Printf("%*s | valid, expires on %s (%s)\n", width, server,
				c.NotAfter.Format("2006-01-02"), humanize.Time(c.NotAfter))
		} else {
			fmt.Printf("%*s | %v\n", width, server, valid)
		}
		return
	}
}

func handleError(server string, badcert string) {
	//open log file
	f, err := os.OpenFile("CertErrors.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != err {
		log.Fatalf("error opening file %v", err)
	}
	defer f.Close()

	//write cert errors to log file
	log.SetOutput(f)
	log.Println(server, badcert) //write cert error to log file
}

func main() {
	// parse command-line args
	flag.Parse()
	if flag.NArg() == 0 && len(*file) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: tlscheck [-f file] servername ...\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// collect list of server names
	names := getNames()

	// for cosmetics
	width := 0
	for _, name := range names {
		if len(name) > width {
			width = len(name)
		}
	}

	// actually check
	fmt.Printf("%*s | Certificate status\n%s-+-%s\n", width, "Server",
		strings.Repeat("-", width), strings.Repeat("-", 80-width-2))
	for _, name := range names {
		check(name, width)
	}

	if alertFlag == 0 {
		return
	} else {
		sendAlert()
	}
}

func getNames() (names []string) {

	// read names from the file
	if len(*file) > 0 {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) > 0 && line[0] != '#' {
				names = append(names, strings.Fields(line)[0])
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
		f.Close()
	}

	// add names specified on the command line
	names = append(names, flag.Args()...)
	return
}

func sendAlert() {
	msg := email.NewMessage("SSL Cert Issue Found", "View attachment for full log of errors")
	msg.From = mail.Address{Name: "From", Address: "from@example.com"}
	msg.To = []string{"to@example.com"}

	if err := msg.Attach("CertErrors.txt"); err != nil {
		log.Fatal(err)
	}

	auth := smtp.PlainAuth("", "from@example.com", "pwd", "smtp.zoho.com")
	if err := email.Send("smtp.zoho.com:587", auth, msg); err != nil {
		log.Fatal(err)
	}
}
