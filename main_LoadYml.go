package main

import (
	"errors"
	"io"
	"log"
	"os"

	"baton-om-data-apiservice/internal/sysenv"

	"gopkg.in/yaml.v2"
)

// conf는 설정 파일의 구조체입니다.
type conf struct {
	System   systemConf   `yaml:"system"`
	Database databaseConf `yaml:"database"`
}

type systemConf struct {
	PrvId int64 `yaml:"prvId"`
}

type databaseConf struct {
	Dsn string `yaml:"dsn"`
}

// main_LoadYml은 지정된 파일에서 설정을 로드하는 함수입니다.
// 
// 매개변수:
//   - fname: 설정 파일의 경로
//
// 반환값:
//   - error: 설정 로드 중 발생한 오류

func main_LoadYml(fname string) error {
	var c conf

	// get
	err := _loadCfgFromFile(fname, &c)
	if err != nil {
		return err
	}

	sysenv.Database.Dsn = c.Database.Dsn

	return nil
}

// _loadCfgFromFile 함수는 지정된 파일에서 설정을 로드하는 함수입니다.
// 만약 파일이 존재하지 않는 경우, 새로운 파일을 생성하고 기본 설정을 저장합니다.
// 파일이 존재하는 경우, 파일에서 설정을 읽어와 구조체에 언마샬링합니다.
// fname은 설정 파일의 경로입니다.
// c는 설정을 저장할 구조체입니다.
// 함수가 에러를 반환할 경우, 해당 에러를 반환합니다.
func _loadCfgFromFile(fname string, c *conf) error {
	var f *os.File
	var err error

	f, err = os.OpenFile(fname, os.O_RDONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("config file not found: %s", fname)
			f, err = os.Create(fname)
			if err != nil {
				log.Printf("creating config error: %v", err)
				return err
			}
			defer f.Close()

			data, err := yaml.Marshal(c)
			if err != nil {
				log.Printf("generating config fail!: %v", err)
				return err
			}

			f.WriteString(string(data))
		} else {
			log.Printf("config is there but couldnot open it!: %v", err)
			return err
		}
	} else {
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			log.Printf("reading config fail!: %v", err)
			return err
		}

		err = yaml.Unmarshal(data, c)
		if err != nil {
			log.Printf("parsing config fail!: %v", err)
			return err
		}
	}

	return nil
}
