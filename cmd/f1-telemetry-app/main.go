package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

func floatFromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func bits32FromBuf(buf []byte, i int) []byte {
	return buf[4*i : 4*(i+1)]
}

func main() {

	sAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:20777")
	if err != nil {
		fmt.Errorf("Error: ", err)
	}

	sConn, err := net.ListenUDP("udp", sAddr)
	if err != nil {
		fmt.Errorf("Error: ", err)
	}
	defer sConn.Close()

	buf := make([]byte, 1289)

	for {
		_, _, err := sConn.ReadFromUDP(buf)
		fmt.Println("Received Speed :", floatFromBytes(bits32FromBuf(buf, 7))*3.6)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
	}
}
