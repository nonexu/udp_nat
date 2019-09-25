package main
 
import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
 	"../misc"
	"../config"
)

func main() {
	serverAddr := config.SERVER_IP + ":" + strconv.Itoa(config.SERVER_PORT)
	conn, err := net.Dial("udp", serverAddr)
	misc.CheckError(err)
	defer conn.Close()
 
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		line := input.Text()
		conn.Write([]byte(line))

 		buff := make([]byte, 100)
		_, err := conn.Read(buff)
		if err != nil {
			fmt.Println(err)
		}else {
			fmt.Println(string(buff))			
		}
	}
}
