package zlib

type Metrics struct {
	Container map[string]int
}

func NewMetrics()*Metrics{
	metrics := new (Metrics)
	metrics.Container = make(map[string]int)
	return metrics
}

func (metrics *Metrics) CreateOneNode(name string){
	metrics.Container[name] = 0
}

func (metrics *Metrics)PlusNode(name string ,num int){
	metrics.Container[name] = metrics.Container[name] + num
}

func (metrics *Metrics)SetNode(name string ,num int){
	metrics.Container[name] = num
}

func (metrics *Metrics)IncNode(name string  ){
	metrics.Container[name] ++
}


func (metrics *Metrics)GetAllAndClear()map[string]int{
	container := metrics.Container
	metrics.Container = make(map[string]int)
	return  container
}

func (metrics *Metrics)GetAll()map[string]int{
	container := metrics.Container
	return  container
}