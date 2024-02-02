package hermes

import (
	"bufio"
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// LineInut read next line from bufio.Scanner
func LineInut(scanner *bufio.Scanner) (line string) {
	if ok := scanner.Scan(); ok {
		line = scanner.Text()
	} else if err := scanner.Err(); err != nil {
		log.Fatalf("ERROR: Failed to read file: %v", err)
	} else {
		log.Fatal("ERROR: EOF")
	}
	return line
}

func NextLineInut(idIdx int, scanner *bufio.Scanner, splitFunc func(string) []string) (id string, tokens []string, valid bool) {
	valid = false
	if ok := scanner.Scan(); ok {
		line := scanner.Text()
		tokens = splitFunc(line)
		if len(tokens) > idIdx {
			id = strings.TrimSpace(tokens[idIdx])
			valid = true
		}
	}
	return id, tokens, valid
}

// ValAsFloat parse text into float
func ValAsFloat(toParse, filename, line string) float64 {
	noSpaces := strings.TrimSpace(toParse)
	value, err := strconv.ParseFloat(noSpaces, 64)
	if err != nil {
		log.Fatalf("Error: parsing float! File: %s \n   Line: %s \n", filename, line)
	}
	return value
}

// TryValAsFloat parse text into float
func TryValAsFloat(toParse string) (float64, error) {
	noSpaces := strings.TrimSpace(toParse)
	value, err := strconv.ParseFloat(noSpaces, 64)
	return value, err
}

// ValAsBool parse text into bool
func ValAsBool(toParse, filename, line string) bool {
	noSpaces := strings.TrimSpace(toParse)
	value, err := strconv.ParseInt(noSpaces, 10, 64)
	if err != nil {
		log.Fatalf("Error: parsing int! File: %s \n   Line: %s \n", filename, line)
	}
	if value == 0 {
		return false
	} else if value == 1 {
		return true
	} else {
		log.Fatalf("Error: parsing int as bool! File: %s \n   Line: %s \n", filename, line)
	}
	return false
}

// ValAsInt parse text into int
func ValAsInt(toParse, filename, line string) int64 {
	noSpaces := strings.TrimSpace(toParse)
	value, err := strconv.ParseInt(noSpaces, 10, 64)
	if err != nil {
		log.Fatalf("Error: parsing int! File: %s \n   Line: %s \n", filename, line)
	}
	return value
}

type seperatorRunes struct {
	Seperator []rune
}

func (sR *seperatorRunes) isRune(r rune) bool {
	for _, sep := range sR.Seperator {
		if r == sep {
			return true
		}
	}
	return false
}

// Explode splits a string by a set of seperator runes
func Explode(str string, seperator []rune) (result []string) {
	var sR seperatorRunes
	sR.Seperator = seperator

	result = strings.FieldsFunc(str, sR.isRune)
	return result
}

// DateConverterFunc type of DateConverter
type DateConverterFunc func(string) (int, int)

// DateConverter get a function that calculates (ztDat = day of the year) and (masDat = total days since 01.01.1901)
func DateConverter(splitCenturyAt int, dateformat DateFormat) func(string) (ztDat, masDat int) {

	format := dateformat
	cent := splitCenturyAt
	return func(ztdatIn string) (ztDat, masDat int) {
		MT := [12]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}

		ztdatInNoSpaces := strings.TrimSpace(ztdatIn)

		var TG int
		var MON int
		var YR int
		var err error
		switch format {
		case DateDEshort:
			TG, MON, YR, err = extractDate(ztdatInNoSpaces, true)
			if YR < cent {
				YR = YR + 100
			}
		case DateDElong:
			TG, MON, YR, err = extractDate(ztdatInNoSpaces, false)
			if YR < 1901 {
				log.Fatalf("Error: parsing date! Date before 1901 are not supported: %s \n", ztdatInNoSpaces)
			}
			YR = YR - 1900
		case DateENshort:
			MON, TG, YR, err = extractDate(ztdatInNoSpaces, true)
			if YR < cent {
				YR = YR + 100
			}
		case DateENlong:
			MON, TG, YR, err = extractDate(ztdatInNoSpaces, false)
			if YR < 1901 {
				log.Fatalf("Error: parsing date! Date before 1901 are not supported: %s \n", ztdatInNoSpaces)
			}
			YR = YR - 1900
		}
		if err != nil {
			log.Fatal(err)
		}

		// add a day if the year is a leap year
		if YR%4 == 0 {
			for i := 1; i <= 12; i++ {
				if i >= 3 {
					MT[i-1] = MT[i-1] + 1
				}
			}
		}

		// number of days starting 01.01.1901
		masDat = (YR-1)*365 + (YR-1)/4 + MT[MON-1] + TG
		ztDat = MT[MON-1] + TG
		return ztDat, masDat
	}
}

func extractDate(date string, short bool) (first, second, third int, err error) {
	if short {
		if len(date) == 6 {
			first = int(ValAsInt(date[0:2], "none", date))
			second = int(ValAsInt(date[2:4], "none", date))
			third = int(ValAsInt(date[4:6], "none", date))
		} else if len(date) == 8 {
			first = int(ValAsInt(date[0:2], "none", date))
			second = int(ValAsInt(date[3:5], "none", date))
			third = int(ValAsInt(date[6:8], "none", date))
		} else {
			return first, second, third, errors.New("wrong date format for short year format")
		}

	} else {
		if len(date) == 8 {
			first = int(ValAsInt(date[0:2], "none", date))
			second = int(ValAsInt(date[2:4], "none", date))
			third = int(ValAsInt(date[4:8], "none", date))
		} else if len(date) == 10 {
			first = int(ValAsInt(date[0:2], "none", date))
			second = int(ValAsInt(date[3:5], "none", date))
			third = int(ValAsInt(date[6:10], "none", date))
		} else {
			return first, second, third, errors.New("wrong date format for long year format")
		}
	}
	return first, second, third, nil
}

// datumOld deprecated calculates (ztDat = day of the year) and (masDat = total days since 1900)
func datumOld(ztdatIn string, cent int) (ztDat, masDat int) {
	// !                        BERECHNUNG DES DATUMS
	MT := [12]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
	TG := int(ValAsInt(ztdatIn[0:2], "none", ztdatIn))
	MON := int(ValAsInt(ztdatIn[2:4], "none", ztdatIn))
	YR := int(ValAsInt(ztdatIn[4:6], "none", ztdatIn))
	if YR < cent {
		YR = YR + 100
	}

	// add a day if the year is a leap year
	if YR%4 == 0 {
		for i := 1; i <= 12; i++ {
			if i >= 3 {
				MT[i-1] = MT[i-1] + 1
			}
		}
	}

	// number of days starting 1901
	masDat = (YR-1)*365 + (YR-1)/4 + MT[MON-1] + TG
	ztDat = MT[MON-1] + TG
	return ztDat, masDat
}

// KalenderConverterFunc KalenderConverter type
type KalenderConverterFunc func(int) string

// KalenderConverter function to convert internal date format into string
func KalenderConverter(dateformat DateFormat, seperator string) func(int) string {
	format := dateformat
	formatStrShort := "%02d" + seperator + "%02d" + seperator + "%02d"
	formatStrLong := "%02d" + seperator + "%02d" + seperator + "%d"

	return func(MASDAT int) (KALDAT string) {

		year, month, day := KalenderDate(MASDAT)
		YR := year - 1900
		switch format {
		case DateDElong:
			KALDAT = fmt.Sprintf(formatStrLong, day, month, year)
		case DateDEshort:
			if YR > 99 {
				KALDAT = fmt.Sprintf(formatStrShort, day, month, YR-100)
			} else {
				KALDAT = fmt.Sprintf(formatStrShort, day, month, YR)
			}
		case DateENlong:
			KALDAT = fmt.Sprintf(formatStrLong, month, day, year)
		case DateENshort:
			if YR > 99 {
				KALDAT = fmt.Sprintf(formatStrShort, month, day, YR-100)
			} else {
				KALDAT = fmt.Sprintf(formatStrShort, month, day, YR)
			}
		}
		return KALDAT
	}
}

// KalenderDate get year, month, day from MASDAT( = total days since 1900)
func KalenderDate(MASDAT int) (year, month, day int) {
	MT := []int{31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365}
	YR := MASDAT / 365
	if MASDAT%365 <= YR/4 {
		YR = YR - 1
	}
	TG := MASDAT - YR*365 - YR/4
	KORR := 0
	if (YR+1)%4 == 0 {
		if TG > 59 {
			KORR = 1
		} else {
			KORR = 0
		}
	} else {
		KORR = 0
	}
	MOZ := 1
	for {
		MOZindex := MOZ - 1
		if MOZ > 1 {
			MT[MOZindex] = MT[MOZindex] + KORR
		}
		if TG <= MT[MOZindex] {
			break
		} else {
			MOZ++
		}
	}
	if MOZ > 1 {
		TG = TG - MT[MOZ-2]
	}
	year = YR + 1900 + 1
	month = MOZ
	day = TG

	return year, month, day
}

// func leftAlignmentFormat(numberOfRunes int, input string) (outStr string) {
// 	lenStr := utf8.RuneCountInString(input)
// 	var builder strings.Builder
// 	if numberOfRunes > lenStr {
// 		numSpaces := numberOfRunes - lenStr
// 		builder.Grow(len(input) + numSpaces)
// 		builder.WriteString(input)
// 		for runeIndex := 0; runeIndex < numSpaces; runeIndex++ {
// 			builder.WriteByte(' ')
// 		}
// 		outStr = builder.String()
// 	} else {
// 		outStr = input
// 	}
// 	return outStr
// }

// AskDirectory returns the current directory
func AskDirectory() string {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return root
}

// create directory if file does not exist
func MakeDir(outPath string) {
	dir := filepath.Dir(outPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("ERROR: Failed to generate output path %s :%v", dir, err)
		}
	}
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
func Open(fd *FileDescriptior) (*os.File, *bufio.Scanner, error) {

	// remove space characters
	fileCorrected := strings.TrimSpace(fd.FilePath)
	fd.FilePath = fileCorrected

	if fd.UseFilePool {
		byteData := HermesFilePool.Get(fd)
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

func ReadFile(fd *FileDescriptior) ([]byte, error) {

	// remove space characters
	fileCorrected := strings.TrimSpace(fd.FilePath)
	fd.FilePath = fileCorrected

	if fd.UseFilePool {
		byteData := HermesFilePool.Get(fd)
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
func OpenResultFile(filePath string, append bool) *Fout {
	var flags int
	if append {
		flags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	} else {
		flags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	}

	file, err := os.OpenFile(filePath, flags, 0600)
	if err != nil {
		log.Fatalf("Error occured while opening result file: %s   \n", filePath)
	}
	fwriter := bufio.NewWriter(file)
	return &Fout{file, fwriter}
}

// PrintTo prints a string into a file
func PrintTo(file *Fout, text string) {
	if _, err := file.Write(text); err != nil {
		log.Fatal(err)
	}
}

// min minimum of 2 int values
func min(val1, val2 int) int {
	if val1 < val2 {
		return val1
	}
	return val2
}

//type printToLimitFunc func(int)

// printToLimit provides println debug output for the first n occurences
func (g *GlobalVarsMain) printToLimit(n int) func(int, interface{}) {
	limit := n
	counter := 0
	return func(zeit int, val interface{}) {
		if counter <= limit {
			date := g.Kalender(zeit)
			fmt.Println(date, val)
			counter++
		}
	}
}

// print out error message to chanel or fail Fatal
func printError(logID, errorMsg string, out, logout chan<- string) {
	if logout != nil {
		logout <- fmt.Sprintf("%s Error: %s", logID, errorMsg)
	}
	if out != nil {
		out <- fmt.Sprintf("%s Error: %s", logID, errorMsg)
	} else {
		log.Fatal(errorMsg)
	}
}

// DumpStructToFile debug dump global variables to a file
func DumpStructToFile(filename string, global *GlobalVarsMain) {

	file := OpenResultFile(filename, false)
	defer file.Close()

	data, err := yaml.Marshal(global)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := file.WriteBytes(data); err != nil {
		log.Fatal(err)
	}
}
func isNil(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
func toYamlNode(in interface{}) (*yaml.Node, error) {
	node := &yaml.Node{}
	// do not wrap yaml.Node into yaml.Node
	if n, ok := in.(*yaml.Node); ok {
		return n, nil
	}

	// if input implements yaml.Marshaler we should use that marshaller instead
	// same way as regular yaml marshal does
	if m, ok := in.(yaml.Marshaler); ok && !isNil(reflect.ValueOf(in)) {
		res, err := m.MarshalYAML()
		if err != nil {
			return nil, err
		}

		if n, ok := res.(*yaml.Node); ok {
			return n, nil
		}

		in = res
	}
	if _, ok := in.(encoding.TextMarshaler); ok && !isNil(reflect.ValueOf(in)) {
		return node, node.Encode(in)
	}

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		node.Kind = yaml.MappingNode
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			// skip unexported fields
			if !v.Field(i).CanInterface() {
				continue
			}

			tag := t.Field(i).Tag.Get("yaml")
			parts := strings.Split(tag, ",")
			fieldName := parts[0]
			parts = parts[1:]

			if fieldName == "" {
				fieldName = strings.ToLower(t.Field(i).Name)
			}

			if fieldName == "-" {
				continue
			}
			// handle omitempty and omitonlyifnil from yml tag
			var (
				empty = isEmpty(v.Field(i))
				null  = isNil(v.Field(i))

				skip bool
			)
			for _, part := range parts {
				if part == "omitempty" && empty {
					skip = true
				}

				if part == "omitonlyifnil" && !null {
					skip = false
				}
			}
			if skip {
				continue
			}
			// just utilize the head comment
			headComment := t.Field(i).Tag.Get("comment")

			childKey, err := toYamlNode(fieldName)
			if err != nil {
				return nil, err
			}
			childKey.HeadComment = headComment

			var value interface{}
			if v.Field(i).CanInterface() {
				value = v.Field(i).Interface()
			}

			childValue, err := toYamlNode(value)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, childKey, childValue)
		}
	case reflect.Slice:
		node.Kind = yaml.SequenceNode
		nodes := make([]*yaml.Node, v.Len())

		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)

			var err error

			nodes[i], err = toYamlNode(element.Interface())
			if err != nil {
				return nil, err
			}
		}
		node.Content = append(node.Content, nodes...)
	default:
		if err := node.Encode(in); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func isEmpty(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.Ptr:
		return value.IsNil()
	case reflect.Map:
		return len(value.MapKeys()) == 0
	case reflect.Slice:
		return value.Len() == 0
	default:
		return value.IsZero()
	}
}
