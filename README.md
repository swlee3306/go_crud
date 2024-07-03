# go_crud

API 방식을 이용한 Database 기본 CRUD 프로그램

## 목차
1. [소개](#소개)
2. [기능](#기능)
3. [파일 구성](#파일-구성)
4. [요구 사항](#요구-사항)
5. [설치](#설치)
6. [사용 방법](#사용-방법)

## 소개

`go_crud`는 Go 언어를 사용하여 API 방식을 통해 데이터베이스의 기본적인 CRUD(Create, Read, Update, Delete) 작업을 수행하는 프로그램입니다. 이 프로그램은 RESTful API를 제공하여 데이터베이스와의 상호작용을 쉽게 할 수 있도록 설계되었습니다.

## 기능

- **Create**: 새로운 데이터 항목을 데이터베이스에 추가합니다.
- **Read**: 데이터베이스에서 데이터를 조회합니다.
- **Update**: 기존 데이터 항목을 수정합니다.
- **Delete**: 데이터 항목을 삭제합니다.


##파일 구성
```
.
├── README.md
├── go.mod
├── go.sum
├── internal
│   ├── dblinker
│   │   ├── dbmd
│   │   │   └── bt_vm.gen.go
│   │   └── main_LoadDb.go
│   ├── sysdef
│   │   └── init.go
│   └── sysenv
│       └── sysenv.go
├── main.go
├── main_LoadEnv.go
├── main_LoadYml.go
├── print_sdk_version.bat
├── print_sdk_version.sh
├── push_sdk_newversion.bat
├── push_sdk_newversion.sh
├── setting.yml
└── utils
    └── router
        └── router.go
```

## 요구 사항

- Go 1.16 이상
- 데이터베이스 (예: MySQL, PostgreSQL 등)
- 필요한 Go 패키지 (Go 모듈을 사용하여 설치 가능)


## 설치

1. 이 저장소를 클론합니다:
    ```sh
    git clone https://github.com/swlee3306/go_crud.git
    cd go_crud
    ```

2. 필요한 Go 패키지를 설치합니다:
    ```sh
    go mod tidy
    ```

3. 데이터베이스 설정을 구성합니다. `config.yaml` 파일을 생성하고 다음과 같은 내용을 추가합니다:
    ```yaml
    database:
      dsn: "root:1234@tcp(localhost:3306)/baton?parseTime=true"

    ```

## 사용 방법

1. 서버를 시작합니다:
    ```sh
    go run main.go
    ```

2. API 엔드포인트를 통해 CRUD 작업을 수행할 수 있습니다. 기본 엔드포인트는 `http://localhost:8080`입니다.

### Create
- **Endpoint**: `POST /api/v1/datastore/data/insert`
- **Request Body**:
    ```json
    {
      "id" : 0,
      "ip" : "1.1.1.2",
      "hostname" : "test7",
      "user" : "tst7",
      "pwd" : "12344",
      "message": "hello7"
    }
    ```

### Read
- **Endpoint**: `GET /api/v1/datastore/data/search`
- **Response**:
    ```json
    {
      "id": 1
    }
    ```

### Update
- **Endpoint**: `PUT /api/v1/datastore/data/update`
- **Request Body**:
    ```json
    {
      "id" : 1,
      "ip" : "1.1.1.2",
      "hostname" : "test7",
      "user" : "tst7",
      "pwd" : "12344",
      "message": "hello7"
    }
    ```

### Delete
- **Endpoint**: `DELETE /api/v1/datastore/data/delete`
- **Response**:
    ```json
    {
      "id": 1
    }
    ```
