package main

import (
	"fmt"
	"net"
	//"os"
	//"strings"

	"../config"
	"../misc"
	"encoding/binary"
	"encoding/json"
)

var (
	addr2SvId map[string]int32
	id        int32
)

func init() {
	addr2SvId = make(map[string]int32)
	id = 0
}

func getId() int32 {
	id++
	return id
}

func main() {

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: config.SERVER_PORT,
	})
	misc.CheckError(err)
	defer conn.Close()
	for {
		data := make([]byte, 100)
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		strData := string(data)
		fmt.Println("Received:", strData, "from", rAddr)

		reply := msgHandle(rAddr)
		data = pkt(reply)
		_, err = conn.WriteToUDP(data, rAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Send:", reply)
	}
}

func pkt(str string) []byte {
	msglen := len(str)
	buff := make([]byte, msglen+2)
	binary.BigEndian.PutUint16(buff, uint16(msglen))
	copy(buff[2:], str)
	return buff
}

func msgHandle(addr *net.UDPAddr) string {
	addStr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	id, ok := addr2SvId[addStr]
	if !ok {
		id = getId()
		addr2SvId[addStr] = id
	}
	return pktNodeInfo(id)
}

func pktNodeInfo(id int32) string {
	cluster := &misc.ClusterInfo{
		Id:   id,
		Node: make([]*misc.AddressInfo, 0),
	}

	for addr, idx := range addr2SvId {
		cluster.Node = append(cluster.Node, &misc.AddressInfo{idx, addr})
	}
	str, _ := json.Marshal(cluster)
	msg := &misc.Msg{
		Type: misc.CLUSTER_TYPE,
		Data: string(str),
	}
	str, _ = json.Marshal(msg)
	return string(str)
}
