package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

var workqueue []string
var procqueue chan int
var itemCounter = 0
var waitCounter sync.WaitGroup
var totalcount = 0
var vmode *bool
var version = "v1.0"

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

//using pure go minify
func doWork(filename string) error {
	defer waitCounter.Done()

	if *vmode {
		fmt.Printf("[%d/%d] - %s\n", itemCounter, totalcount, filename)
	}
	fi, err := os.Open(filename)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, fi)
	if err != nil {
		return err
	}
	fi.Close()

	// make a read buffer
	r := bufio.NewReader(buf)

	// open output file
	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fo.Close()

	// make a write buffer
	w := bufio.NewWriter(fo)
	//

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepWhitespace:      false,
	})
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	err = m.Minify("text/javascript", w, r)
	w.Flush()
	itemCounter++
	procqueue <- itemCounter
	return err
}

func visit(fileloc string, f os.FileInfo, err error) error {
	//make sure it is a file
	if f.Mode().IsRegular() == true {
		ext := path.Ext(fileloc)

		//make sure we are dealing with javascript files (.js)
		if strings.Compare(strings.ToLower(ext), string(".js")) == 0 {
			workqueue = append(workqueue, fileloc)
		}
	}
	return nil
}

func processWorkQueue() {
	fmt.Printf("%v - Processing %d items.\n", time.Now(), len(workqueue))
	for item := range workqueue {
		//wait here for your turn
		<-procqueue
		go doWork(workqueue[item])
	}
}

func usage() {
	fmt.Println("go-jsminify", version)
	fmt.Println("usage:")
	fmt.Println("\tgo-jsminify -s=<sourcefolder> -w=<workercount> -v")
	fmt.Println("\tgo-jsminify c:\\temp\\jsfiles 4 -v")
}

func main() {
	//debug flag
	vmode = flag.Bool("v", false, "turn on verbose.")

	proccount := flag.Int("w", 1, "worker count")
	sourcedir := flag.String("s", "", "sourcefolder")
	flag.Parse()

	//we need at least the source code directory
	if len(*sourcedir) == 0 {
		usage()
		os.Exit(1)
	}

	t1 := time.Now()
	procqueue = make(chan int, *proccount)

	err := filepath.Walk(*sourcedir, visit)
	totalcount = len(workqueue)
	waitCounter.Add(totalcount)
	for i := 0; i < *proccount; i++ {
		procqueue <- i
	}

	processWorkQueue()
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	waitCounter.Wait()
	close(procqueue)
	elapsed := time.Since(t1)
	fmt.Printf("%v - Processed %d files in %v time.", time.Now(), totalcount, elapsed)
}
