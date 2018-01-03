package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"os/exec"
	"strings"
	"time"
)

func init() {
	flag.StringVar(&flagvar, "f", "test", "help message for flagname")
	flag.Parse()

}

var flagvar string

func main() {
	fmt.Printf("Flag: %v\n", flagvar)

	switch flag.Arg(0) {
	case "":
		Err("ARGERR", flag.ErrHelp, 1)
		break
	case "-":
		fmt.Println("stdin")
		os.Exit(0)
		break

	default: // Try to run jobs
		fmt.Printf("Try: %v\n", strings.Join( flag.Args(),";"))
		cwd := Workdir()
		fmt.Printf("CWD: %s\n", cwd)
		sources := Enumerate(cwd)
		fmt.Printf("Files: \n\t%s\n", strings.Join( sources,"\n\t"))
		makefile := Compile(sources)
		tempfile := WriteTemp(makefile)
		// defer RemoveFile(tempfile)
		RunMake(tempfile,  flag.Args())
		os.Exit(0)
		break
	}
}

// Enumerate

func Enumerate(path string) (list []string) {
	list = append(list, ReadDirOrFail(path)...)

	list = append(list, ReadDirOrFail(os.Getenv("HOME"))...)
	bakepath := os.Getenv("BAKEPATH")
	
	if bakepath == "" { 
		return list
	}

	for _, bakedir := range strings.Split(bakepath, ";") {
		list = append(list, ReadDirOrFail(bakedir)...)
	}

	return list
}

func ReadDirOrFail(path string) (list []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		Err("DIRERR", err, 1)
	}

	for _, file := range files {
		if name := file.Name(); isBakefile(name) {
			list = append(list, name)
		} else {
			// fmt.Println(name)
		}
	}

	return list
}

func onlyBakefiles(files []string) (match []string) {
	for _, file := range files {
		match = append(match, file)
	}
	return match
}

func isBakefile(file string) bool {
	ext := filepath.Ext(file)
	// fmt.Println(ext)
	return ext == ".bake" || ext == ".Bake" || ext == ".BAKE"
}

// Compile

func Compile(sources []string) string {
	contents := []string{}
	for _, source := range sources {
		if content, err := ioutil.ReadFile(source); err == nil {
			contents = append(contents, Parse(content))
		} else {
			Err("READERR", err, 1)
		}
	}
	return strings.Join(contents, "\n")
}

func Parse(src []byte) string {
	return string(src)
}

// WriteTemp

func WriteTemp(src string) string {
	file := fmt.Sprintf("/tmp/bakejob_%d", time.Now().Unix())

	err := ioutil.WriteFile(file, []byte(src), 777)
	if err != nil {
		Err("WRITERR", err, 1)
	}
	return file
}

// RunMake

func RunMake(file string, jobs []string) []byte {
	args := append([]string{"-f", file}, jobs...)

	fmt.Printf("command is::\n%s %s\n", "/usr/bin/make", strings.Join(args," "))

	out, err  := exec.Command ("/usr/bin/make", args...).Output()
	if err!=nil {
		fmt.Printf("output is::\n%s\n", out)
		Err("EXECERR", err, 1)
	}
	fmt.Printf("output is::\n%s\n", out)
	return out
}

// Util

func Err(msg string, err error, code int) {
	fmt.Println(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(code)
}

func Workdir() string {
	cwd, err := os.Getwd()
	if err != nil {
		Err("CWDERR", err, 1)
	}
	return cwd
}

func RemoveFile(file string) {
	if err:=os.Remove(file); err != nil {
		Err("RMERR",err,1)
	}
}