package main

import (
	//"bufio"
	"fmt"
	"net"
	//"os"
	"../config"
	"../misc"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"time"
	"strings"
)

var (
	conns map[string]net.Conn
	srcAddr *net.UDPAddr
	dstAddr *net.UDPAddr
	conn net.Conn
)

func init() {
	conns = make(map[string]net.Conn)
	srcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: 10005}
	dstAddr = &net.UDPAddr{IP: net.ParseIP(config.SERVER_IP), Port: config.SERVER_PORT}
}

func main() {
	var err error
	conn, err = net.DialUDP("udp", srcAddr, dstAddr)
	misc.CheckError(err)
	defer conn.Close()

	i := 1
	for {
		time.Sleep(time.Second * 1)
		line := fmt.Sprintf("send message :%d", i)
		conn.Write([]byte(line))
		i++
			buff := make([]byte, 1000)
			_, err := conn.Read(buff)
			if err != nil {
				fmt.Println(err)
			} else {
				msgHandle(buff)
		}

	}
}

func msgHandle(msg []byte) {
	msgLen := uint32(binary.BigEndian.Uint16(msg[0:2]))
	msgInfo := &misc.Msg{}
	err := json.Unmarshal(msg[2:2+msgLen], msgInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msgInfo.Type {
	case misc.CLUSTER_TYPE:
		clusterTypeHandle(msgInfo.Data)
	case misc.CLUSTER_MSG:
		clusterMsgHandle(msgInfo.Data)
	default:
		fmt.Println("unknown msg")
	}
}

func clusterTypeHandle(msg string) {
	cluster := &misc.ClusterInfo{}
	err := json.Unmarshal([]byte(msg), cluster)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, node := range cluster.Node {
		if node.Id == cluster.Id {
			continue
		}
		sendUdpMsg(node.Address)
	}
}

func clusterMsgHandle(msg string) {
	fmt.Println(msg)
}

func parse(address string) *net.UDPAddr{
     t := strings.Split(address, ":")
     port, _ := strconv.Atoi(t[1])
	return &net.UDPAddr{IP: net.ParseIP(t[0]), Port: port}
}

func sendUdpMsg(address string) {
	otherdstAddr := parse(address)
	conn.Close()
	var err error
	conn, err = net.DialUDP("udp", srcAddr, otherdstAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(srcAddr,otherdstAddr)
	go read()
	for {
		time.Sleep(time.Second*3)
		info := &misc.Msg{
			misc.CLUSTER_MSG,
			"hello",
		}
		data, _ := json.Marshal(info)
		n, err := conn.Write(data)
		fmt.Println(n ,err)
	}
}

func read() {
	for {
		buff := make([]byte, 100)
		_, err := conn.Read(buff)
		fmt.Println(string(buff), err)

	}
}
