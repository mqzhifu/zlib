package zlib

import "reflect"

type Goroutine struct {
	Log *Log
	Metrics *Metrics
}


func NewGoroutine(log *Log,metrics *Metrics)*Goroutine{
	self := new(Goroutine)
	self.Log = log
	self.Metrics = metrics
	return self
}

func (goroutine *Goroutine)  CreateExec(class interface{},funcNme string,argc ...interface{} ) {
	//return  []reflect.Value
	goroutine.Log.Warning("NewGoroutine ",funcNme," ,len: ",len(argc), " argc:",argc)
	//getType := reflect.TypeOf(gamematch)
	//for i:= 0; i<getType.NumMethod();i++{
	//	met := getType.Method(i)
	//	//if met.Name == funcNme{
	//	//
	//	//	break
	//	//}
	//	zlib.MyPrint(met.Name, met.Index)
	//	//fmt.Printf("%s, %s, %s, %d\n", met.Type, met.Type.Kind(), met.Name, met.Index)
	//}
	//myfunctype,err := reflect.TypeOf(gamematch).MethodByName(funcNme)
	//zlib.MyPrint(err,myfunctype.Type)
	//var inputs []reflect.Value
	inputs := make([]reflect.Value, len(argc))
	for i, _ := range argc {
		inputs[i] = reflect.ValueOf(argc[i])
	}

	myfunc := reflect.ValueOf(class).MethodByName(funcNme)
	//var aaaa []reflect.Value
	//in := make([]reflect.Value)
	//aa := []reflect.Value{inputs}
	goroutine.Metrics.IncNode("goroutineNum")
	go myfunc.Call( inputs)
	//goroutine.Metrics.LessNode("goroutineNum")
	//for k,v := range rs{
	//	zlib.MyPrint(k,v.Interface())
	//}
	//gg := rs.(Group)
	//gg := rs[0].Interface().(Group)
	//zlib.ExitPrint(gg,1111)
	//return rs
	//zlib.MyPrint(gamematch.PidFilePath)
	//sign := gamematch.getContainerSuccessByRuleId(999)
	//zlib.MyPrint(sign.Rule)
}