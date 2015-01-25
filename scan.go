package main

import (
	"bytes"
	"fmt"
	flag "github.com/ogier/pflag"
	"log"
	"net"
	"net/smtp"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func compareMapAlternate(X, Y []int) []int {
	counts := make(map[int]int)
	var total int
	for _, val := range X {
		counts[val] += 1
		total += 1
	}
	for _, val := range Y {
		if count := counts[val]; count > 0 {
			counts[val] -= 1
			total -= 1
		}
	}
	difference := make([]int, total)
	i := 0
	for val, count := range counts {
		for j := 0; j < count; j++ {
			difference[i] = val
			i++
		}
	}
	return difference
}

func main() {
	// Flags
	//// Host and it's expected ports
	host := flag.StringP("host", "h", "", "IP or resolvable hostname")
	port_range := flag.StringP("range", "p", "", "dash separated port range to scan, E.G. '1-65535'")
	expected_ports_str := flag.StringP("expected", "e", "",
		"space separated list of ports that are expected to be found unfiltered (unfiltered = open or closed)")

	//// SMTP stuff
	to_addresses := flag.StringP("to", "t", "", "space separated list of 'to' address(es)")
	from_address := flag.StringP("from", "f", "", "the 'from' email address")
	server_address := flag.StringP("server", "s", "", "the SMTP server address, E.G 'smtp.example.com:587'")
	user_address := flag.StringP("username", "u", "", "the SMTP username or email address")
	password := flag.StringP("password", "x", "", "the SMTP user password")

	//// Nmap --min-parallelism
	parallelism := flag.StringP("parallelism", "m", "", "the Nmap --min-parrallelism setting, E.G. '1024'")

	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NFlag() != 9 {
		flag.Usage()
		return
	}

	// Convert EXPECTED_PORTS string to int slice
	expected_ports_str_slice := strings.Fields(*expected_ports_str)
	expected_ports_ints_slice := []int{}
	for _, expected_port_str := range expected_ports_str_slice {
		expected_port_int, err := strconv.Atoi(expected_port_str)
		if err != nil {
			log.Fatalf("Error converting string '%s' to an integer: %v\n", expected_port_str, err)
		}
		expected_ports_ints_slice = append(expected_ports_ints_slice, expected_port_int)
	}

	// Run Nmap and record output
	command := exec.Command("nmap", "-PN", "--min-parallelism", *parallelism,
		"-n", "-sS", fmt.Sprintf("-p%s", *port_range), "--reason", *host)
	var stdout bytes.Buffer
	command.Stdout = &stdout
	var stderr bytes.Buffer
	command.Stderr = &stderr
	err := command.Run()
	var std_output = stdout.String()
	var std_error = stderr.String()
	if err != nil {
		log.Printf("Error with command.Run():\n%v\n", err)
		log.Fatalf("This is the stderr:\n%s\n", std_error)
	}

	//  grep Nmap for ports
	ports_regex := regexp.MustCompilePOSIX("^[0-9]*")
	found_ports_strings := strings.Fields(strings.Join(ports_regex.FindAllString(std_output, -1), " "))
	//  grep Nmap for 'Not shown'
	not_shown_regex := regexp.MustCompilePOSIX("^Not shown:.*")
	not_shown_strings := not_shown_regex.FindAllString(std_output, -1)
	////  grep 'Not shown' for 'closed'
	closed_regex := regexp.MustCompilePOSIX("closed")
	not_shown_closed := closed_regex.MatchString(strings.Join(not_shown_strings, " "))
	////  grep 'Not shown' for number of ports
	not_shown_ports_regex := regexp.MustCompilePOSIX("[0-9]*")
	not_shown_number := not_shown_ports_regex.FindAllString(strings.Join(not_shown_strings, " "), -1)

	// Convert greped ports string slice to int slice
	found_ports_ints_slice := []int{}
	for _, found_port_string := range found_ports_strings {
		found_port_int, err := strconv.Atoi(found_port_string)
		if err != nil {
			log.Fatalf("Error converting string '%s' to an integer: %v\n", found_port_string, err)
			//log.Fatalf("Error with strconv.Atoi: %v", err)
		}
		found_ports_ints_slice = append(found_ports_ints_slice, found_port_int)
	}

	// Compare and setup subject and message
	difference := compareMapAlternate(found_ports_ints_slice, expected_ports_ints_slice)
	var subject = fmt.Sprintf("Unfiltered ports found on %s", *host)
	var message string
	if len(difference) != 0 && not_shown_closed {
		message = fmt.Sprintf(
			"The following ports were found unfiltered and are not part of the expected set:\n\n%d.\n\nThere are %s closed ports that are not shown.\n\nThe expected set was:\n\n%d.",
			difference, not_shown_number, expected_ports_ints_slice,
		)
	} else if len(difference) != 0 {
		message = fmt.Sprintf(
			"The following ports were found unfiltered and are not part of the expected set:\n\n%d.\n\nThe expected set was:\n\n%d.",
			difference, expected_ports_ints_slice,
		)
	} else if not_shown_closed {
		message = fmt.Sprintf(
			"There are %s closed ports that are not shown.\n\nThe expected set is:\n\n%d.",
			not_shown_number, expected_ports_ints_slice,
		)
	}

	// Email if needed
	if len(difference) != 0 || not_shown_closed {
		const layout = "Mon, 2 Jan 2006 15:04:05 -0700"
		body := "From: " + *from_address + "\r\nTo: " + *to_addresses + "\r\nSubject: " + subject +
			"\r\nDate: " + time.Now().Format(layout) + "\r\n\r\n" + message
		domain, _, err := net.SplitHostPort(*server_address)
		if err != nil {
			log.Fatalf("Error with net.SplitHostPort: %v", err)
		}
		auth := smtp.PlainAuth("", *user_address, *password, domain)
		err = smtp.SendMail(*server_address, auth, *from_address,
			strings.Fields(*to_addresses), []byte(body))
		if err != nil {
			log.Printf("Error with smtp.SendMail: %v\n\n", err)
			log.Printf("Body: %v\n\n", body)
		}
	}
}
