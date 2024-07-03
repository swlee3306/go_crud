package dblinker

import (
	"baton-om-data-apiservice/internal/dblinker/dbmd"
	"baton-om-data-apiservice/internal/sysenv"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error){
	var err error
	db, err := gorm.Open(mysql.Open(sysenv.Database.Dsn), &gorm.Config{})
	if err != nil {
			log.Printf("failed to connect database: %v", err)
			return nil, err
	}
	return db, nil
}


func LoadModule(db *gorm.DB) ([]*sysenv.Data, error) {

	var vmList []*dbmd.BtVM
	res := db.Raw("select * from bt_vm br where (br.del_yn= 'n' or br.del_yn= 'N')").Scan(&vmList)
	if res.Error != nil {
		return nil, res.Error
	}

	VmsList := make([]*sysenv.Data, len(vmList))

	for i, v := range vmList {
		VmsList[i] = &sysenv.Data{
			Id:       v.ID,
			Hostname: v.Hostname,
			User:     v.Hostuser,
			Ip:       v.HostIP,
			Pwd:      v.HostPwd,
			Message:  v.Message,
		}
	}

	return VmsList, nil
}

func Insert(db *gorm.DB, vm *sysenv.Data) error {

	{
		v := &dbmd.BtVM{
			Hostname: vm.Hostname,
			Hostuser: vm.User,
			HostIP:   vm.Ip,
			HostPwd:  vm.Pwd,
			Message:  vm.Message,
		}

		if res := db.Create(v); res.Error != nil {
			log.Printf("db.Insert failed: Hostname(%s): %v", vm.Hostname, res.Error)
		} else {
			log.Printf("db.Insert succeed: Hostname(%s)", vm.Hostname)
		}
	}

	return nil
}

func Search(db *gorm.DB, vm *sysenv.DataSearch) (*dbmd.BtVM, error) {

	var vmList []*dbmd.BtVM
	res := db.Raw("select * from bt_vm br where br.id = ?", vm.Id).Scan(&vmList)
	if res.Error != nil {
		return nil, res.Error
	}

	return vmList[0], nil
}

func Delete(db *gorm.DB, vm *sysenv.DataSearch) error {

	if res := db.Table("bt_vm").Where("id = ?", vm.Id).Updates(map[string]interface{}{"del_dt": time.Now()}); res.Error != nil {
		log.Printf("db.Delete failed Id(%d): %s", vm.Id, res.Error.Error())
	}

	if res := db.Table("bt_vm").Where("id = ?", vm.Id).Updates(map[string]interface{}{"del_yn": "Y", "del_dt": time.Now()}); res.Error != nil {
		log.Printf("db.Delete sucess: Id(%d): %s", vm.Id, res.Error.Error())
	}

	return nil
}

func Update(db *gorm.DB, vm *sysenv.Data) error {

	var vmList []*dbmd.BtVM
	res := db.Raw("select * from bt_vm br where br.id = ?", vm.Id).Scan(&vmList)
	if res.Error != nil {
		return res.Error
	}

	if vmList[0].ID != vm.Id {
		vmList[0].ID = vm.Id
	}

	if vmList[0].Hostname != vm.Hostname {
		vmList[0].Hostname = vm.Hostname
	}

	if vmList[0].HostIP != vm.Ip {
		vmList[0].HostIP = vm.Ip
	}

	if vmList[0].Hostuser != vm.User {
		vmList[0].Hostuser = vm.User
	}

	if vmList[0].HostPwd != vm.Pwd {
		vmList[0].HostPwd = vm.Pwd
	}

	if vmList[0].Message != vm.Message {
		vmList[0].Message = vm.Message
	}

	if err := db.Save(vmList[0]).Error; err != nil {
		return err
	}
	
	return nil
}
