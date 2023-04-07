package bat2txt

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net"
	"os"
)

//https://www.biaodianfu.com/qqwry-dat.html
func Convert()  {
	fs, err := os .Open("../app/qqwry.dat")
	if err != nil {
		log.Panic(err)
	}
	defer fs.Close()

	f, err := os.Create("../pure.txt")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	var byte4 = [4]byte{}
	fs.ReadAt(byte4[:], 0)
	var firstIdx uint32
	binary.Read(bytes.NewReader(byte4[:]), binary.LittleEndian, &firstIdx)
	log.Printf("first index %d", firstIdx)
	byte4 = [4]byte{}
	fs.ReadAt(byte4[:], 4)
	var lastIdx uint32
	binary.Read(bytes.NewReader(byte4[:]), binary.LittleEndian, &lastIdx)
	log.Printf("last index %d", lastIdx)

	for firstIdx <= lastIdx {
		var recordA string
		var recordB string

		var startIpBytes = [4]byte{}
		fs.ReadAt(startIpBytes[:], int64(firstIdx))
		var startIp uint32
		binary.Read(bytes.NewReader(startIpBytes[:]), binary.LittleEndian, &startIp)
		var offset = readOffset(fs, firstIdx+4)
		var endIpBytes = [4]byte{}
		fs.ReadAt(endIpBytes[:], int64(offset))
		var endIp uint32
		binary.Read(bytes.NewReader(endIpBytes[:]), binary.LittleEndian, &endIp)

		var modeByte = [1]byte{}
		fs.ReadAt(modeByte[:], int64(offset+4))
		var mode = modeByte[0]
		switch mode {
		case 1:
			var offset2 = readOffset(fs, offset + 5)
			fs.ReadAt(modeByte[:], int64(offset2))
			mode = modeByte[0]

			if mode == 2 {
				var offset3 = readOffset(fs, offset2 + 1)
				recordA,_ = readA(fs, offset3)
				recordB = readB(fs, offset2 + 4)
			} else {
				recordA, offset2 = readA(fs, offset2)
				recordB = readB(fs, offset2)
			}
		case 2:
			var offset2 = readOffset(fs, offset + 5)
			recordA, _ = readA(fs, offset2)
			recordB = readB(fs, offset + 8)
		default:
			recordA, offset = readA(fs, offset + 4)
			recordB = readB(fs, offset)
		}
		firstIdx+=7
		f.WriteString(fmt.Sprintf("%-16s%-16s%s %s\n", Uint32ToIpStr(startIp), Uint32ToIpStr(endIp), recordA, recordB))
	}
}

func readA(fs *os.File, offset uint32) (string, uint32){
	r, err := fs.Seek(int64(offset), 0)
	if err != nil {
		log.Panic(err)
	}
	reader := bufio.NewReader(fs)
	b, err := reader.ReadBytes(byte(0))
	if err != nil {
		log.Panic(err)
	}
	readed := len(b)
	b = b[:len(b)-1]
	b, err = simplifiedchinese.GBK.NewDecoder().Bytes(b)
	if err != nil {
		log.Panic(err)
	}
	return string(b), uint32(r + int64(readed))
}

func readB(fs *os.File, offset uint32) (recordB string) {
	var modeByte = [1]byte{}
	fs.ReadAt(modeByte[:], int64(offset))
	mode := modeByte[0]
	if mode == 1 || mode == 2 {
		var offset3 = readOffset(fs, offset + 1)
		recordB,_ = readA(fs, offset3)
	} else {
		recordB,_ = readA(fs, offset)
	}
	return
}

func Uint32ToIpStr(ip uint32) string {
	var b [4]byte
	b[0] = byte(ip & 0xff)
	b[1] = byte((ip >> 8) & 0xff)
	b[2] = byte((ip >> 16) & 0xff)
	b[3] = byte((ip >> 24) & 0xff)
	return net.IPv4(b[3], b[2], b[1], b[0]).String()
}

func readOffset(fs *os.File, from uint32)  uint32 {
	tmp := make([]byte, 3, 4)
	fs.ReadAt(tmp, int64(from))
	var offset uint32
	tmp = append(tmp, 0)
	binary.Read(bytes.NewReader(tmp), binary.LittleEndian, &offset)
	return offset
}