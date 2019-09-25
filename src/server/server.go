package main
 
import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
 
	"../config"
)
 
func main() {

	address := config.SERVER_IP + ":" + strconv.Itoa(config.SERVER_PORT)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
 
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
 
	defer conn.Close()
 
	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte, 100)
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
 
		strData := string(data)
		fmt.Println("Received:", strData)
 
		upper := strings.ToUpper(strData)
		_, err = conn.WriteToUDP([]byte(upper), rAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
 
		fmt.Println("Send:", upper)
	}
}
