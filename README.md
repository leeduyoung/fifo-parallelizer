# FIFO Parallelizer - AWS SQS FIFO Queue Worker

AWS SQS FIFO 큐에서 메시지를 병렬 처리하는 Go 워커 애플리케이션입니다.

## 주요 특징

### SQS FIFO 큐 지원
- **순서 보장**: MessageGroupId가 같은 메시지들의 FIFO 순서 보장
- **중복 제거**: ContentBasedDeduplication을 통한 자동 중복 제거
- **배치 처리**: 한 번에 여러 메시지 전송 및 처리
- **Long Polling**: 효율적인 메시지 수신을 위한 장시간 폴링

### 고성능 병렬 처리
- **멀티 워커**: 설정 가능한 워커 수로 병렬 처리
- **비동기 처리**: 고루틴 기반 비동기 메시지 처리
- **안전한 종료**: Graceful shutdown 지원
- **에러 핸들링**: 재시도 로직과 백오프 전략

### LocalStack 지원
- **로컬 개발**: LocalStack을 통한 완전한 로컬 개발 환경
- **Docker Compose**: 한 명령어로 전체 환경 구성
- **자동화**: Makefile을 통한 개발 작업 자동화

## 아키텍처
### 데이터 플로우

```bash
AWS SQS FIFO Queue ──┐
                     │
LocalStack (Dev) ────┼──► SQS Client ──► Message Processor
                     │                        │
Real AWS (Prod) ─────┘                        │
                                              ▼
                                        Worker Pool
                                      (Multiple Workers)
                                               │
                                               ▼
                                      Message Handler
                                     (Business Logic)
```

## 프로젝트 구조

```bash
fifo-parallelizer/
├── cmd/
│   └── main.go                    # 애플리케이션 진입점
├── internal/
│   ├── interfaces/
│   │   └── interface.go           # 핵심 인터페이스 정의
│   ├── container/
│   │   └── container.go           # 의존성 주입 컨테이너
│   ├── worker/
│   │   ├── worker_pool.go         # 워커 풀 구현
│   │   └── message_processor.go   # 메시지 처리 로직
│   ├── client/
│   │   └── client.go              # SQS 클라이언트 구현
│   ├── handler/
│   │   └── handler.go             # 메시지 핸들러 구현
│   ├── config/
│   │   └── config.go              # 설정 관리
│   └── types/
│       └── types.go               # 공통 타입 정의
├── scripts/
│   ├── send-same-message-groupd.sh      # 동일 그룹 ID 테스트
│   └── send-different-message-group.sh  # 다른 그룹 ID 테스트
├── docker-compose.yml             # LocalStack 환경 설정
├── Makefile                       # 개발 작업 자동화
├── go.mod                         # Go 모듈 정의
└── README.md                      # 프로젝트 문서
```

## 시작하기

### 필수 요구사항

- **Go 1.24.2+**
- **Docker & Docker Compose**
- **AWS CLI** (테스트용)
- **Make** (선택사항)

### 설치

1. **저장소 클론**
```bash
git clone https://github.com/leeduyoung/fifo-parallelizer.git
cd fifo-parallelizer
```

2. **의존성 설치**
```bash
go mod download
```

3. **LocalStack 환경 시작**
```bash
make test_setup
```

4. **SQS FIFO 큐 생성**
```bash
make create_queue
```

## 사용법

### 기본 실행

#### 1. LocalStack 환경에서 실행
```bash
# 전체 환경 시작
make test_setup

# 큐 생성
make create_queue

# 워커 애플리케이션 실행
make run
```

#### 2. 테스트 메시지 전송

**동일한 MessageGroupId (순서 보장)**
```bash
make send_same_message_groupd
```

**서로 다른 MessageGroupId (병렬 처리)**
```bash
make send_different_message_group
```

#### 3. 큐 상태 모니터링
```bash
# 큐 목록 확인
make list_queues

# 실시간 모니터링 (별도 터미널)
watch -n 2 'make list_queues'
```

### 실제 AWS 환경에서 실행

#### 1. 환경 변수 설정
```bash
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
export AWS_DEFAULT_REGION=ap-northeast-2
export SQS_QUEUE_URL=https://sqs.ap-northeast-2.amazonaws.com/123456789012/your-fifo-queue.fifo
```

#### 2. 애플리케이션 실행
```bash
go run cmd/main.go
```

## 설정

### 환경 변수

| 변수명 | 기본값 | 설명 |
|--------|--------|------|
| `SQS_QUEUE_URL` | `http://localhost:4566/000000000000/test-fifo-queue.fifo` | SQS 큐 URL |
| `MAX_WORKERS` | `5` | 워커 고루틴 수 |
| `VISIBILITY_TIMEOUT` | `30` | 메시지 가시성 타임아웃 (초) |
| `WAIT_TIME_SECONDS` | `20` | Long polling 대기 시간 (초) |
| `MAX_MESSAGES` | `1` | 한 번에 수신할 최대 메시지 수 |
| `ENDPOINT_URL` | `http://localhost:4566` | AWS 엔드포인트 URL (LocalStack용) |


## 테스트

### FIFO 동작 테스트

#### 순서 보장 테스트
```bash
# 동일한 MessageGroupId로 메시지 전송
make send_same_message_groupd

# 워커 실행하여 순서대로 처리되는지 확인
make run
```

예상 결과: 메시지가 1→2→3→4→5 순서로 처리됨

#### 병렬 처리 테스트
```bash
# 서로 다른 MessageGroupId로 메시지 전송
make send_different_message_group

# 워커 실행하여 병렬 처리되는지 확인
make run
```

예상 결과: 여러 워커가 동시에 메시지를 처리

### 성능 테스트

#### 처리량 테스트
```bash
TODO: 
```

## Makefile 명령어

| 명령어 | 설명 |
|--------|------|
| `make test_setup` | LocalStack 환경 시작 |
| `make test_down` | LocalStack 환경 종료 |
| `make create_queue` | FIFO 큐 생성 |
| `make list_queues` | 큐 목록 확인 |
| `make run` | 워커 애플리케이션 실행 |
| `make send_same_message_groupd` | 동일 그룹 ID 테스트 메시지 전송 |
| `make send_different_message_group` | 다른 그룹 ID 테스트 메시지 전송 |
