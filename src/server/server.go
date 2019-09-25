package main
 
import (
	"fmt"
	"net"
	"os"
	"strings"
 
	"../config"
)
 
func main() {

      conn, err := net.ListenUDP("udp4", &net.UDPAddr{
        IP:   net.IPv4(0, 0, 0, 0),
        Port: config.SERVER_PORT,
    })

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
