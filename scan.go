package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/smtp"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Execute(output_buffer *bytes.Buffer, stack ...*exec.Cmd) (err error) {
	var error_buffer bytes.Buffer
	pipe_stack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdin_pipe, stdout_pipe := io.Pipe()
		stack[i].Stdout = stdout_pipe
		stack[i].Stderr = &error_buffer
		stack[i+1].Stdin = stdin_pipe
		pipe_stack[i] = stdout_pipe
	}
	stack[i].Stdout = output_buffer
	stack[i].Stderr = &error_buffer

	if err := call(stack, pipe_stack); err != nil {
		log.Fatalln(string(error_buffer.Bytes()), err)
	}
	return err
}

func call(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = call(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}

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
	// Host and it's expected ports
	host := os.Getenv("HOST")
	expected_ports_str := os.Getenv("EXPECTED_PORTS")
	port_range := os.Getenv("PORT_RANGE")

	// Nmap --min-parallelism
	parallelism := os.Getenv("PARALLELISM")

	// SMTP stuff
	to_addresses := strings.Fields(os.Getenv("TO_ADDRESSES"))
	from_address := os.Getenv("FROM_ADDRESS")
	server_address := os.Getenv("SMTP_SERVER_ADDRESS")
	user_address := os.Getenv("SMTP_USER_ADDRESS")
	password := os.Getenv("SMTP_PASSWORD")

	// Convert EXPECTED_PORTS string to int slice
	expected_ports_str_slice := strings.Fields(expected_ports_str)
	expected_ports_ints_slice := []int{}
	for _, expected_port_str := range expected_ports_str_slice {
		expected_port_int, err := strconv.Atoi(expected_port_str)
		if err != nil {
			log.Fatalf("Error converting string '%s' to an integer: %v\n", expected_port_str, err)
		}
		expected_ports_ints_slice = append(expected_ports_ints_slice, expected_port_int)
	}

	// Run Nmap and grep output
	var found_ports_data bytes.Buffer
	err := Execute(&found_ports_data,
		exec.Command(
			"nmap", "-PN", "--min-parallelism", parallelism,
			"-n", "-sS", fmt.Sprintf("-p%s", port_range), "--reason", host,
		),
		exec.Command("grep", "-o", "^[0-9]*"),
	)
	if err != nil {
		log.Fatalln("Error with Execute: %v", err)
	}

	// Convert greped output string to int slice
	found_ports_ints_slice := []int{}
	for _, found_port_string := range strings.Fields(found_ports_data.String()) {
		found_port_int, err := strconv.Atoi(found_port_string)
		if err != nil {
			log.Fatalln("Error with strconv.Atoi: %v", err)
		}
		found_ports_ints_slice = append(found_ports_ints_slice, found_port_int)
	}

	// Compare and email if needed
	difference := compareMapAlternate(found_ports_ints_slice, expected_ports_ints_slice)
	if len(difference) != 0 {

		subject := fmt.Sprintf("Unfiltered ports found on %s", host)
		message := fmt.Sprintf(
			"The following ports were found unfiltered and are not part of the expected set:\r\r%d.\r\rThe expected set is:\r\r%d.",
			difference, expected_ports_ints_slice,
		)
		const layout = "Mon, 2 Jan 2006 15:04:05 -0700"
		body := "From: " + from_address + "\r\nTo: " + strings.Join(to_addresses, " ") + "\r\nSubject: " + subject +
			"\r\nDate: " + time.Now().Format(layout) + "\r\n\r\n" + message
		domain, _, err := net.SplitHostPort(server_address)
		if err != nil {
			log.Fatalf("Error with net.SplitHostPort: %v", err)
		}
		auth := smtp.PlainAuth("", user_address, password, domain)
		err = smtp.SendMail(server_address, auth, from_address,
			to_addresses, []byte(body))
		if err != nil {
			log.Fatalf("Error with smtp.SendMail: %v", err)
		}
	}
}
