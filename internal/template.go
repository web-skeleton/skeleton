package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Data map[string]interface{}

func (d Data) Parse(content string) (string, error) {
	funcMap := template.FuncMap{
		"implode":     strings.Join,
		"explode":     strings.Split,
		"datetime":    datetimeFormat,
		"starts_with": startsWith,
		"ends_with":   endsWith,
		"trim":        strings.Trim,
		"trim_right":  strings.TrimRight,
		"trim_left":   strings.TrimLeft,
		"trim_space":  strings.TrimSpace,
		"format":      fmt.Sprintf,
		"integer":     toInteger,
	}

	var buffer bytes.Buffer
	if err := template.Must(template.New("").Funcs(funcMap).Parse(content)).Execute(&buffer, d); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// NewData create data object from json
func NewData(source []byte) (Data, error) {
	var data Data
	if err := json.Unmarshal(source, &data); err != nil {
		return data, err
	}

	return data, nil
}

// datetimeFormat 时间格式化
func datetimeFormat(datetime time.Time) string {
	loc, _ := time.LoadLocation("Asia/Chongqing")

	return datetime.In(loc).Format("2006-01-02 15:04:05")
}

// startsWith 判断是字符串开始
func startsWith(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.HasPrefix(haystack, n) {
			return true
		}
	}

	return false
}

// endsWith 判断字符串结尾
func endsWith(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.HasSuffix(haystack, n) {
			return true
		}
	}

	return false
}

// toInteger 转换为整数
func toInteger(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return val
}
