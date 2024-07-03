package router

import (
	"baton-om-data-apiservice/internal/dblinker"
	"baton-om-data-apiservice/internal/sysenv"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartGinServer() {
	fnc := "StartGinServer"
	log.Printf("%s Run", fnc)

	r := gin.Default()

	db, err := dblinker.InitDB()
	if err != nil {
		log.Printf("%s DB connection Fail : %s", fnc, err)
	}

	// 패닉 및 에러 처리를 위해 리커버리 미들웨어 사용
	r.Use(gin.Recovery())

	dataApi := r.Group("/api/v1/datastore")
	{
		crud := dataApi.Group("data")
		{
			crud.POST("/insert", func(c *gin.Context) {
				var err error
				var req sysenv.Data
				if err = c.BindJSON(&req); err != nil {
					log.Printf("%s: c.BindJSON failed: %s", fnc, err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"message": ":" + err.Error()})
					return
				}

				err = dblinker.Insert(db, &req)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"message": "Data insert fail"})
				}
				c.JSON(http.StatusOK, gin.H{"message": "Data insert successfully"})
			})
			crud.GET("/search", func(c *gin.Context) {
				var err error
				var req sysenv.DataSearch
				if err = c.BindJSON(&req); err != nil {
					log.Printf("%s: c.BindJSON failed: %s", fnc, err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"message": ":" + err.Error()})
					return
				}

				vm, err := dblinker.Search(db, &req)
				if err != nil {
					log.Printf("%s: dblinker.Search: %s", fnc, err.Error())
					c.JSON(http.StatusOK, gin.H{"message": "Data Search fail", "err": err.Error()})
				}

				var res sysenv.Data

				res.Id = vm.ID
				res.Hostname = vm.Hostname
				res.User = vm.Hostuser
				res.Ip = vm.HostIP
				res.Pwd = vm.HostPwd
				res.Message = vm.Message

				c.JSON(http.StatusOK, res)
			})
			crud.DELETE("/delete", func(c *gin.Context) {
				var err error
				var req sysenv.DataSearch
				if err = c.BindJSON(&req); err != nil {
					log.Printf("%s: c.BindJSON failed: %s", fnc, err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"message": ":" + err.Error()})
					return
				}

				err = dblinker.Delete(db, &req)
				if err != nil {
					log.Printf("%s: dblinker.Delete: %s", fnc, err.Error())
					c.JSON(http.StatusOK, gin.H{"message": "Data Delete fail", "err": err.Error()})
				}

				c.JSON(http.StatusOK, gin.H{"message": "Data delete successfully"})
			})
			crud.PUT("/update", func(c *gin.Context) {
				var err error
				var req sysenv.Data
				if err = c.BindJSON(&req); err != nil {
					log.Printf("%s: c.BindJSON failed: %s", fnc, err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"message": ":" + err.Error()})
					return
				}

				err = dblinker.Update(db, &req)
				if err != nil {
					log.Printf("%s: dblinker.Update: %s", fnc, err.Error())
					c.JSON(http.StatusOK, gin.H{"message": "Data Update fail", "err": err.Error()})
				}

				c.JSON(http.StatusOK, gin.H{"message": "Data update successfully"})
			})
			crud.GET("/all", func(c *gin.Context) {

				if Vmlist, err := dblinker.LoadModule(db); err != nil {
					log.Printf("%s: dblinker.LoadModule() failed: %s", fnc, err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"message": ":" + err.Error()})
					return
				} else if len(Vmlist) == 0 {
					c.JSON(http.StatusOK, gin.H{"message": "Data none"})
					return
				} else {

					var vmlists []*sysenv.DataRes

					for _, vm := range Vmlist {
						res := sysenv.DataRes{
							Id:       vm.Id,
							Hostname: vm.Hostname,
							Ip:       vm.Ip,
							User:     vm.User,
							Message:  vm.Message,
						}
						vmlists = append(vmlists, &res)
					}

					res_e := sysenv.VmList{
						Vm: vmlists,
					}
					c.JSON(http.StatusOK, res_e)
				}
			})
		}
	}

	err = r.Run(":8080")
	if err != nil {
		log.Printf("%s error: %v", fnc, err)
	}

}
