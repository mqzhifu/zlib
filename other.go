package zlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//这个函数只是懒......
func MyPrint(a ...interface{}) (n int, err error) {
	if LogLevelFlag == LOG_LEVEL_DEBUG{
		return fmt.Println(a)
	}
	return
}

func ExitPrint(a ...interface{})   {
	fmt.Println(a)
	fmt.Println("ExitPrint...22")
	os.Exit(-22)
}

func GetRandIntNum(max int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func CheckStrEmpty(str string)bool{
	if str == ""{
		return true
	}
	str = strings.Trim(str," ")
	if str == ""{
		return true
	}
	return false
}


func MapCovertArr( myMap map[int]int) (arr []int){
	for _,v := range myMap {
		arr = append(arr,v)
	}
	return arr
}

func ArrCovertMap(arr []int )map[int]int{
	mapArr := make(map[int]int)
	for k,v := range arr {
		mapArr[k] = v
	}
	return mapArr
}

func ArrStringCoverArrInt(arr []string )(arr2 []int){
	for i:=0;i<len(arr);i++{
		arr2 = append(arr2, Atoi(arr[i]))
	}
	return arr2
}
func GetSpaceStr(n int)string{
	str := ""
	for i:=0;i<n;i++{
		str += " "
	}
	return str
}
//检查已经make过的，二维map int 类型，是否为空
func CheckMap2IntIsEmpty(hashMap map[int]map[int]int)bool{
	if len(hashMap) == 0{
		return true
	}

	for _,v := range hashMap{
		if len(v) > 0{
			return false
		}
	}
	return true

}


func ArrCoverStr(arr []int,IdsSeparation string)string{
	if len(arr) == 0{
		ExitPrint("ArrCoverStr arr len = 0")
	}
	str := ""
	for _,v := range arr{
		str +=  strconv.Itoa(v) + IdsSeparation
	}
	str = str[0:len(str)-1]
	return str
}

func StructCovertMap(inStruct interface{})interface{}{
	jsonStr ,_:= json.Marshal(inStruct)
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {

	}
	return mapResult
}
//strconv.Atoi 返回的是两个参数，很麻烦，这里简化，只返回一个参数
func Atoi(str string)int{
	num, _ := strconv.Atoi(str)
	return num
}

func MapCovertStruct(inMap map[string]interface{},outStruct interface{})interface{}{
	setFiledValue := func(	outStruct interface{},name string , v interface{}) {

		structValue := reflect.ValueOf(outStruct).Elem()
		structFieldValue := structValue.FieldByName(name)

		structFieldType := structFieldValue.Type() //结构体的类型
		val := reflect.ValueOf(v)              //map值的反射值

		var err error
		if structFieldType != val.Type() {
			val, err = TypeConversion(fmt.Sprintf("%v", v), structFieldValue.Type().Name()) //类型转换
			if err != nil {
				ExitPrint(err.Error())
			}
		}
		//MyPrint(val,val.Type(),v)

		structFieldValue.Set(val)
	}
	for k,v := range inMap{
		//MyPrint("MapCovertStruct for range:",outStruct,k,v)
		setFiledValue(outStruct,k,v)
	}
	//outStructV := reflect.ValueOf(outStruct)
	//outStructT := reflect.TypeOf(outStruct)

	return outStruct
}


//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
//判断一个元素，在一个数组中的位置
func ElementInArrIndex(arr []int ,element int )int{
	for i:=0;i<len(arr);i++{
		if arr[i] == element{
			return i
		}
	}
	return -1
}

func FloatToString(number float32,little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(number), 'f', little, 64)
}

func Float64ToString(number float64,little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat( number, 'f', little, 64)
}

func GetNowTimeSecondToInt()int{
	return int( time.Now().Unix() )
}

func StringToFloat(str string)float32{
	v1,_ := strconv.ParseFloat(str, 32)
	number  := float32(v1)
	return number
}
