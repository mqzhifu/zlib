package zlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

const (
	LEVEL_INFO 		=	 1 << iota
	LEVEL_DEBUG		//2
	LEVEL_ERROR		//4
	LEVEL_PANIC		//8

	LEVEL_EMERGENCY	//16
	LEVEL_ALERT		//32
	LEVEL_CRITICAL	//64
	LEVEL_WARNING	//128
	LEVEL_NOTICE	//256
)
var levelContentPrefixes = map[int]string{
	LEVEL_INFO: "INFO",
	LEVEL_DEBUG: "DEBUG",
	LEVEL_ERROR: "ERROR",
	LEVEL_PANIC: "PANIC",
	LEVEL_EMERGENCY: "EMERG",
	LEVEL_ALERT: "ALERT",
	LEVEL_CRITICAL: "CRITI",
	LEVEL_WARNING: "WARNI",
	LEVEL_NOTICE: "NOTIC",
}

const(
	LEVEL_ALL = LEVEL_INFO | LEVEL_DEBUG | LEVEL_ERROR | LEVEL_PANIC | LEVEL_EMERGENCY |LEVEL_ALERT| LEVEL_CRITICAL |LEVEL_WARNING |LEVEL_NOTICE
	LEVEL_DEV = LEVEL_INFO | LEVEL_DEBUG | LEVEL_ERROR | LEVEL_PANIC
	LEVEL_ONLINE =  LEVEL_INFO | LEVEL_ERROR | LEVEL_PANIC
)

const(
	OUT_TARGET_SC = 	 1 << iota
	OUT_TARGET_FILE
	OUT_TARGET_NET
)

const(
	OUT_TARGET_ALL = OUT_TARGET_SC|OUT_TARGET_FILE|OUT_TARGET_NET
)

const(
	OUT_TARGET_NET_TCP  = 1
	OUT_TARGET_NET_UDP = 2
)



type Log struct {
	option LogOption
	Op LogOption
}

type LogOption struct {
	OutFilePath 	string
	OutFileName 	string
	OutFileHashType	int
	OutFileFd		*os.File
	Level 			int
	Target 			int
}
func NewLog( logOption LogOption)(log *Log ,errs error){
	//MyPrint("New log class ,OutFilePath : ",logOption.OutFilePath , " level : ",logOption.Level ," target : ",logOption.Target)

	if logOption.OutFilePath == ""{
		return nil,errors.New(" OutFilePath is empty ")
	}

	if logOption.Level == 0{
		return nil,errors.New(" level is empty ")
	}

	if logOption.Target == 0 {
		return nil,errors.New(" target is empty ")
	}

	log = new(Log)
	log.option = logOption

	errs = log.checkOutFilePathPower(logOption.OutFilePath)
	if errs != nil{
		return nil,errs
	}

	log.option = logOption
	log.Op = logOption
	if log.checkTargetIncludeByBit(OUT_TARGET_FILE){
		if log.option.OutFileHashType == 0{//未开启hash
			pathFile := log.GetPathFile()
			fd, err  := os.OpenFile(pathFile, os.O_WRONLY | os.O_CREATE | os.O_APPEND , 0777)
			if err != nil{
				return nil,errors.New(" log out file , OpenFile :  " + err.Error())
			}
			log.option.OutFileFd = fd
		}
	}
	msg := "NewLogClass , OutFilePath : "+logOption.OutFilePath +" level : " + strconv.Itoa(logOption.Level) +" target : "+ strconv.Itoa(logOption.Target)
	log.Debug(msg)

	return log,nil
}

//permission
func  (log *Log)  checkOutFilePathPower(path string)error{
	if path == ""{
		return errors.New(" checkOutFilePathPower ("+path+") : path is empty")
	}

	fd,e :=  os.Stat(path)
	if e != nil{
		return errors.New(" checkOutFilePathPower ("+path+"): os.Stat : "+ e.Error())
	}

	if !fd.IsDir(){
		return errors.New(" checkOutFilePathPower ("+path+"): path is not a dir ")
	}
	perm := fd.Mode().Perm().String()
	//MyPrint(perm,os.FileMode(0755).String())
	//log.Debug(fd.Mode(),fd.Mode().Perm())
	if perm < os.FileMode(0755).String(){
		return errors.New(" checkOutFilePathPower ("+path+"):path permission 0777 ")
	}
	return nil
}

func (log *Log) Info(content ...interface{}){
	log.Out(LEVEL_INFO,content)
}

func (log *Log) Debug(content ...interface{}){
	log.Out(LEVEL_DEBUG,content...)
}

func (log *Log) Error(content ...interface{}){
	log.Out(LEVEL_ERROR,content...)
}

func (log *Log) Notice(content ...interface{}){
	log.Out(LEVEL_NOTICE,content...)
}

func (log *Log) Warning(content ...interface{}){
	log.Out(LEVEL_WARNING,content...)
}

func (log *Log) OutScreen(a ...interface{}){
	if a[0] == "[INFO]"{
		fmt.Printf("%c[1;40;33m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[ERROR]" {
		fmt.Printf("%c[1;40;31m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[NOTIC]" {
		fmt.Printf("%c[1;40;34m%s%c[0m", 0x1B, a[0], 0x1B)
	}else if a[0] == "[WARNI]" {
		fmt.Printf("%c[1;40;35m%s%c[0m", 0x1B, a[0], 0x1B)
	}else{
		fmt.Printf("%c[1;40;32m%s%c[0m", 0x1B, a[0], 0x1B)
	}

	newlist := append(a[:0], a[(0+1):]...)
	fmt.Println(newlist...)
}
func (log *Log) GetPathFile()string{
	return log.option.OutFilePath + "/" + log.option.OutFileName
}
func (log *Log) OutFile(content string){
	_, err := io.WriteString(log.option.OutFileFd, content)
	if err != nil{
		log.Error("OutFile io.WriteString : ",err.Error())
	}
}

func (log *Log) CloseFileFd()error{
	if !log.checkTargetIncludeByBit(OUT_TARGET_FILE){
		return errors.New("checkTargetIncludeByBit OUT_TARGET_FILE :false")
	}
	err := log.option.OutFileFd.Close()
	return err
}

func (log *Log)getHeaderContentStr()string{
	//timeStr:=time.Now().Format("2006-01-02 15:04:05")
	//unixstamp := GetNowTimeSecondToInt()
	//uuid4 := getUuid4()
	pid  := os.Getpid()
	//str := timeStr + "[" + strconv.Itoa(pid) + "]"
	str :=   "[" + strconv.Itoa(pid) + "]"
	return str
}

func (log *Log) checkTargetIncludeByBit(flag int)bool{
	if log.option.Target & flag == flag {
		return true
	}
	return false
}

func (log *Log) checkLevelIncludeByBit(level int)bool{
	//MyPrint(log.option.Level,level)
	if log.option.Level & level == level {
			return true
	}
	return false
}

func  (log *Log)Out(level int ,argcs ...interface{}){
	if !log.checkLevelIncludeByBit(level){
		return
	}
	contentLevelPrefix := "[" + levelContentPrefixes[level] +"]"
	content := ""
	for _,argc := range argcs{
		content += " " + log.String(argc)
	}

	if log.checkTargetIncludeByBit(OUT_TARGET_FILE){
		log.OutFile(contentLevelPrefix + log.getHeaderContentStr() + content + "\n")
	}

	if log.checkTargetIncludeByBit(OUT_TARGET_SC){
		log.OutScreen(contentLevelPrefix,log.getHeaderContentStr(),content)
	}

	if log.checkTargetIncludeByBit(OUT_TARGET_NET){

	}

	//ExitPrint(-200)
}

//https://github.com/gogf/gf/tree/master/os/glog
type apiString interface {
	String() string
}

type apiError interface {
	Error() string
}

func  (log *Log)String(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	//case time.Time:
	//	if value.IsZero() {
	//		return ""
	//	}
	//	return value.String()
	//case *time.Time:
	//	if value == nil {
	//		return ""
	//	}
	//	return value.String()
	//case gtime.Time:
	//	if value.IsZero() {
	//		return ""
	//	}
	//	return value.String()
	//case *gtime.Time:
	//	if value == nil {
	//		return ""
	//	}
	//	return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(apiString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(apiError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return log.String(rv.Elem().Interface())
		}
		// Finally we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

//func  (log *Log)parseInterfaceValueCovertStr(interValue interface{})string{
//
//	switch f := interValue.(type) {
//		case bool:
//			if f {
//				return "true"
//			}else{
//				return "false"
//			}
//		case float32:
//			return FloatToString(interValue.(float32),3)
//		case float64:
//			return Float64ToString(interValue.(float64),3)
//		//case complex64:
//		//	p.fmtComplex(complex128(f), 64, verb)
//		//case complex128:
//		//	p.fmtComplex(f, 128, verb)
//		case int:
//			return strconv.Itoa(interValue.(int))
//		case int8:
//			return strconv.Itoa(int(interValue.(int8)))
//		case int16:
//			strconv.Itoa(int (interValue.(int16)))
//		case int32:
//			strconv.FormatInt(int64 (interValue.(int32)),10)
//		case int64:
//			strconv.FormatInt(interValue.(int64),10)
//		case uint:
//			strconv.FormatUint(uint64(interValue.(uint)),10)
//		case uint8:
//			strconv.FormatUint(uint64(interValue.(uint8)),10)
//		case uint16:
//			strconv.FormatUint(uint64(interValue.(uint16)),10)
//		case uint32:
//			strconv.FormatUint(uint64(interValue.(uint32)),10)
//		case uint64:
//			strconv.FormatUint(interValue.(uint64),10)
//		case uintptr:
//			p.fmtInteger(uint64(f), unsigned, verb)
//		case string:
//			return interValue.(string)
//		case []byte:
//			return interValue.(string)
//		case reflect.Value:
//			if f.IsValid() && f.CanInterface() {
//				p.arg = f.Interface()
//				if p.handleMethods(verb) {
//					return
//				}
//			}
//			p.printValue(f, verb, 0)
//		default:
//			if !p.handleMethods(verb) {
//				p.printValue(reflect.ValueOf(f), verb, 0)
//			}
//
//	}
//}

//var levelStringMap = map[string]int{
//	"ALL":      LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEV":      LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEVELOP":  LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"PROD":     LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"PRODUCT":  LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEBU":     LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"DEBUG":    LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"INFO":     LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"NOTI":     LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"NOTICE":   LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"WARN":     LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"WARNING":  LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT,
//	"ERRO":     LEVEL_ERRO | LEVEL_CRIT,
//	"ERROR":    LEVEL_ERRO | LEVEL_CRIT,
//	"CRIT":     LEVEL_CRIT,
//	"CRITICAL": LEVEL_CRIT,
//}