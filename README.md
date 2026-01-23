# Chat-Server
Go 기반 실시간 채팅 서버

## 📋 프로젝트 개요
실시간 채팅 기능을 제공하는 백엔드 서버입니다.

## 🛠 기술 스택
- **Language**: Go 1.21+
- **Framework**: Gin 
- **Database**: PostgreSQL
- **ORM**: GORM
- **WebSocket**: Gorilla WebSocket

## 📁 프로젝트 구조

```
Chat-Server/
├── cmd/
│   └── server/
│       └── main.go                  # 메인 애플리케이션
├── internal/
│   ├── config/                      # 설정 
│   │   └── config.go
│   ├── handler/                     # HTTP 핸들러
│   ├── service/                     # 비즈니스 로직
│   ├── repository/                  # DB 접근 레이어
│   ├── middleware/                  # 미들웨어
│   ├── models/                      # DB 모델 
│   └── websocket/                   # WebSocket 핸들러
├── pkg/                             # 외부에서 사용 가능한 라이브러리 (필요시)
│   └── utils/                       # 공용 유틸리티
├── api/                             # API 정의 (OpenAPI/Swagger 등)
├── configs/                         # 설정 파일들 (.yaml, .toml 등)
│   └── config.yaml
├── scripts/                         # 빌드/설치 스크립트
├── deployments/                     # 배포 설정 (docker-compose, k8s 등)
├── test/                            # 추가 테스트 앱/데이터
├── docs/                            # 문서
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md
```

## 🔧 개발 환경 설정

### 의존성 설치
`go mod download`

### 서버 실행
`go run cmd/server/main.go`

