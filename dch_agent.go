package main

import (
	//"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	//"regexp"
	"runtime"
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

var did_list = [5]string{
	"12345678901234567890123456789001",
	"12345678901234567890123456789002",
	"12345678901234567890123456789003",
	"12345678901234567890123456789004",
	"12345678901234567890123456789005"}

var (
	urlip            = "52.68.172.23:80"
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
	runtime.SetBlockProfileRate(1)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Println("Agent Start to connect ... ")
	//go AutoGC()

	deviceRoutine(did_list[0], did_list[0])

	for {
		time.Sleep(time.Second * 2)
	}
	//printMemory()

	log.Println("Program is going to exit")

}

//func SendRoutine(cs chan bool)
func deviceRoutine(did string, hash string) {

	post_msg := "POST /connect HTTP/1.1\r\n\r\n" + "\"did\":\"" + did + "\"" +
		"\r\n" + "\"hash\":\"" + hash + "\"" + "\r\n"

	connction_id := 0
	tcpaddr, err := net.ResolveTCPAddr("tcp", urlip)

	if checkError(err, "Agent ResolveTCPAddr") {
		os.Exit(0)
	}

	//var conn *net.TCPConn

	for {

		conn, err := net.DialTCP("tcp", nil, tcpaddr)

		if checkError(err, "Agent DialTCP") {
			continue
		}
		connction_id = connction_id + 1
		log.Println("Agent TCP Connect ... ", strconv.Itoa(connction_id))

		_, err = conn.Write([]byte(post_msg))

		if checkError(err, "Agent Write") {
			continue
		}

		buf := make([]byte, RECV_BUF_LEN)
		_, get_err := GetMessage(conn, buf)

		if get_err != nil {
			if get_err == io.EOF {
				log.Println("Disconnected")
				break
			} else {
				log.Println("error")
				continue
			}
		}

		//c := make(chan bool)

		go contiRead(conn, connction_id)

		//<-c
		//printMemory()
		log.Println("recieve msg go next round")

	}
	//defer closeConn(conn)
	printMemory()
	//closeConn(conn)
	log.Println("SendRoutine Exit")
}

func contiRead(conn *net.TCPConn, connid int) {
	//var read_msg string
	buf_recever := make([]byte, RECV_BUF_LEN)
	var read_msg string
	var get_err error

	for {
		read_msg, get_err = GetMessage(conn, buf_recever)

		if get_err != nil {
			if get_err == io.EOF {
				log.Println("Disconnected")
				break
			} else {
				log.Println("error")
				continue
			}
		}

		//cs <- true
		log.Println("Recieved msg = ", read_msg)
		SendMessage(conn, version_response)
	}
	conn.Close()
	printMemory()
	log.Println("Exist contiRead: ", strconv.Itoa(connid))
}

func parseGet(msg string) int {
	return 1
}

func AutoGC() {
	for {
		//log.Println("System auto GC...")
		runtime.GC()
		//fmt.Println("Auto GC")
		time.Sleep(1 * time.Second)

	}
}

func GetMessage(conn *net.TCPConn, buf []byte) (string, error) {
	//fmt.Println("Prepare GetMessage ...")
	//conn.SetReadDeadline (time.Time{})
	n, err := conn.Read(buf)

	//read_data := make([]byte, n)
	//copy(read_data, buf_recever)
	return string(buf[0:n]), err
}

func SendMessage(conn *net.TCPConn, msg string) {
	log.Println("Prepare SendMessage ...")
	_, err := conn.Write([]byte(msg))
	if err != nil {
		println("Error send request:", err.Error())
	} else {
		println("Request sent")
	}
}

func closeConn(c *net.TCPConn) {
	log.Println("connection close")
	c.Close()
}

func checkError(err error, act string) bool {

	if err != nil {
		log.Println(act + " Error Occur!")
		//fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		log.Printf("%s Fatal error: %s", err.Error())
		return true
	}
	//fmt.Println(act + " no Error")
	return false
}

func printMemory() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Println("Memory Usage")
	log.Println("Alloc : ", mem.Alloc)
	log.Println("Total Alloc : ", mem.TotalAlloc)
	log.Println("HeapAlloc : ", mem.HeapAlloc)
	log.Println("HeapSys : ", mem.HeapSys)
}
