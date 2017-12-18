package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var (
	codeMap map[string]int
)

const (
	filePath = "test.txt"
	symbols  = "+-*/:=<>;()$"
)

//Unit : 存放关键字、数字、标识符、运算符等的结构体
type Unit struct {
	length  int    //内容长度
	typenum int    //种别码
	value   string //内容
}

//RDP : recursive-descent parsing
//		递归下降分析
type RDP struct {
	data    []byte //扫描的程序
	index   int    //指针
	unit    Unit   //当前扫描的单词
	isValid bool   //程序是否合法
	tag     string //出错信息
}

func (u Unit) String() string {
	var res string
	switch u.typenum {
	case 2:
		res = "(" + strconv.Itoa(u.typenum) + ", '" + u.value + "')"
	default:
		res = "(" + strconv.Itoa(u.typenum) + ", " + u.value + ")"
	}
	return res
}

func main() {
	//file, err := os.Open(filePath)
	//checkError(err)
	//defer file.Close()
	//data, err := ioutil.ReadAll(file)
	//checkError(err)
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	str = strings.Trim(str, "\r\n")
	data := []byte(str)
	rdp := RDP{data: data}
	rdp.parse()
	if rdp.isValid {
		fmt.Println("success")
	} else {
		fmt.Println("error")
		fmt.Println(rdp.tag)
	}
}

func init() {
	codeMap = map[string]int{
		"$":     -1,
		"begin": 0,
		"end":   1,
		"ID":    2,
		"NUM":   3,
		"SPACE": 4,
		"if":    10,
		"then":  11,
		"while": 12,
		"do":    13,
		"+":     20,
		"-":     21,
		"*":     22,
		"/":     23,
		":":     24,
		":=":    25,
		"<":     26,
		"<>":    27,
		"<=":    28,
		">":     29,
		">=":    30,
		"=":     31,
		";":     32,
		"(":     33,
		")":     34,
		"ERR":   99,
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func scanner(data []byte, index int) Unit {
	if index >= len(data) {
		return Unit{}
	}
	var length, code int
	var value string
	switch {
	case isLetter(data[index]):
		length, value = getWord(data, index)
		code = checkID(value)
	case isDigit(data[index]):
		length, value = getNUM(data, index)
		code = codeMap["NUM"]
	case isSymbol(data[index]):
		length, value = getOperator(data, index)
		code = checkOperator(value)
	case isSpace(data[index]):
		length = 1
		code = codeMap["SPACE"]
		value = ""
	default:
		length = 1
		code = codeMap["ERR"]
		value = string(data[index])
	}
	return Unit{
		length:  length,
		typenum: code,
		value:   value,
	}
}

func isLetter(key byte) bool {
	return unicode.IsLetter(rune(key))
}

func isDigit(key byte) bool {
	return unicode.IsDigit(rune(key))
}

func isSpace(key byte) bool {
	return unicode.IsSpace(rune(key))
}

func isLetterOrDigit(key byte) bool {
	return unicode.IsLetter(rune(key)) || unicode.IsDigit(rune(key))
}

func isSymbol(key byte) bool {
	return strings.ContainsRune(symbols, rune(key))
}

func checkID(str string) (code int) {
	code, ok := codeMap[str]
	if !ok {
		code = codeMap["ID"]
	}
	return
}

func checkOperator(str string) (code int) {
	code, ok := codeMap[str]
	if !ok {
		code = codeMap["ERR"]
	}
	return
}

func getWord(data []byte, index int) (length int, value string) {
	for index+length < len(data) && isLetterOrDigit(data[index+length]) {
		value += string(data[index+length])
		length++
	}
	return
}

func getNUM(data []byte, index int) (length int, value string) {
	for index+length < len(data) && isDigit(data[index+length]) {
		value += string(data[index+length])
		length++
	}
	return
}

func getOperator(data []byte, index int) (length int, value string) {
	switch string(data[index]) {
	case "+":
		fallthrough
	case "-":
		fallthrough
	case "*":
		fallthrough
	case "/":
		fallthrough
	case "=":
		fallthrough
	case ";":
		fallthrough
	case "(":
		fallthrough
	case ")":
		fallthrough
	case "$":
		length = 1
		value = string(data[index])
	case ":":
		//case ":=":
		fallthrough
	case "<":
		//case "<>":
		//case "<=":
		fallthrough
	case ">":
		//case ">=":
		if index+1 < len(data) && strings.ContainsRune(symbols, rune(data[index+1])) {
			length = 2
			value = string(data[index]) + string(data[index+1])
		} else {
			length = 1
			value = string(data[index])
		}
	}
	return
}

func (rdp *RDP) nextUnitButNotSpace() {
	rdp.unit = scanner(rdp.data, rdp.index)
	rdp.index += rdp.unit.length
	for rdp.unit.typenum == codeMap["SPACE"] {
		rdp.unit = scanner(rdp.data, rdp.index)
		rdp.index += rdp.unit.length
	}
}

func (rdp *RDP) parse() {
	rdp.isValid = true
	rdp.nextUnitButNotSpace()
	if rdp.unit.typenum != codeMap["begin"] {
		rdp.tag = "no begin"
		rdp.isValid = false
		return
	}
	rdp.nextUnitButNotSpace()
	rdp.multiStatements()
	if rdp.unit.typenum != codeMap["end"] {
		rdp.tag = "no end"
		rdp.isValid = false
		return
	}
	rdp.nextUnitButNotSpace()
	if rdp.unit.typenum != codeMap["$"] {
		rdp.tag = "no $"
		rdp.isValid = false
		return
	}
}

func (rdp *RDP) multiStatements() {
	rdp.statement()
	for rdp.unit.typenum == codeMap[";"] {
		rdp.nextUnitButNotSpace()
		rdp.statement()
	}
}

func (rdp *RDP) statement() {
	if rdp.unit.typenum != codeMap["ID"] {
		rdp.tag = "statement error"
		rdp.isValid = false
		return
	}
	rdp.nextUnitButNotSpace()
	if rdp.unit.typenum != codeMap[":="] {
		rdp.tag = "assignment error"
		rdp.isValid = false
		return
	}
	rdp.nextUnitButNotSpace()
	rdp.expression()
}

func (rdp *RDP) expression() {
	rdp.term()
	for rdp.index < len(rdp.data) && (rdp.unit.typenum == codeMap["+"] || rdp.unit.typenum == codeMap["-"]) {
		rdp.nextUnitButNotSpace()
		rdp.term()
	}
}

func (rdp *RDP) term() {
	rdp.factor()
	for rdp.index < len(rdp.data) && (rdp.unit.typenum == codeMap["*"] || rdp.unit.typenum == codeMap["/"]) {
		rdp.nextUnitButNotSpace()
		rdp.factor()
	}
}

func (rdp *RDP) factor() {
	switch rdp.unit.typenum {
	case codeMap["ID"]:
		fallthrough
	case codeMap["NUM"]:
		rdp.nextUnitButNotSpace()
	case codeMap["("]:
		rdp.nextUnitButNotSpace()
		rdp.expression()
		if rdp.unit.typenum != codeMap[")"] {
			rdp.tag = "no )"
			rdp.isValid = false
			return
		}
		rdp.nextUnitButNotSpace()
	default:
		rdp.tag = "expression error"
		rdp.isValid = false
	}
}
