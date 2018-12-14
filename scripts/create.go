package main

import (
	"os"
	"log"
	"time"
	"bufio"
	"bytes"
	"fmt"
)






func writer(path string, msg []byte) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	for {
		_, err = f.Write(msg)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = f.Sync()
		if err != nil {
			log.Fatal(err)
			return
		}
		time.Sleep(time.Second)
	}


	defer f.Close()
}

func main() {


	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	for i := 1; i < 1000; i++ {
		w.Write([]byte("0"))
	}
	w.Write([]byte("\n"))
	w.Flush()
	msg := b.Bytes()


	for i := 1; i <= 3000; i++ {
		go writer(fmt.Sprintf("/tmp/watch_test/test/%d.log", i), msg)
	}
	for  {
		time.Sleep(time.Second)
	}
}