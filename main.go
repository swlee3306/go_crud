package main

import (
	"baton-om-data-apiservice/internal/sysdef"
	"baton-om-data-apiservice/internal/sysenv"
	"baton-om-data-apiservice/utils/router"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"bitbucket.org/okestrolab/baton-ao-sdk/btoutil"
	"gorm.io/gorm"
)

var (
	X_buildDatetime, X_buildRevision, X_buildRevisionShort, X_buildBranch, X_buildTag string
	Db *gorm.DB
)

func main() {

	// 함수 이름과 에러 변수를 초기화합니다.
	fnc := "main"
	err := error(nil)

	// 빌드 정보를 로그에 출력합니다.
	log.Printf("%s: 빌드 정보", fnc)
	log.Printf("\t buildDatetime: %s", X_buildDatetime)
	log.Printf("\t buildRevision: %s (%s)", X_buildRevisionShort, X_buildRevision)
	log.Printf("\t buildBranch: %s", X_buildBranch)
	log.Printf("\t buildTag: %s", X_buildTag)

	{
		// 시간 위치 설정
		//btoutil.SetDefaultTimeZone("UTC")
		btoutil.SetDefaultTimeZone("Asia/Seoul")

		// 디버그 모드 설정
		if val, ok := os.LookupEnv("OKE_DEBUG"); ok && (len(val) > 0) {
			sysenv.Mode.IsDebug, _ = strconv.ParseBool(val)
		}

		// 설정 파일 이름 설정
		if val, ok := os.LookupEnv("BATON_SETTING_FILENAME"); ok && (len(val) > 0) {
			sysdef.ConfFilename = val
		}

		// 설정 파일 로드
		if err = main_LoadYml(sysdef.ConfFilename); err != nil {
			panic(fmt.Sprintf("%s: Cfg load 실패: %s", fnc, err.Error()))
		}
	}

	// Gin 서버를 시작합니다.
	router.StartGinServer()

	// 무한 루프
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		log.Printf("인터럽트 발생: %v", sig)

		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			break
		}
	}

}
