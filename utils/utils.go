package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const HourKeyFactor = 10
const PartitionCountOfHourKey = 10
const DayKeyFactor = 100
const PartitionCountOfDayKey = 10

func GetDayPartitionKeyByUuid(iDay int, uuid int64) int {
	return GetDayPartitionKeyByLong(iDay, uuid)
}
func GetDayPartitionKeyByMac(iDay int, apMac int64) int {
	return GetDayPartitionKeyByLong(iDay, apMac)
}
func GetHourPartitionKeyByUuid(iHour int, uuid int64) int {
	return GetHourPartitionKeyByLong(iHour, uuid)
}
func GetDayPartitionKeyByLong(iDay int, l int64) int {
	// use the last byte of apMac to compose Partitionkey
	key := (iDay * DayKeyFactor) + int((l>>40)%PartitionCountOfDayKey)
	return key
}

/**
历史原因，小时的分片位数只有一位
*/
func GetHourPartitionKeyByLong(iHour int, l int64) int {
	// use the last byte of apMac to compose Partitionkey
	key := (iHour * HourKeyFactor) + int((l>>40)%PartitionCountOfHourKey)
	return key
}

func GetIntDay(time time.Time) int {
	return time.Year()*10000 + int(time.Month())*100 + time.Day()
}
func GetIntHour(time time.Time) int {
	return time.Year()*1000000 + int(time.Month())*10000 + time.Day()*100 + time.Hour()
}
func ConvertMacLongToStr(lMac int64) string {
	buffer := &bytes.Buffer{}
	for i := 0; i < 6; i++ {
		AppendByte(buffer, byte(lMac%0x100))
		lMac /= 0x100
		if i < 5 {
			buffer.WriteString(":")
		}
	}
	return buffer.String()
}
func AppendByte(buffer *bytes.Buffer, b byte) {
	str := strconv.FormatInt(int64(b), 16)
	if len(str) < 2 {
		buffer.WriteString("0")
	}
	buffer.WriteString(str)
}
func ConvertMacStrToLong(macStr string) (int64, error) {
	bytes, err := ConvertToBytes(macStr)
	if err != nil {
		return 0, err
	}
	lMac := ConvertMacBytesToLong(bytes, 0)
	return lMac, nil
}
func ConvertMacBytesToLong(bytes []byte, start int) int64 {
	var lMac int64
	lMac = int64(bytes[start+5])*0x10000000000 + int64(bytes[start+4])*0x100000000 + int64(bytes[start+3])*0x1000000 + int64(bytes[start+2])*0x10000 + int64(bytes[start+1])*0x100 + int64(bytes[start])
	//lMac = 1;
	return lMac
}
func ConvertMacBytesToString(mac []byte) string {
	return ConvertMacBytesToStringPart(mac, 0)
}
func ConvertMacBytesToStringPart(mac []byte, start int) string {
	var buffer bytes.Buffer
	for i := start; i < start+6; i++ {
		byteStr := strconv.FormatInt(int64(mac[i]), 16)
		if len(byteStr) < 2 {
			buffer.WriteString("0")
		}
		buffer.WriteString(byteStr)
		if i < start+5 {
			buffer.WriteString(":")
		}
	}
	return buffer.String()
}
func ConvertToBytes(mac string) ([]byte, error) {
	byteStrs := strings.Split(mac, ":")
	byteLen := len(byteStrs)
	if byteLen != 6 {
		return nil, errors.New("Invalid Mac: " + mac)
	}
	bytes := make([]byte, 6)
	for i := 0; i < byteLen; i++ {
		// use 16 cause sometimes 8 will overflow
		b, err := strconv.ParseInt(byteStrs[i], 16, 16)
		if err != nil {
			return nil, err
		}
		bytes[i] = byte(b)
	}
	return bytes, nil
}

func BuildInClause(vals []string) string {
	inClause := "("
	for i := 0; i < len(vals); i++ {
		inClause += "'" + vals[i] + "'"
		if i != len(vals)-1 {
			inClause += ","
		}
	}
	inClause += ")"
	return inClause
}

/** Belows are line format message**/

func GetNextString(data []byte, start int) (string, int) {
	nextPos := FindNextLinePos(data, start)
	strData := data[start:nextPos]
	return string(strData), nextPos
}
func FindNextMessagePos(data []byte, start int) int {
	length := len(data)
	for i := start; i < length; i++ {
		if data[i] == '#' {
			return i
		}
	}
	return -1
}
func FindNthLinePos(data []byte, start int, lineNum int) int {
	currLineNum := 0
	length := len(data)
	for i := start; i < length; i++ {
		if data[i] == '\n' {
			currLineNum++
			if currLineNum == lineNum {
				return i
			}
		}
	}
	return -1
}
func FindNextLinePos(data []byte, start int) int {
	length := len(data)
	for i := start; i < length; i++ {
		if data[i] == '\n' {
			return i
		}
	}
	return -1
}

func GetCurrentExeDir() string {
	dirPath := ""
	file, _ := exec.LookPath(os.Args[0])
	filePath, _ := filepath.Abs(file)
	pathSeparator := string(os.PathSeparator)
	if runtime.GOOS == "windows" {
		lastPos := strings.LastIndex(filePath, pathSeparator)
		if lastPos >= 0 {
			dirPath = Substr(filePath, 0, lastPos)
		}
	} else {
		dirPath = path.Dir(filePath)
	}
	return dirPath
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func GenerateUuid() int64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Int63()
}

func GetUnixTimeStampMs(time time.Time) string {
	return strconv.FormatInt(time.UnixNano()/1000000, 10)
}

func isAnyEmpty(arr ...string) bool {
	for _, item := range arr {
		if len(item) == 0 {
			return true
		}
	}

	return false
}

func encryptMD5(input string) string {
	m := md5.New()
	m.Write([]byte(input))
	data := m.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(data))
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
