package hermes

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type HermesSession struct {
	HermesFilePool   FilePool
	HermesRPCService RPCService
	HermesOutWriter  OutWriterGenerator
}

func NewHermesSession() *HermesSession {
	return &HermesSession{
		HermesFilePool:   FilePool{},
		HermesRPCService: RPCService{},
		HermesOutWriter:  DefaultFoutGenerator,
	}
}

func (hs *HermesSession) Close() {
	hs.HermesFilePool.Close()
}

// FileDescriptior describes a file with name, usage description and debug log output channel
type FileDescriptior struct {
	FilePath        string        // absolute file path
	FileDescription string        // file useage/description (optional)
	UseFilePool     bool          // store file content, load file only once (optional, default false)
	debugOut        chan<- string // debug output channel for concurrent execution (optional, default nil)
	logID           string        // log ID to identify run (optional)
	ContinueOnError bool          // terminate the programm on error or continue (optional, default false)
}

// Open will Open a file and create a scanner, to iterate through the file
func (s *HermesSession) Open(fd *FileDescriptior) (*os.File, *bufio.Scanner, error) {

	// remove space characters
	fileCorrected := strings.TrimSpace(fd.FilePath)
	fd.FilePath = fileCorrected

	if fd.UseFilePool {
		byteData := s.HermesFilePool.Get(fd)
		r := bytes.NewReader(byteData)
		scanner := bufio.NewScanner(r)
		return nil, scanner, nil
	}

	file, err := os.Open(fd.FilePath)
	if err != nil {
		if fd.ContinueOnError {
			if fd.debugOut != nil {
				fd.debugOut <- fmt.Sprintf("%s Error occured while reading %s: %s   \n", fd.logID, fd.FileDescription, fd.FilePath)
			} else {
				fmt.Printf("Error occured while reading %s: %s   \n", fd.FileDescription, fd.FilePath)
			}

			return nil, nil, err
		}
		log.Fatalf("Error occured while reading %s: %s   \n", fd.FileDescription, fd.FilePath)
	}
	// if fd.debugOut != nil {
	// 	fd.debugOut <- fmt.Sprintf("%s %s", fd.logID, fd.filePath)
	// }
	scanner := bufio.NewScanner(file)
	return file, scanner, nil
}

func (s *HermesSession) ReadFile(fd *FileDescriptior) ([]byte, error) {

	// remove space characters
	fileCorrected := strings.TrimSpace(fd.FilePath)
	fd.FilePath = fileCorrected

	if fd.UseFilePool {
		byteData := s.HermesFilePool.Get(fd)
		return byteData, nil
	}

	byteData, err := os.ReadFile(fd.FilePath)
	if err != nil {
		if fd.ContinueOnError {
			if fd.debugOut != nil {
				fd.debugOut <- fmt.Sprintf("%s Error occured while reading %s: %s   \n", fd.logID, fd.FileDescription, fd.FilePath)
			} else {
				fmt.Printf("Error occured while reading %s: %s   \n", fd.FileDescription, fd.FilePath)
			}

			return nil, err
		}
		log.Fatalf("Error occured while reading %s: %s   \n", fd.FileDescription, fd.FilePath)
	}
	return byteData, nil
}

// OpenResultFile opens a file for writing - options append, if append is false, it will truncate the file and override
func (s *HermesSession) OpenResultFile(filePath string, append bool) OutWriter {

	if s.HermesOutWriter == nil {
		s.HermesOutWriter = DefaultFoutGenerator
	}
	res, err := s.HermesOutWriter(filePath, append)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

// WriteYamlConfig write a default config file
func (s *HermesSession) WriteYamlConfig(filename string, structIn interface{}) {
	file := s.OpenResultFile(filename, false)
	defer file.Close()
	data, err := yaml.Marshal(structIn)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := file.WriteBytes(data); err != nil {
		log.Fatal(err)
	}
}

// DumpStructToFile debug dump global variables to a file
func (s *HermesSession) DumpStructToFile(filename string, global interface{}) {

	file := s.OpenResultFile(filename, false)
	defer file.Close()

	data, err := yaml.Marshal(global)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := file.WriteBytes(data); err != nil {
		log.Fatal(err)
	}
}
