package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type MyReader struct {
	io.Reader
	UpdateReadSize func(n int)
}

func (r *MyReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.UpdateReadSize(n)
	return
}

type Parameters struct {
	ParallelNum int
	OutPath     string
	File        string
	FileList    string
}

func handleFlag() (*Parameters, error) {
	p := Parameters{}
	parallelNum := new(int)
	outPath := new(string)
	file := new(string)
	fileList := new(string)
	fs := flag.NewFlagSet("sfd", flag.ExitOnError)
	fs.IntVar(parallelNum, "n", runtime.NumCPU(), "Number of parallel downloads")
	fs.StringVar(outPath, "o", "", "Output path for remote file download（default current path）")
	fs.StringVar(file, "f", "", "Individual remote files to download")
	fs.StringVar(fileList, "l", "", "List of remote files to download")
	fs.Parse(os.Args[1:])
	p.ParallelNum = *parallelNum
	if len(*outPath) > 0 {
		absOutPath, _ := filepath.Abs(*outPath)
		if !isFileExist(absOutPath) {
			err := os.MkdirAll(absOutPath, 0777)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("create output path dir err: %s", err.Error()))
			}
		}
	}
	p.OutPath, _ = filepath.Abs(*outPath)
	if len(*file) > 0 {
		isValid, _ := regexp.MatchString("^(http|https)", *file)
		if !isValid {
			return nil, errors.New("invalid remote file address")
		}
		p.File = *file
	}
	if len(*fileList) > 0 {
		absFilePath, _ := filepath.Abs(*fileList)
		if !isFileExist(absFilePath) {
			return nil, errors.New("list of remote files not exist")
		}
		p.FileList = absFilePath
	}
	return &p, nil
}

func main() {
	parameters, err := handleFlag()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	parallelChannel := make(chan string, parameters.ParallelNum)
	links := make([]string, 0)
	if len(parameters.File) > 0 {
		links = append(links, parameters.File)
	}
	if len(parameters.FileList) > 0 {
		f, err := os.Open(parameters.FileList)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			body, _, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			links = append(links, string(body))
		}
	}
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWidth(80), mpb.WithWaitGroup(&wg))
	wg.Add(len(links))
	errStr := make([]string, 0)
	fmt.Printf("total download files: %d\n", len(links))
	for _, link := range links {
		parallelChannel <- link
		go func(pc chan string, url string) {
			defer wg.Done()
			err := downLoad(parameters.OutPath, url, p)
			if err != nil {
				errStr = append(errStr, fmt.Sprintf("%s: %s", url, err.Error()))
			}
			<-pc
		}(parallelChannel, link)
	}
	p.Wait()
	fmt.Printf("download result, success: %d, error: %d\n", len(links)-len(errStr), len(errStr))
	if len(errStr) > 0 {
		fmt.Println(strings.Join(errStr, "\n"))
	}
}

func isFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	return false
}

func downLoad(base, url string, p *mpb.Progress) error {
	dist := base
	idx := strings.LastIndex(url, "/")
	var name string
	if idx < 0 {
		name = url
	} else {
		name = url[idx+1:]
	}
	dist = dist + "/" + name
	v, err := http.Get(url)
	if err != nil {
		return errors.New(fmt.Sprintf("http get failed: %s", err.Error()))
	}
	defer v.Body.Close()
	if v.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("http get failed, httpcode: %d", v.StatusCode))
	}
	picFile, err := os.Create(dist)
	if err != nil {
		return errors.New(fmt.Sprintf("create [%s] failed: %s", dist, err.Error()))
	}
	defer picFile.Close()
	bar := p.AddBar(v.ContentLength,
		mpb.PrependDecorators(
			decor.Name(name),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
		),
	)
	myReader := &MyReader{
		Reader: v.Body,
		UpdateReadSize: func(n int) {
			bar.IncrBy(n)
		},
	}
	_, err = io.Copy(picFile, myReader)
	if err != nil {
		return errors.New(fmt.Sprintf("save file failed: %s", err.Error()))
	}
	return nil
}
