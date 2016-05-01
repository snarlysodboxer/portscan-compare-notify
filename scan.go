package main

import (
	"bytes"
	"fmt"
	flag "github.com/ogier/pflag"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Returns the expected-unfound and unexpected-found
func comparePorts(expectedPorts, foundPorts []int) ([]int, []int) {
	expectedUnfound := compare(expectedPorts, foundPorts)
	unexpectedFound := compare(foundPorts, expectedPorts)
	return expectedUnfound, unexpectedFound
}

func compare(X, Y []int) []int {
	counts := make(map[int]int)
	var total int
	for _, val := range X {
		counts[val] += 1
		total += 1
	}
	for _, val := range Y {
		if _, ok := counts[val]; ok {
			counts[val] -= 1
			total -= 1
			if counts[val] <= 0 {
				delete(counts, val)
			}
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

// Converts a space separated string of whole numbers to a slice of integers
func convertStringToIntSlice(str *string) []int {
	substrings := strings.Fields(*str)
	integers := []int{}
	for _, substring := range substrings {
		integer, err := strconv.Atoi(substring)
		if err != nil {
			log.Fatalf("Error converting string '%s' to an integer:\n%v\n", substring, err)
		}
		integers = append(integers, integer)
	}
	return integers
}

func nmapRun(nmapOptionsSlice *[]string) string {
	command := exec.Command("nmap", *nmapOptionsSlice...)
	log.Printf("Running: \"nmap %s\"\n", strings.Join(*nmapOptionsSlice, " "))
	var stdout bytes.Buffer
	command.Stdout = &stdout
	var stderr bytes.Buffer
	command.Stderr = &stderr
	err := command.Run()
	stdOut := stdout.String()
	stdErr := stderr.String()
	if err != nil || stdErr != "" {
		if err != nil {
			log.Printf("Error with command.Run(): %v\n", err)
		}
		if stdErr != "" {
			log.Printf("Stderr:\n%s", stdErr)
		}
		log.Fatal("Exiting!")
	}
	return stdOut
}

func removeEmptyStrings(strings []string) []string {
	var newStrings []string
	for _, str := range strings {
		if str != "" {
			newStrings = append(newStrings, str)
		}
	}
	return newStrings
}

// foundPortsString, notShownBool, notShownQuantity
func grepNmap(nmapOutput string) (string, bool, []string) {
	portsRegex := regexp.MustCompilePOSIX("^[0-9]*")
	foundPortsString := strings.Join(removeEmptyStrings(portsRegex.FindAllString(nmapOutput, -1)), " ")
	//  grep Nmap for 'Not shown'
	notShownRegex := regexp.MustCompilePOSIX("^Not shown:.*")
	notShownString := strings.Join(removeEmptyStrings(notShownRegex.FindAllString(nmapOutput, -1)), " ")
	////  grep 'Not shown' for 'closed'
	closedRegex := regexp.MustCompilePOSIX("closed")
	notShownBool := closedRegex.MatchString(notShownString)
	////  grep 'Not shown' for number of ports
	notShownPortsRegex := regexp.MustCompilePOSIX("[0-9]*")
	notShownQuantity := removeEmptyStrings(notShownPortsRegex.FindAllString(notShownString, -1))
	return foundPortsString, notShownBool, notShownQuantity
}

func message(expectedPorts, foundPorts []int, notShownBool bool, notShownQuantity []string, expectedUnfoundPorts, unexpectedFoundPorts []int) string {
	var message string
	if len(expectedUnfoundPorts) != 0 {
		message += fmt.Sprintf(
			"The following ports were found filtered but were expected to be unfiltered:\n\n%d.\n\n",
			expectedUnfoundPorts,
		)
	}
	if len(unexpectedFoundPorts) != 0 {
		message += fmt.Sprintf(
			"The following ports were found unfiltered and are not part of the expected set:\n\n%d.\n\n",
			unexpectedFoundPorts,
		)
	}
	if notShownBool {
		message += fmt.Sprintf("There are %s unfiltered ports that are not shown.\n\n", notShownQuantity)
	}
	message += fmt.Sprintf("The expected set was:\n\n%d.\n\n", expectedPorts)
	message += "\t[[Unfiltered = 'open' or 'closed' (firewall is open)]]\n"
	message += "\t[[Filtered = no response/timeout (firewall is closed)]]\n"
	message += "\t[['open' usually = SYN-ACK]]\n"
	message += "\t[['closed' usually = RST]]\n"
	message += "\t[[See Nmap manual for more details.]]"
	return message
}

func UID(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func main() {
	log.SetPrefix(fmt.Sprintf("[%s] ", UID(8)))
	// Flags
	//// NMAP options and expected ports
	nmapOptions := flag.StringP("nmapoptions", "o", "", "options to pass to NMAP")
	expectedPortsString := flag.StringP("expected", "e", "",
		"space separated list of ports that are expected to be found unfiltered (unfiltered = open or closed)")

	//// SMTP stuff
	toAddresses := flag.StringP("to", "t", "", "space separated list of 'to' address(es)")
	fromAddress := flag.StringP("from", "f", "", "the 'from' email address")
	serverAddress := flag.StringP("smtpserver", "s", "", "the SMTP server address, E.G 'smtp.example.com:587'")
	userAddress := flag.StringP("username", "u", "", "the SMTP username or email address")
	password := flag.StringP("password", "x", "", "the SMTP user password")

	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NFlag() != 7 {
		flag.Usage()
		return
	}

	nmapOptionsSlice := strings.Split(*nmapOptions, " ")

	// Get host out of nmapOptions
	host := nmapOptionsSlice[len(nmapOptionsSlice)-1]

	// Convert EXPECTED_PORTS string to int slice
	expectedPorts := convertStringToIntSlice(expectedPortsString)

	// Run Nmap and record output
	startTime := time.Now()
	nmapOutput := nmapRun(&nmapOptionsSlice)
	endTime := time.Now()
	log.Printf("The NMAP run took %s", endTime.Sub(startTime))

	// Grep Nmap output
	foundPortsString, notShownBool, notShownQuantity := grepNmap(nmapOutput)

	// Convert greped ports string to int slice
	foundPorts := convertStringToIntSlice(&foundPortsString)

	// Compare
	expectedUnfoundPorts, unexpectedFoundPorts := comparePorts(expectedPorts, foundPorts)

	// Setup message
	message := message(expectedPorts, foundPorts, notShownBool, notShownQuantity,
		expectedUnfoundPorts, unexpectedFoundPorts)

	// Email if needed
	if len(expectedUnfoundPorts) != 0 || len(unexpectedFoundPorts) != 0 || notShownBool {
		log.Printf("Unexpected unfiltered ports found on %s, sending Alert Email\n", host)
		log.Printf("Results:\n%s\n", message)
		var subject = fmt.Sprintf("Unexpected unfiltered ports found on %s", host)
		const dateLayout = "Mon, 2 Jan 2006 15:04:05 -0700"
		body := "From: " + *fromAddress + "\r\nTo: " + *toAddresses + "\r\nSubject: " + subject +
			"\r\nDate: " + time.Now().Format(dateLayout) + "\r\n\r\n" + message
		domain, _, err := net.SplitHostPort(*serverAddress)
		if err != nil {
			log.Fatalf("Error with net.SplitHostPort: %v", err)
		}
		auth := smtp.PlainAuth("", *userAddress, *password, domain)
		err = smtp.SendMail(*serverAddress, auth, *fromAddress,
			strings.Fields(*toAddresses), []byte(body))
		if err != nil {
			log.Printf("Error with smtp.SendMail: %v\n\n", err)
			log.Printf("Body: %v\n\n", body)
		}
	} else {
		log.Println("Results: Everything is as expected")
	}
}
