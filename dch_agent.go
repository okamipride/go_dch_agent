package main

import (
	"fmt"
	"io"
	"log"
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

//func SendRoutine(cs chan bool)
func SendRoutine() {

	tcpaddr, err := net.ResolveTCPAddr("tcp", urlip)
	if checkError(err, "Agent ResolveTCPAddr") {
		os.Exit(0)
	}

RECONNECT:
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	defer conn.Close() // close when leave the loop
	if checkError(err, "Agent DialTCP") {
		goto RECONNECT
	}

	fmt.Println("App TCP Connect ... ")

	_, err = conn.Write([]byte(post_msg))
	log.Printf("Agent write:\n%s", post_msg)
	if checkError(err, "Agent Write") {
		goto RECONNECT
	}

	for {
		echo := GetMessage(conn)
		if echo == "" {
			fmt.Println("echo empty")
			break
		}
		println("receive success")
		SendMessage(conn, version_response)
		println("sendEcho")
	}

	fmt.Println("SendRoutine Exist")
}

func parseGet(msg string) int {
	return 1
}

func main() {
	fmt.Println("Agent Start to connect ... ")
	go SendRoutine()
	for {
		time.Sleep(time.Second * 1)
	}
	fmt.Println("Program is going to exit")

}

func GetMessage(conn *net.TCPConn) string {
	buf_recever := make([]byte, RECV_BUF_LEN)
	//conn.SetReadDeadline(time.Time{})
	n, err := conn.Read(buf_recever)
	if err != nil {
		if err != io.EOF {
			println("Error while receive response:", err.Error())
			fmt.Println("Get Message Error")
			return ""
		}
	}

	echodata := make([]byte, n)
	copy(echodata, buf_recever)
	fmt.Println("App recieve message:%s", string(echodata))
	return string(echodata)
}

func SendMessage(conn *net.TCPConn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		println("Error send request:", err.Error())
	} else {
		println("Request sent")
	}
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

func checkReadError(err error, act string) bool {

	if err != nil {
		if err != io.EOF {
			fmt.Println(act + " Error Occur!")
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			return true
		} else {
			fmt.Println("EOF")
		}
	}
	//fmt.Println(act + " no Error")
	return false
}
