package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

type FlagsX struct {
	Path string `long:"path" description:"Path to files. Example: --path=files/*.jpg"`
	Help bool   `short:"h" long:"help" description:"Show this help message"`
}

var (
	Flags = &FlagsX{}
)

func main() {
	// Flags
	parser := flags.NewParser(Flags, flags.PrintErrors)
	parser.LongDescription = `
===================================
MD5 hash file changer written in Go
https://github.com/lexesv/md5-file-hash-changer
===================================`
	parser.Usage = "[OPTIONS]"
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		LogFatal(err)
	}
	if Flags.Path == "" || Flags.Help {
		var b bytes.Buffer
		parser.WriteHelp(&b)
		fmt.Println(b.String())
	}
	files, err := filepath.Glob(Flags.Path)
	if err != nil {
		LogFatal(err)
	}
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			LogFatal(err)
		}
		fd, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, info.Mode())
		if err != nil {
			LogFatal(err)
		}
		md5_old, err := MD5(fd);
		if err != nil {
			LogFatal(err)
		}
		b := make([]byte, 1)
		rand.Seed(time.Now().UTC().UnixNano())
		rand.Read(b)
		fd.Write(b)
		fd.Close()
		fd, _ = os.OpenFile(file, os.O_RDONLY, info.Mode())
		defer fd.Close()
		if md5_new, err := MD5(fd); err != nil {
			LogFatal(err)
		} else {
			Log("File:", file, md5_old, "->", md5_new)
		}
	}

}

func Log(s ...interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	s = preparePrintLog(fn, line, s...)
	log.Println(s...)
}

func LogFatal(e ...interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	e = preparePrintLog(fn, line, e...)
	log.Println(e)
	os.Exit(1)

}
func preparePrintLog(file string, line int, s ...interface{}) []interface{} {
	file = strings.Split(filepath.Base(file), ".")[0]
	return append([]interface{}{file + ":" + strconv.Itoa(line)}, s...)
}

func MD5(fd *os.File) (s string, err error) {
	hash := md5.New()
	if _, err := io.Copy(hash, fd); err != nil {
		return s, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}

func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
