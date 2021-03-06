package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

const (
	RECV_BUF_LEN     = 1024 * 10
	CONFIG_LINE_SIZE = 1024
)

type Device struct {
	usr_did  string
	usr_hash string
}

var did_prefix = "1234567890" + "1234567890"

//12345678901234567890000000000001
var (
	url = "r0401.dch.dlink.com:80"
	//url       = "r0101.dch.dlink.com:80"
	resp_data = "\"status\":\"ok\"," +
		"\"errno\":\"\"," +
		"\"errmsg\":\"\"," +
		"\"version\":\"1.0.1\"," +
		"\"detail\":\"70b3852=2015-04-30 11:45:58 +80800\"\r\n"
	version_response = "HTTP/1.1 200 OK\r\n" + "Server: Spaced/0.1\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Type: application/javascript; charset=utf-8\r\n" +
		"Connection: close\r\n\r\n" +
		"1\r\n" +
		"{\r\n" +
		"64\r\n" +
		resp_data +
		"1\r\n" +
		"}\r\n" +
		"0\r\n"
)

func main() {
	log.Println("Agent Start to connect ... ")
	num_dev := readNumDevice()
	go AutoGC()

	for i := int64(1); i <= num_dev; i++ {
		device := Device{usr_did: genDid(i), usr_hash: genDid(i)}
		go device.deviceRoutine()
	}

	for {
		time.Sleep(time.Second * 2)
	}

	log.Println("Program is going to exit")

}

func genDid(num int64) string {
	return did_prefix + fmt.Sprintf("%012x", num)
}

//func SendRoutine(cs chan bool)
func (dev *Device) deviceRoutine() {

	post_msg := "POST /connect HTTP/1.1\r\n\r\n" + "\"did\":\"" + dev.usr_did + "\"" +
		"\r\n" + "\"hash\":\"" + dev.usr_hash + "\"" + "\r\n"

	connction_id := 0
	tcpaddr, err := net.ResolveTCPAddr("tcp", url)

	if checkError(err, "Agent ResolveTCPAddr") {
		os.Exit(0)
	}

	var conn *net.TCPConn

	conn, err = net.DialTCP("tcp", nil, tcpaddr)
	conn.SetKeepAlive(true)
	if checkError(err, "Agent DialTCP") {
		return
	}

	connction_id = connction_id + 1
	//log.Println("Agent TCP Connect ... ", strconv.Itoa(connction_id))

	_, err = conn.Write([]byte(post_msg))

	if checkError(err, "Agent Write") {
		return
	}

	buf_recever := make([]byte, RECV_BUF_LEN)
	_, err = conn.Read(buf_recever)

	if err != nil {
		if err == io.EOF {
			log.Println("Disconnected")
			closeConn(conn)
			return
		} else {
			log.Println("Read error", err)
			closeConn(conn)
			return
		}
	}
	closeConn(conn)
	log.Println("Read & connection closed")

	//closeConn(conn)
	//log.Println("SendRoutine Exit")
}

func contiRead(conn *net.TCPConn, connid int, cs chan bool) {
	//var read_msg string
	for {

		buf_recever := make([]byte, RECV_BUF_LEN)
		_, err := conn.Read(buf_recever)

		if err != nil {
			if err == io.EOF {
				log.Println("Disconnected")
				break
			} else {
				log.Println("error")
				continue
			}
		}
		cs <- true
		SendMessage(conn, version_response)

		//log.Println("Recieved msg = ", string(buf_recever[0:n]))

	}

	defer closeConn(conn)
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

func readNumDevice() int64 {
	var num_did int64 = 0

	for {
		consolereader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Number of Devices: ")
		input, err := consolereader.ReadString('\n') // this will prompt the user for input

		if err != nil {
			fmt.Println(err)
			fmt.Println("ReadString error! Retry again")
			continue
		}
		reg, _ := regexp.Compile("^[1-9][0-9]*") // Remove special character only take digits
		num := reg.FindString(input)

		if err != nil {
			fmt.Println(err)
			fmt.Println("ReadString error! Please enter digits")
			continue
		}

		fmt.Println(string(num))

		num_devices, err := strconv.ParseInt(string(num), 0, 64)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Convert Number failed! Please re-enter")
			continue
		}
		num_did = num_devices
		break
		//return num_devices
	}

	return num_did
}
