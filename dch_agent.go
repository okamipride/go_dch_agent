package main

import (
	"bufio"
	"flag"
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

var did_prefix string = "1234567890" + "1234567890"
var DELAY_MS = 50 * time.Millisecond

//12345678901234567890000000000001
var (
	//url = "172.31.4.183:80"

	//url = "172.31.13.171:80"
	//url = "52.68.253.9"
	url = "r0402.dch.dlink.com:80"
	//url = "r0101.dch.dlink.com:80"
	//url       = "52.68.198.236:80"
	resp_data = "\"status\":\"ok\"," +
		"\"errno\":\"\"," +
		"\"errmsg\":\"\"," +
		"\"version\":\"1.0.1\"," +
		"\"detail\":\"70b3852=2015-04-30 11:45:58 +80800\"\r\n"
		/*
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
		*/
	version_response = "HTTP/1.1 200 OK\r\nContent-type: text\r\nContent-Length: 2\r\n\r\nOK"
)

func main() {
	log.Println("Agent Start to connect ... ")

	relay_url, num_dev := readArg()

	num_concurrence := int64(1)
	var my_delay time.Duration = 10 * time.Millisecond
	//relay_url, num_dev, num_concurrence, my_delay := readNumDevice()

	if relay_url != "" {
		url = relay_url + ":80"

	}

	//go AutoGC()

	for i := int64(1); i <= num_dev; i++ {
		device := Device{usr_did: genDid(i), usr_hash: genDid(i)}
		go device.deviceRoutine()
		if num_dev%num_concurrence == 0 {
			time.Sleep(my_delay)
		}
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

	tcpaddr, err := net.ResolveTCPAddr("tcp", url)
	if err != nil {
		log.Println("error", err)
		log.Println("url = ", url)
		return // retry escape
	}

	conn, err := net.DialTCP("tcp", nil, tcpaddr)

	if err != nil {
		log.Println("connect error", err, "url = ", url)
		return // retry escape
	}
	log.Println("Connected to : " + url)

	conn.Write([]byte(post_msg))

	defer closeConn(conn)

	first := make([]byte, 1)
	buf := make([]byte, 1024*32)

	n, err := conn.Read(first)

	if err == io.EOF {
		log.Println("Read EOF")
		return
	}

	go dev.deviceRoutine()

	println(string(buf[0:n]))

	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Println("Read EOF")
			return
		}
		if n > 0 {
			SendMessage(conn, version_response)
			log.Println(string(buf[0:n]))
		}
	}

}

func contiRead(conn *net.TCPConn, connid int, cs chan bool) {
	//var read_msg string
	buf_recever := make([]byte, RECV_BUF_LEN)
	for {
		_, err := conn.Read(buf_recever)

		if err != nil {
			if err == io.EOF {
				log.Println("EOF:Disconnected")
				break
			} else {
				log.Println("ContiRead Error = ", err)
				continue
			}
		}
		cs <- true
		SendMessage(conn, version_response)
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
		log.Println("Error send request:", err.Error())
	} else {
		log.Println("Request sent")
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

func readArg() (string, int64) {

	serverPtr := flag.String("serv", "r0401.dch.dlink.com", "relay server address , no port included")
	//serverPtr := flag.String("serv", "172.31.4.183:80", "relay server address , no port included")
	devPtr := flag.Int64("dev", 1, "number of devices want to connect to relay server")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "command line arguments")
	flag.Parse()
	fmt.Println("serverPtr:", *serverPtr)
	fmt.Println("devPtr:", *devPtr)
	fmt.Println("tail:", flag.Args())

	return *serverPtr, *devPtr

}

func readNumDevice() (string, int64, int64, time.Duration) {
	var relay_addr = ""
	var num_did int64 = 0
	var num_concur int64 = 0
	var delay time.Duration = 100 * time.Millisecond

	for {
		// Number of devices connect to relayd
		consolereader := bufio.NewReader(os.Stdin)

		log.Print("Enter Relay Server Address : ")
		input, err := consolereader.ReadString('\n')

		reg, _ := regexp.Compile("^[a-zA-Z0-9.]*$") // only alphanumeric and dot

		if err != nil {
			log.Println("ReadString error! Retry again = ", err)
			continue
		}

		relay_addr = reg.FindString(input)

		log.Print("Enter Number of Devices: ")
		input, err = consolereader.ReadString('\n')
		if err != nil {
			log.Println("ReadString error! Retry again = ", err)
			continue
		}
		reg, _ = regexp.Compile("^[1-9][0-9]*") // Remove special character only take digits
		num := reg.FindString(input)

		if err != nil {
			log.Println("ReadString error! Please enter digits. error = ", err)
			continue
		}

		log.Println(string(num))

		num_devices, err := strconv.ParseInt(string(num), 0, 64)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Convert Number failed! Please re-enter")
			continue
		}
		num_did = num_devices

		// Number of devices connect to relayd continousely without delay
		log.Print("Enter Number of Concurrent Connect: ")
		input, err = consolereader.ReadString('\n')

		if err != nil {
			log.Println("ReadString error! Use number of devices. error = ", err)
			num_concur = num_did
		}

		concur := reg.FindString(input)
		num_concur, err = strconv.ParseInt(string(concur), 0, 64)
		if err != nil {
			log.Println("ReadString error! Use number of devices. error = ", err)
			num_concur = num_did
		}

		// Number of devices connect to relayd continousely without delay
		log.Print("Enter delay ms: ")
		input, err = consolereader.ReadString('\n')

		if err != nil {
			log.Println("ReadString error! Use 100ms. error = ", err)
			delay = 100 * time.Millisecond
		}

		delayStr := reg.FindString(input)
		ms, err := strconv.ParseInt(string(delayStr), 0, 64)

		if err != nil {
			log.Println("ReadString error! Use 100ms,  error = ", err)
			delay = 100 * time.Millisecond
		}

		delay = time.Duration(ms) * time.Millisecond

		break
	}
	return relay_addr, num_did, num_concur, delay
}
