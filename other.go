package zlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
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
//debug 调试使用
func ExitPrint(a ...interface{})   {
	fmt.Println(a)
	fmt.Println("ExitPrint...22")
	os.Exit(-22)
}
//获取一个随机整数
func GetRandIntNum(max int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
//获取一个随机整数 范围
func GetRandIntNumRange(min int ,max int) int{
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
//判断一个字符串是否为空，包括  空格
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
//检查一个文件是否已存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename);os.IsNotExist(err){
		exist = false
	}
	return exist
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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
//将字符串的首字母转大写
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return str
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122  {
		strArry[0] = strArry[0] - 32
	}
	return string(strArry)
}

func MapCovertStruct(inMap map[string]interface{},outStruct interface{})interface{}{
	fmt.Printf("%+v",inMap)
	fmt.Printf("%+v",outStruct)

	setFiledValue := func(	outStruct interface{},name string , v interface{}) {
		MyPrint(name)
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

func ElementStrInArrIndex(arr []string ,element string )int{
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
var GoRoutineList = make( map[string]int )
func AddRoutineList(name string){
	GoRoutineList[name] = GetNowTimeSecondToInt()
}
func GetLocalIp()(ip string,err error){
	netInterfaces, err := net.Interfaces()
	//MyPrint(netInterfaces, err)
	if err != nil {
		return ip,errors.New("net.Interfaces failed, err:" +  err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(),nil
					}
				}
			}
		}
	}

	return ip,nil
}
//在一个：一维数组中，找寻最大数
func FindMaxNumInArrFloat32(arr []float32  )float32{
	number := arr[0]
	for _,v := range arr{
		if v > number{
			number = v
		}
	}
	return number
}

//在一个：一维数组中，找寻最小数
func FindMinNumInArrFloat32(arr []float32  )float32{
	number := arr[0]
	for _,v := range arr{
		if v < number{
			number = v
		}
	}
	return number
}
//4舍5入，保留2位小数
//func round(x float32)string{
//	numberStr := FloatToString(x,4)
//	numberStrSplit :=  strings.Split(numberStr,".")
//	if len(numberStrSplit) == 1{
//		return numberStrSplit[0]
//	}
//	numberLittleStrByte := []byte(numberStrSplit[1])
//	numberLittle := Atoi(numberStrSplit[1])
//	if len(numberLittleStrByte) == 4{//4位小数
//		three := numberLittleStrByte[0] + numberLittleStrByte[1] + numberLittleStrByte[2]
//		if strByte[4] >= 5{
//			three = Atoi(three) + 1
//		}
//	}else{//3位小数
//
//	}
//
//}