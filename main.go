package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
)

func dirTree(writer io.Writer, path string, flag bool) error {
	err := searchDir(writer, path, "", flag)
	if err != nil {
		return err
	}
	return nil
}

func getPrefix(isLast bool) string {
	if isLast {
		return "└───"
	}
	return "├───"
}

func getNewPrefix(prefix string, isLast bool) string {
	if isLast {
		return prefix + "\t"
	}
	return prefix + "│\t"
}

func getFileSize(info fs.FileInfo) string {
	size := info.Size()

	if size == 0 {
		return " (empty)"
	}
	return " (" + strconv.FormatInt(size, 10) + "b)"
}

func searchDir(writer io.Writer, path, prefix string, flag bool) error {
	startFiles, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var files []os.DirEntry
	for _, file := range startFiles {
		if flag || file.IsDir() {
			files = append(files, file)
		}
	}

	for i, f := range files {
		isLast := i == len(files)-1

		filePrefix := prefix + getPrefix(isLast)

		if f.IsDir() {
			fmt.Fprintln(writer, filePrefix+f.Name())

			dirPrefix := getNewPrefix(prefix, isLast)
			err = searchDir(writer, path+"/"+f.Name(), dirPrefix, flag)
			if err != nil {
				return err
			}
		} else {
			info, err := f.Info()
			if err != nil {
				return err
			}

			sizeStr := getFileSize(info)
			fmt.Fprintln(writer, filePrefix+f.Name()+sizeStr)
		}
	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
