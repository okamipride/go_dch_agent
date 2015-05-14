package main

import (
	"fmt"
	"io"
	//"log"
	"net"
	"os"
	//"regexp"
	"time"
)

const (
	RECV_BUF_LEN     = 1024 * 10
	CONFIG_LINE_SIZE = 1024
)

type Device struct {
	procid   int
	usr_did  string
	usr_hash string
}

var (
	urlip    = "52.68.172.23:80"
	usr_did  = "12345678901234567890123456789012"
	usr_hash = "12345678901234567890123456789012"
	post_msg = "POST /connect HTTP/1.1\r\n\r\n" + "\"did\":\"" + usr_did + "\"" +
		"\r\n" + "\"hash\":\"" + usr_hash + "\"" + "\r\n"

	version_response = "HTTP/1.1 200 OK\r\n" + "Server: Spaced/0.1\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Type: application/javascript; charset=utf-8\r\n" +
		"Connection: close\r\n\r\n" +
		"1\r\n" +
		"{\r\n" +
		"64\r\n" +
		"\"status\":\"ok\"," +
		"\"errno\":\"\"," +
		"\"errmsg\":\"\"," +
		"\"version\":\"1.0.1\"," +
		"\"detail\":\"70b3852=2015-04-30 11:45:58 +80800\"\r\n" +
		"1\r\n" +
		"}\r\n" +
		"0\r\n"
)

func main() {
	fmt.Println("Agent Start to connect ... ")
	go SendRoutine()
	for {
		time.Sleep(time.Second * 1)
	}
	fmt.Println("Program is going to exit")

}

//func SendRoutine(cs chan bool)
func SendRoutine() {
	connNum := 0
	tcpaddr, err := net.ResolveTCPAddr("tcp", urlip)
	if checkError(err, "Agent ResolveTCPAddr") {
		os.Exit(0)
	}

RECONNECT:
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	defer closeConn(conn)

	if checkError(err, "Agent DialTCP") {
		goto RECONNECT
	}
	connNum = connNum + 1
	fmt.Println("Agent TCP Connect ... ")

	_, err = conn.Write([]byte(post_msg))

	if checkError(err, "Agent Write") {
		goto RECONNECT
	}

	fmt.Println("Agent Write did & hash")

	for {
		readMsg := GetMessage(conn)
		if readMsg == "" {
			fmt.Println("echo empty")
			continue
		}

		if readMsg == "EOF" {
			fmt.Println("Disconnected")
			break
		}

		fmt.Println("Recieved msg = ", readMsg)
		SendMessage(conn, version_response)
	}

	fmt.Println("SendRoutine Exist")
}

func parseGet(msg string) int {
	return 1
}

func GetMessage(conn *net.TCPConn) string {
	//fmt.Println("Prepare GetMessage ...")
	buf_recever := make([]byte, RECV_BUF_LEN)
	//conn.SetReadDeadline (time.Time{})
	n, err := conn.Read(buf_recever)
	if err != nil {
		if err != io.EOF {
			println("Error while receive response:", err.Error())
			fmt.Println("Get Message Error")
			return ""
		} else {
			return "EOF"
		}

	}

	read_data := make([]byte, n)
	copy(read_data, buf_recever)
	return string(read_data)
}

func SendMessage(conn *net.TCPConn, msg string) {
	fmt.Println("Prepare SendMessage ...")
	_, err := conn.Write([]byte(msg))
	if err != nil {
		println("Error send request:", err.Error())
	} else {
		println("Request sent")
	}
}

func closeConn(c *net.TCPConn) {
	fmt.Println("connection close")
	c.Close()
}

func checkError(err error, act string) bool {

	if err != nil {
		fmt.Println(act + " Error Occur!")
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return true
	}
	//fmt.Println(act + " no Error")
	return false
}
