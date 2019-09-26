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
)

var (
	conns map[string]net.Conn
)

func init() {
	conns = make(map[string]net.Conn)
}

func main() {
	serverAddr := config.SERVER_IP + ":" + strconv.Itoa(config.SERVER_PORT)
	conn, err := net.Dial("udp", serverAddr)
	misc.CheckError(err)
	defer conn.Close()
	fmt.Println(serverAddr)

	go func() {
		for {
			buff := make([]byte, 1000)
			_, err := conn.Read(buff)
			if err != nil {
				fmt.Println(err)
			} else {
				msgHandle(buff)
			}
		}
	}()
	i := 1
	for {
		time.Sleep(time.Second * 3)
		line := fmt.Sprintf("send message :%d", i)
		conn.Write([]byte(line))
		i++
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
	fmt.Println(msgInfo)
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

func sendUdpMsg(address string) {
	conn, ok := conns[address]
	var err error
	if !ok {
		conn, err = net.Dial("udp", address)
		if err != nil {
			fmt.Println(err)
			return
		}
		conns[address] = conn
	}

	info := &misc.Msg{
		misc.CLUSTER_MSG,
		"hello",
	}

	data, err := json.Marshal(info)
	n, err := conn.Write(data)
	fmt.Println(n, err)
}
