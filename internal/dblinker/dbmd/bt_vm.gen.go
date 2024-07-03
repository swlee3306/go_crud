package dbmd

import "time"

const TableNameBtVM = "bt_vm"

// BtVM mapped from table <bt_vm>
type BtVM struct {
	ID       int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Hostname string    `gorm:"column:hostname" json:"hostname"`
	Hostuser string    `gorm:"column:host_user" json:"host_user"`
	HostIP   string    `gorm:"column:host_ip" json:"host_ip"`
	HostPwd  string    `gorm:"column:host_pwd" json:"host_pwd"`
	Message  string    `gorm:"column:message" json:"message"`
	DelYn    string    `gorm:"column:del_yn;not null;default:N"`
	RegDt    time.Time `gorm:"column:reg_dt;default:now()"`
	ModDt    time.Time `gorm:"column:mod_dt"`
	DelDt    time.Time `gorm:"column:del_dt"`
}

// TableName BtVM's table name
func (*BtVM) TableName() string {
	return TableNameBtVM
}
