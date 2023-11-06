package unpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	_ "os"
	"strconv"
	"strings"
)

type Protocol struct {
	Format []string
}

// 封包
func (p *Protocol) Pack(args ...any) []byte {
	la := len(args)
	ls := len(p.Format)
	ret := []byte{}
	if ls > 0 && la > 0 && ls == la {
		for i := 0; i < ls; i++ {
			if p.Format[i] == "H" {
				ret = append(ret, IntToBytes2(args[i].(int))...)
			} else if p.Format[i] == "I" {
				ret = append(ret, IntToBytes4(args[i].(int))...)
			} else if strings.Contains(p.Format[i], "s") {
				num, _ := strconv.Atoi(strings.TrimRight(p.Format[i], "s"))
				ret = append(ret, []byte(fmt.Sprintf("%s%s", args[i].(string), strings.Repeat("\x00", num-len(args[i].(string)))))...)
			}
		}
	}
	return ret
}

// 解包
func (p *Protocol) UnPack(msg []byte) []any {
	la := len(p.Format)
	ret := make([]interface{}, la)
	if la > 0 {
		for i := 0; i < la; i++ {
			if p.Format[i] == "H" {
				ret[i] = Bytes4ToInt(msg[0:2])
				msg = msg[2:len(msg)]
			} else if p.Format[i] == "I" {
				ret[i] = Bytes4ToInt(msg[0:4])
				msg = msg[4:len(msg)]
			} else if strings.Contains(p.Format[i], "s") {
				num, _ := strconv.Atoi(strings.TrimRight(p.Format[i], "s"))
				ret[i] = string(msg[0:num])
				msg = msg[num:len(msg)]

			}
		}
	}
	return ret
}

func (p *Protocol) Size() int {
	size := 0
	ls := len(p.Format)
	if ls > 0 {
		for i := 0; i < ls; i++ {
			if p.Format[i] == "H" {
				size = size + 2
			} else if p.Format[i] == "I" {
				size = size + 4
			} else if strings.Contains(p.Format[i], "s") {
				num, _ := strconv.Atoi(strings.TrimRight(p.Format[i], "s"))
				size = size + num
			}
		}
	}
	return size
}

// 整形转换成字节
func IntToBytes(n int) []byte {
	m := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, m)

	gbyte := bytesBuffer.Bytes()

	return gbyte
}

// 整形转换成字节4位
func IntToBytes4(n int) []byte {
	m := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, m)

	gbyte := bytesBuffer.Bytes()
	//c++ 高低位转换
	k := 4
	x := len(gbyte)
	nb := make([]byte, k)
	for i := 0; i < k; i++ {
		nb[i] = gbyte[x-i-1]
	}
	return nb
}

// 整形转换成字节2位
func IntToBytes2(n int) []byte {
	m := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, m)

	gbyte := bytesBuffer.Bytes()
	//c++ 高低位转换
	k := 2
	x := len(gbyte)
	nb := make([]byte, k)
	for i := 0; i < k; i++ {
		nb[i] = gbyte[x-i-1]
	}
	return nb
}

// 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

// 4个字节转换成整形
func Bytes4ToInt(b []byte) int {
	xx := make([]byte, 4)
	if len(b) == 2 {
		xx = []byte{b[0], b[1], 0, 0}
	} else {
		xx = b
	}

	m := len(xx)
	nb := make([]byte, 4)
	for i := 0; i < 4; i++ {
		nb[i] = xx[m-i-1]
	}
	bytesBuffer := bytes.NewBuffer(nb)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
