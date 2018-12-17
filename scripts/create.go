package main

import (
	"log"
	"time"
	"flag"
	"math/rand"
	"fmt"
	"bytes"
	"bufio"
	"os"
)



var msgSize int
var processSize int
var writePath string
var isFsync bool



func writer(path string, msg []byte) {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000)))
	f, err := os.Create(path)
	for {
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = f.Write(msg)
		if err != nil {
			log.Fatal(err)
			return
		}
		if isFsync {
			err = f.Sync()
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		time.Sleep(time.Second)
	}
	f.Close()
}

func init() {
	flag.IntVar(&msgSize, "msg-size", 100, "Message size per second")
	flag.IntVar(&processSize, "process-size", 100, "How much process per second")
	flag.BoolVar(&isFsync, "is-fsync", false, "make sure file is sync")
	flag.StringVar(&writePath, "write-path", "/tmp", "Write file to a path")
}

func main() {

	flag.Parse()
	fmt.Printf("path:%s ,%d process * %d bytes / second \n",writePath, processSize, msgSize)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	for i := 1; i < msgSize; i++ {
		w.Write([]byte("0"))
	}
	w.Write([]byte("\n"))
	w.Flush()
	msg := b.Bytes()


	for i := 1; i <= processSize; i++ {
		go writer(fmt.Sprintf( writePath+ "/%d.log", i), msg)
	}
	for  {
		time.Sleep(time.Second)
	}
}