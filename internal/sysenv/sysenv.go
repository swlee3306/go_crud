package sysenv

var Mode struct {
	IsDebug bool
}
var Database struct {
	Dsn             string
	MaxIdleConns    int
	MaxOpenConns    int
	MaxLifetimeHour int
}

type Data struct{
	Id int64 `json:"id"`
	Ip string `json:"ip"`
	Hostname string `json:"hostname"`
	User string `json:"user"`
	Pwd string `json:"pwd"`
	Message string `json:"message"`
}

type DataRes struct{
	Id int64 `json:"id"`
	Ip string `json:"ip"`
	Hostname string `json:"hostname"`
	User string `json:"user"`
	Message string `json:"Message"`
}

type DataSearch struct{
	Id int64 `json:"id"`
}

type VmList struct{
	Vm []*DataRes `json:"vm_list"`
}
