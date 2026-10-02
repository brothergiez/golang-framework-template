Berikut PRD yang saya sarankan untuk boilerplate ini. Fokusnya bukan sekadar template CRUD, tetapi menjadi foundation backend production yang bisa dipakai untuk berbagai jenis project dan tetap mudah dikembangkan ketika traffic, jumlah service, atau kompleksitas domain bertambah.

# PRD: Go Gin Production Backend Boilerplate

## 1. Product Overview

Nama sementara:

Go Gin Production Boilerplate

Tujuan utama adalah menyediakan base template backend menggunakan Go 1.27.0 yang memiliki karakteristik:

* Robust
* Scalable
* Reliable
* Maintainable
* Testable
* SOLID
* Database agnostic
* Mudah digunakan untuk monolith maupun microservice
* Mendukung beberapa service type: `api`, `cron`, `worker`, dan `lambda`
* Mendukung beberapa database
* Redis bersifat optional
* Memiliki structured logging
* Memiliki configuration management
* Memiliki standard error handling
* Memiliki observability dasar
* Siap digunakan sebagai foundation aplikasi production

Satu codebase harus dapat dijalankan sebagai HTTP API (Gin), scheduled job (cron), background worker, atau serverless function (lambda) tanpa menduplikasi business logic.

Boilerplate ini harus menghindari ketergantungan langsung antar layer sehingga database, cache, HTTP framework, logger, dan external service dapat diganti tanpa mengubah business logic.

---

# 2. Goals

## 2.1 Primary Goals

Boilerplate harus menyediakan:

1. Multi service type support (`api`, `cron`, `worker`, `lambda`)
2. HTTP server menggunakan Gin untuk service type `api`
3. Cron / scheduler runner untuk service type `cron`
4. Background worker runner untuk service type `worker`
5. Serverless / Lambda adapter untuk service type `lambda`
6. Go 1.27.0
7. Clean architecture
8. SOLID principles
9. Dependency Injection
10. Shared domain & application layer lintas service type
11. Multi database support
12. PostgreSQL
13. MySQL / MariaDB
14. Microsoft SQL Server
15. MongoDB
16. Redis optional
17. Structured logging
18. Configuration management
19. Request validation
20. Standard API response
21. Centralized error handling
22. Middleware architecture untuk `api`
23. Graceful shutdown per service type
24. Health check untuk `api`
25. Readiness check untuk `api`
26. Liveness check untuk `api`
27. Unit testing
28. Integration testing
29. Docker support per service type
30. Environment based configuration
31. Secure default configuration
32. Request ID / correlation ID
33. API documentation support
34. Metrics support
35. Distributed tracing ready

---

# 3. Non Goals

Versi pertama tidak harus menyediakan:

* Business domain tertentu
* Authentication implementation yang terlalu spesifik
* Admin dashboard
* Frontend
* Message broker implementation wajib
* Kubernetes deployment yang kompleks
* Cloud provider specific implementation yang terlalu sempit
* ORM specific business abstraction
* Event sourcing
* CQRS sebagai default architecture

Namun architecture harus memungkinkan fitur tersebut ditambahkan kemudian.

Service type `lambda` menyediakan adapter/entrypoint, bukan mengunci ke satu cloud vendor.

---

# 4. Design Principles

Architecture harus mengikuti prinsip berikut.

## 4.1 SOLID

### Single Responsibility

Setiap package memiliki tanggung jawab yang jelas.

Contoh:

```text
handler
    HTTP handling

service
    business logic

repository
    data persistence

middleware
    HTTP middleware

config
    configuration

logger
    logging
```

### Open / Closed

Implementasi repository dapat ditambahkan tanpa mengubah business logic.

Contoh:

```text
UserRepository
    PostgreSQLUserRepository
    MongoUserRepository
    MySQLUserRepository
    MSSQLUserRepository
```

### Liskov Substitution

Service tidak boleh mengetahui implementasi database.

Contoh:

```go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
}
```

Service hanya bergantung kepada interface tersebut.

### Interface Segregation

Hindari interface besar.

Buruk:

```go
type Repository interface {
    Create()
    Update()
    Delete()
    Find()
    Transaction()
    Cache()
    Upload()
}
```

Lebih baik:

```go
type UserReader interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserWriter interface {
    Create(ctx context.Context, user *User) error
}
```

### Dependency Inversion

Business layer tidak boleh bergantung langsung kepada:

```text
MongoDB
PostgreSQL
Redis
Gin
Zap
GORM
```

Sebaliknya:

```text
Domain / Application
        |
        v
    Interfaces
        |
        v
Infrastructure
```

---

# 5. High Level Architecture

Architecture utama berbasis shared core dengan multiple entrypoint:

```text
        +----------+   +----------+   +----------+   +----------+
        |   API    |   |   Cron   |   |  Worker  |   |  Lambda  |
        |  (Gin)   |   | Scheduler|   | Consumer |   | Adapter  |
        +----+-----+   +----+-----+   +----+-----+   +----+-----+
             |              |              |              |
             +------+-------+------+-------+------+-------+
                    |
                    v
             Application Service
                    |
                    v
               Domain Layer
                    |
                    v
            Repository Interface
                    |
          +---------+---------+---------+
          |         |         |         |
          v         v         v         v
      PostgreSQL  MySQL     MSSQL    MongoDB

Optional:

Application
     |
     v
 Cache Interface
     |
     v
   Redis
```

Service type hanya menentukan cara aplikasi di-trigger dan di-lifecycle-kan.
Business logic tetap berada di application/domain layer yang sama.

---

# 6. Recommended Architecture

Saya menyarankan menggunakan kombinasi:

```text
Clean Architecture
+
Hexagonal Architecture
+
Dependency Injection
+
Multi Entrypoint Service Types
```

Bukan Clean Architecture yang terlalu akademis, tetapi pragmatic clean architecture.

Struktur dependency:

```text
Entrypoint (api | cron | worker | lambda)
 |
 v
Transport Adapter (HTTP Handler / Job / Consumer / Lambda Handler)
 |
 v
Application Service
 |
 v
Domain
 |
 v
Interfaces
 |
 v
Infrastructure
```

Infrastructure tidak boleh naik ke domain.
Transport adapter tidak boleh berisi business logic.
---

# 7. Project Structure

Struktur awal:

```text
go-gin-boilerplate/
│
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── cron/
│   │   └── main.go
│   ├── worker/
│   │   └── main.go
│   └── lambda/
│       └── main.go
│
├── internal/
│   │
│   ├── app/
│   │   ├── bootstrap/
│   │   ├── server/          # API HTTP server
│   │   ├── cron/            # Cron runner
│   │   ├── worker/          # Worker runner
│   │   └── lambda/          # Lambda adapter
│   │
│   ├── config/
│   │
│   ├── domain/
│   │   ├── user/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   │
│   │   └── common/
│   │
│   ├── application/
│   │   ├── user/
│   │   │   ├── service.go
│   │   │   ├── command.go
│   │   │   └── query.go
│   │   │
│   │   └── common/
│   │
│   ├── infrastructure/
│   │   │
│   │   ├── database/
│   │   │   ├── postgres/
│   │   │   ├── mysql/
│   │   │   ├── mssql/
│   │   │   └── mongo/
│   │   │
│   │   ├── cache/
│   │   │   └── redis/
│   │   │
│   │   ├── logger/
│   │   ├── metrics/
│   │   └── tracing/
│   │
│   ├── interfaces/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   ├── request/
│   │   │   ├── response/
│   │   │   └── router/
│   │   │
│   │   ├── cron/
│   │   │   └── jobs/
│   │   │
│   │   ├── worker/
│   │   │   └── consumers/
│   │   │
│   │   └── lambda/
│   │       └── handlers/
│   │
│   └── shared/
│       ├── errors/
│       ├── pagination/
│       ├── validator/
│       ├── response/
│       └── utils/
│
├── migrations/
│
├── docs/
│
├── tests/
│   ├── integration/
│   └── e2e/
│
├── deployments/
│   ├── docker/
│   └── compose/
│
├── scripts/
│
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

`cmd/*` hanya sebagai entrypoint tipis.
Semua business logic harus reusable lintas service type melalui `internal/application` dan `internal/domain`.
---

# 8. Dependency Direction

Ini bagian yang sangat penting.

Dependency harus mengikuti arah:

```text
interfaces
      |
      v
application
      |
      v
domain
      ^
      |
infrastructure
```

Contoh:

```text
HTTP Handler
    |
    v
UserService
    |
    v
UserRepository interface
    ^
    |
PostgresUserRepository
```

Bukan:

```text
UserService
    |
    v
PostgreSQL
```

Dengan begitu database dapat diganti tanpa mengubah business logic.

---

# 9. Database Abstraction

Boilerplate harus mendukung beberapa database.

## Supported

```text
PostgreSQL
MySQL
MariaDB
Microsoft SQL Server
MongoDB
```

Tetapi jangan membuat satu interface database global seperti:

```go
type Database interface {
    Query()
    Insert()
    Update()
    Delete()
}
```

Karena SQL dan MongoDB memiliki paradigma berbeda.

Lebih baik abstraction berada pada repository/domain level.

Contoh:

```go
type UserRepository interface {
    FindByID(
        ctx context.Context,
        id string,
    ) (*User, error)

    Create(
        ctx context.Context,
        user *User,
    ) error
}
```

Kemudian:

```text
UserRepository
       |
       +---- PostgreSQL
       |
       +---- MySQL
       |
       +---- MSSQL
       |
       +---- MongoDB
```

---

# 10. Database Driver Strategy

Recommended:

```text
PostgreSQL
    pgx

MySQL / MariaDB
    go-sql-driver/mysql

MSSQL
    go-mssqldb

MongoDB
    official MongoDB Go Driver
```

Untuk SQL, gunakan `database/sql` atau driver native yang kompatibel dengan `database/sql` bila abstraction tersebut cukup.

Jangan menjadikan ORM sebagai dependency wajib boilerplate.

Jika project membutuhkan ORM, dapat ditambahkan sebagai adapter.

---

# 11. Database Configuration

Configuration harus menggunakan environment variable.

Contoh:

```env
APP_ENV=development
APP_NAME=my-service
APP_PORT=8080
SERVICE_TYPE=api
```

`SERVICE_TYPE` wajib salah satu dari:

```text
api
cron
worker
lambda
```

Database:

```env
DATABASE_DRIVER=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=mydb
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_SSLMODE=disable
```

Mongo:

```env
MONGO_ENABLED=false
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=mydb
```

Redis:

```env
REDIS_ENABLED=false
REDIS_URL=redis://localhost:6379
REDIS_DATABASE=0
```

---

# 12. Redis

Redis bersifat optional.

Jika:

```env
REDIS_ENABLED=false
```

aplikasi harus tetap berjalan tanpa Redis.

Redis dapat digunakan untuk:

```text
Caching
Session
Rate limiting
Distributed lock
Idempotency
Temporary data
Pub/Sub
Distributed coordination
```

Application layer hanya menggunakan interface.

Contoh:

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

Redis implementation:

```text
RedisCache
```

Jika Redis disabled:

```text
NoopCache
```

Dengan demikian application tidak perlu:

```go
if redis != nil
```

di seluruh codebase.

---

# 13. Configuration Management

Configuration harus centralized.

Contoh:

```go
type Config struct {
    App      AppConfig
    Service  ServiceConfig
    HTTP     HTTPConfig
    Cron     CronConfig
    Worker   WorkerConfig
    Lambda   LambdaConfig
    Database DatabaseConfig
    Mongo    MongoConfig
    Redis    RedisConfig
    Logger   LoggerConfig
}
```

`ServiceConfig` minimal:

```go
type ServiceConfig struct {
    Type string // api | cron | worker | lambda
}
```

Loader:

```text
Environment
      |
      v
Config Loader
      |
      v
Validated Config
      |
      v
Application
```

Configuration harus divalidasi saat startup.

Jika konfigurasi mandatory invalid:

```text
Application startup failed
```

Jangan membiarkan application berjalan dalam kondisi configuration invalid.

---

# 14. Logging

Logger harus structured.

Recommended abstraction:

```go
type Logger interface {
    Debug(ctx context.Context, msg string, fields ...Field)
    Info(ctx context.Context, msg string, fields ...Field)
    Warn(ctx context.Context, msg string, fields ...Field)
    Error(ctx context.Context, msg string, fields ...Field)
}
```

Log format:

```json
{
  "timestamp": "2026-10-02T10:00:00Z",
  "level": "info",
  "service": "user-service",
  "request_id": "abc123",
  "message": "user created"
}
```

Logger implementation dapat menggunakan:

```text
slog
```

atau logger production lain melalui adapter.

Saya menyarankan menggunakan Go standard library `log/slog` sebagai default supaya boilerplate tidak terlalu bergantung pada library eksternal.

---

# 15. Request ID

Setiap HTTP request harus memiliki:

```text
X-Request-ID
```

Flow:

```text
Client
   |
   v
Request ID Middleware
   |
   +---- existing request ID
   |
   +---- generate new ID
   |
   v
Handler
   |
   v
Service
   |
   v
Repository
```

Request ID harus masuk ke log context.

---

# 16. Correlation ID

Selain request ID, boilerplate dapat mendukung:

```text
X-Correlation-ID
```

Tujuannya untuk distributed system.

Contoh:

```text
API Gateway
    |
    | correlation-id: abc
    v
Service A
    |
    | correlation-id: abc
    v
Service B
    |
    | correlation-id: abc
    v
Service C
```

Hal ini sangat berguna jika boilerplate digunakan untuk microservices.

---

# 17. HTTP Middleware

Default middleware:

```text
Recovery
Request ID
Correlation ID
Logger
CORS
Security Headers
Request Timeout
Body Limit
Rate Limit optional
Authentication optional
```

Middleware harus dapat diaktifkan atau dinonaktifkan melalui configuration.

---

# 18. Error Handling

Semua error harus memiliki standard format.

Contoh:

```json
{
  "success": false,
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User not found",
    "request_id": "abc123"
  }
}
```

Internal error tidak boleh membocorkan:

```text
SQL query
Database credentials
Stack trace
Internal hostname
Redis connection information
File path
Infrastructure details
```

Development environment dapat menampilkan informasi lebih detail.

Production harus menggunakan sanitized error.

---

# 19. Error Classification

Error dibagi menjadi:

```text
Domain Error
Application Error
Infrastructure Error
Validation Error
Authentication Error
Authorization Error
Not Found Error
Conflict Error
Timeout Error
Internal Error
```

Contoh:

```go
var ErrUserNotFound = errors.New("user not found")
```

HTTP handler kemudian memetakan:

```text
ErrUserNotFound
      |
      v
HTTP 404
```

Database error:

```text
Database timeout
      |
      v
HTTP 503
```

---

# 20. API Response Standard

Success:

```json
{
  "success": true,
  "data": {}
}
```

List:

```json
{
  "success": true,
  "data": [],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

Error:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request"
  }
}
```

---

# 21. Validation

Request validation harus centralized.

Contoh:

```go
type CreateUserRequest struct {
    Name  string `validate:"required,min=3,max=100"`
    Email string `validate:"required,email"`
}
```

Response validation tidak harus dipaksakan pada semua endpoint.

---

# 22. Graceful Shutdown

Application harus menangani:

```text
SIGINT
SIGTERM
```

Shutdown flow menyesuaikan service type.

## API

```text
Signal
  |
  v
Stop accepting requests
  |
  v
Wait active requests
  |
  v
Close Redis
  |
  v
Close MongoDB
  |
  v
Close SQL connection
  |
  v
Flush logger
  |
  v
Exit
```

## Worker

```text
Signal
  |
  v
Stop consuming new jobs
  |
  v
Wait in-flight jobs
  |
  v
Close dependencies
  |
  v
Exit
```

## Cron

```text
Signal
  |
  v
Stop scheduler
  |
  v
Wait running job
  |
  v
Close dependencies
  |
  v
Exit
```

## Lambda

Lambda mengikuti lifecycle runtime provider.
Setiap invocation harus selesai dengan context timeout yang jelas dan resource dibersihkan setelah handler selesai.

Dengan timeout:

```text
SHUTDOWN_TIMEOUT=30s
```

---

# 23. Health Endpoints

Health endpoints wajib untuk service type `api`.

Minimal:

```text
GET /health
GET /health/live
GET /health/ready
```

## Liveness

Memastikan process masih hidup.

```json
{
  "status": "ok"
}
```

## Readiness

Memastikan service siap menerima traffic.

Contoh check:

```text
Application
    |
    +---- Database
    |
    +---- Redis optional
    |
    +---- Required dependencies
```

Jika database unavailable:

```text
/health/live  -> 200
/health/ready -> 503
```

Ini penting untuk Docker dan Kubernetes.

Untuk `cron` dan `worker`, health check dapat disediakan secara opsional melalui:
- internal HTTP probe ringan, atau
- platform-native health/liveness mechanism

Untuk `lambda`, health mengikuti lifecycle dan monitoring platform serverless.

---

# 24. Metrics

Boilerplate sebaiknya metrics-ready.

Minimal metrics:

```text
HTTP request count
HTTP request duration
HTTP response status
Active requests
Database query duration
Database errors
Cache hit
Cache miss
```

Endpoint:

```text
GET /metrics
```

Format dapat mengikuti Prometheus.

---

# 25. Distributed Tracing

Tracing bersifat optional tetapi architecture harus siap.

Contoh:

```text
HTTP Request
     |
     v
Trace
     |
     +---- Handler
     |
     +---- Service
     |
     +---- Database
     |
     +---- Redis
     |
     +---- External API
```

OpenTelemetry dapat digunakan sebagai adapter.

---

# 26. Authentication

Authentication tidak harus menjadi bagian mandatory.

Tetapi middleware interface harus disediakan.

Contoh:

```go
type Authenticator interface {
    Authenticate(ctx context.Context, token string) (*Principal, error)
}
```

Kemudian implementation dapat berupa:

```text
JWT
OAuth2
API Key
Session
External Identity Provider
```

---

# 27. Authorization

Authorization sebaiknya dipisahkan dari authentication.

```text
Authentication
    |
    v
Who are you?

Authorization
    |
    v
What are you allowed to do?
```

Interface:

```go
type Authorizer interface {
    Can(
        ctx context.Context,
        principal *Principal,
        action string,
        resource string,
    ) error
}
```

---

# 28. Transaction Management

Untuk database SQL, boilerplate harus memiliki transaction abstraction.

Contoh:

```go
type TransactionManager interface {
    WithinTransaction(
        ctx context.Context,
        fn func(ctx context.Context) error,
    ) error
}
```

Application service tidak boleh mengetahui detail:

```text
BEGIN
COMMIT
ROLLBACK
```

---

# 29. Database Connection Pool

SQL database harus mendukung:

```text
MaxOpenConnections
MaxIdleConnections
ConnectionMaxLifetime
ConnectionMaxIdleTime
```

Contoh:

```env
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m
DB_CONN_MAX_IDLE_TIME=10m
```

Nilai default harus conservative dan configurable.

---

# 30. MongoDB

MongoDB client harus mendukung:

```text
Connection pooling
Timeout
Ping
Graceful shutdown
Read preference
Write concern
```

Configuration:

```env
MONGO_URI=
MONGO_DATABASE=
MONGO_MAX_POOL_SIZE=
MONGO_MIN_POOL_SIZE=
MONGO_CONNECT_TIMEOUT=
```

---

# 31. Repository Pattern

Contoh:

```go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
}
```

Implementation:

```text
PostgresUserRepository
MySQLUserRepository
MSSQLUserRepository
MongoUserRepository
```

Business logic tidak berubah.

---

# 32. Service Layer

Service bertanggung jawab terhadap business rules.

Contoh:

```go
type UserService struct {
    repository UserRepository
    cache      Cache
}
```

Service tidak boleh melakukan:

```text
gin.Context parsing
HTTP response
SQL query langsung
Mongo query langsung
Redis command langsung
```

---

# 33. Handler Layer

Handler bertanggung jawab:

```text
HTTP request
    |
    v
Parse
    |
    v
Validate
    |
    v
Call service
    |
    v
Map response
```

Handler tidak boleh memiliki business logic.

---

# 34. Dependency Injection

Semua dependency diinisialisasi saat bootstrap.

Contoh:

```text
main
 |
 v
Config
 |
 v
Logger
 |
 v
Database
 |
 v
Redis
 |
 v
Repositories
 |
 v
Services
 |
 v
Handlers
 |
 v
Router
 |
 v
HTTP Server
```

Hindari global singleton seperti:

```go
var DB *sql.DB
var Redis *redis.Client
var Logger *slog.Logger
```

Dependency harus explicit.

---

# 35. Bootstrap

Bootstrap menjadi composition root.

Contoh:

```text
cmd/{api|cron|worker|lambda}/main.go
        |
        v
bootstrap.NewApplication(cfg)
        |
        +---- Config
        +---- Logger
        +---- Database
        +---- Redis
        +---- Repository
        +---- Service
        |
        +---- by SERVICE_TYPE:
               api    -> Handler + Router + HTTP Server
               cron   -> Jobs + Scheduler Runner
               worker -> Consumers + Worker Runner
               lambda -> Lambda Handlers + Adapter
```

Bootstrap mereuse dependency yang sama.
Yang berbeda hanya runner / adapter sesuai service type.

Ini membuat dependency graph mudah dipahami dan dites.

---

# 36. Testing Strategy

Target testing:

```text
Unit Test
Integration Test
E2E Test
```

## Unit Test

Target:

```text
Domain
Service
Use Case
Utility
Validation
```

Tidak membutuhkan database asli.

## Integration Test

Menggunakan:

```text
PostgreSQL
MySQL
MSSQL
MongoDB
Redis
```

Docker dapat digunakan untuk test dependency.

## E2E Test

Flow:

```text
HTTP
 |
 v
Handler
 |
 v
Service
 |
 v
Repository
 |
 v
Database
```

---

# 37. Testability

Service harus dapat dites dengan mock/fake repository.

Contoh:

```text
UserService
    |
    v
UserRepository interface
    |
    +---- PostgreSQL implementation
    |
    +---- Mock implementation
```

Dengan demikian unit test tidak membutuhkan database.

---

# 38. Docker

Harus tersedia:

```text
Dockerfile
docker-compose.yml
```

Docker image harus menggunakan multi-stage build.

Flow:

```text
Builder
    |
    v
Go binary (api | cron | worker | lambda)
    |
    v
Minimal runtime image
```

Build target harus mendukung service type, misalnya:

```text
docker build --build-arg SERVICE_TYPE=api
docker build --build-arg SERVICE_TYPE=cron
docker build --build-arg SERVICE_TYPE=worker
docker build --build-arg SERVICE_TYPE=lambda
```

Container harus menjalankan application sebagai non-root user.

---

# 39. Docker Compose

Development compose minimal:

```text
API
Cron optional
Worker optional
PostgreSQL
MongoDB
Redis
```

Database yang tidak digunakan sebaiknya dapat dimatikan.

Contoh:

```text
docker compose up api postgres
```

atau:

```text
docker compose up api worker cron mongo redis
```

---

# 40. Environment

Minimal environment:

```text
development
staging
production
test
```

Configuration harus berasal dari environment.

Jangan menyimpan credential di repository.

---

# 41. Security Baseline

Boilerplate harus menyediakan default security baseline:

```text
HTTPS-ready
CORS configuration
Request body limit
Timeout
Security headers
Input validation
SQL parameterization
Secret melalui environment
No sensitive logs
Graceful error response
Non-root Docker
```

Password atau token tidak boleh muncul di log.

---

# 42. API Versioning

Default API:

```text
/api/v1
```

Contoh:

```text
GET /api/v1/users
POST /api/v1/users
GET /api/v1/users/:id
```

Versioning harus berada pada HTTP layer.

Domain tidak mengetahui API version.

---

# 43. Pagination

Standard pagination:

```text
page
limit
```

Contoh:

```text
GET /api/v1/users?page=1&limit=20
```

Untuk dataset besar, architecture juga harus memungkinkan cursor pagination.

Contoh:

```text
GET /api/v1/users?cursor=abc&limit=20
```

---

# 44. Idempotency

Boilerplate sebaiknya menyediakan optional idempotency middleware.

Header:

```text
Idempotency-Key
```

Useful untuk:

```text
Payment
Order
Booking
Transaction
Webhook
```

Redis dapat menjadi backend idempotency jika enabled.

---

# 45. Rate Limiting

Rate limiting optional.

Architecture:

```text
RateLimiter interface
        |
        +---- InMemoryRateLimiter
        |
        +---- RedisRateLimiter
```

Untuk single instance:

```text
InMemory
```

Untuk multiple instances:

```text
Redis
```

---

# 46. Service Type Support

Boilerplate wajib mendukung 4 service type sebagai first-class citizen.

```text
api
cron
worker
lambda
```

Pemilihan service type dilakukan saat runtime melalui:

```env
SERVICE_TYPE=api
```

atau melalui binary entrypoint di `cmd/`.

Domain dan application layer harus independent dari service type.
Hanya entrypoint, transport, dan lifecycle yang berbeda.

## 46.1 API

Service type default untuk HTTP REST API.

Karakteristik:

```text
Gin HTTP Server
Middleware
Handler
Router
Request validation
Standard API response
Health / readiness / liveness
Graceful shutdown
```

Flow:

```text
HTTP Request
  -> Middleware
  -> Handler
  -> Application Service
  -> Domain / Repository
```

Entrypoint:

```text
cmd/api/main.go
```

## 46.2 Cron

Service type untuk scheduled / periodic job.

Karakteristik:

```text
Scheduler / ticker / crontab expression
Job registration
Single-run atau recurring
Distributed lock optional untuk multi-instance
Timeout per job
Structured logging per job execution
Graceful stop menunggu job yang sedang berjalan
```

Abstraction:

```go
type Job interface {
    Name() string
    Run(ctx context.Context) error
}

type Scheduler interface {
    Register(spec string, job Job) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

Contoh job:

```text
CleanupExpiredSessions
GenerateDailyReport
SyncExternalData
PurgeOldLogs
```

Flow:

```text
Scheduler Trigger
  -> Job Adapter
  -> Application Service
  -> Domain / Repository
```

Entrypoint:

```text
cmd/cron/main.go
```

## 46.3 Worker

Service type untuk background processing / queue consumer.

Karakteristik:

```text
Long-running process
Consume message / event / queue
Concurrency control
Retry / dead-letter ready
Ack / nack semantics
In-flight job tracking
Graceful shutdown
```

Abstraction:

```go
type Worker interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type Consumer interface {
    Handle(ctx context.Context, message Message) error
}
```

Contoh worker:

```text
Email Worker
Webhook Worker
Queue Consumer
Notification Worker
Image Processing Worker
```

Flow:

```text
Queue / Broker / Event Source
  -> Consumer Adapter
  -> Application Service
  -> Domain / Repository
```

Message broker implementation spesifik (Kafka, RabbitMQ, SQS, dll) tidak wajib di versi pertama,
tetapi abstraction worker/consumer harus tersedia agar adapter mudah ditambahkan.

Entrypoint:

```text
cmd/worker/main.go
```

## 46.4 Lambda

Service type untuk serverless / function-style execution.

Karakteristik:

```text
Stateless per invocation
Cold start aware
Context timeout dari runtime
Event-driven adapter
Minimal bootstrap yang reusable
Tidak bergantung pada long-running HTTP server
```

Abstraction:

```go
type LambdaHandler interface {
    Handle(ctx context.Context, event any) (any, error)
}
```

Contoh event:

```text
HTTP API Gateway event
SQS event
EventBridge / cron event
Custom JSON event
```

Flow:

```text
Cloud Event
  -> Lambda Adapter
  -> Application Service
  -> Domain / Repository
```

Entrypoint:

```text
cmd/lambda/main.go
```

Lambda adapter boleh provider-agnostic di layer interface,
dengan implementasi awal yang siap diarahkan ke AWS Lambda atau runtime sejenis.

## 46.5 Shared Rules Across Service Types

```text
1. Business logic tidak boleh hidup di cmd/ atau transport adapter.
2. Satu use case harus reusable dari api, cron, worker, maupun lambda.
3. Dependency injection tetap melalui bootstrap.
4. Logging, metrics, tracing, config, dan error handling harus konsisten.
5. Service type tidak mengubah domain model.
6. Resource lifecycle harus jelas untuk setiap type.
7. Konfigurasi yang tidak relevan per type boleh diabaikan, tetapi tidak boleh mematahkan startup bila optional.
```

Contoh:

```text
CreateUser use case
  <- dipanggil dari API handler
  <- dipanggil dari Worker consumer (bulk import)
  <- dipanggil dari Cron job (sync)
  <- dipanggil dari Lambda handler (event create)
```

---

# 47. External Service Client

Untuk external API, gunakan interface:

```go
type PaymentClient interface {
    CreatePayment(
        ctx context.Context,
        request CreatePaymentRequest,
    ) (*PaymentResponse, error)
}
```

Infrastructure:

```text
HTTP Payment Client
```

Service hanya mengetahui interface.

HTTP client harus memiliki:

```text
Timeout
Retry optional
Circuit breaker optional
Request ID
Tracing
Structured logging
```

---

# 48. Reliability

Boilerplate harus menerapkan:

```text
Timeout
Retry
Graceful shutdown
Connection pooling
Health check
Readiness check
Error classification
Idempotency
Circuit breaker ready
Observability
```

Retry tidak boleh diterapkan secara sembarangan.

Retry hanya untuk operation yang memang aman atau idempotent.

---

# 49. Scalability

Application harus stateless secara default.

Contoh:

```text
                Load Balancer
                     |
       +-------------+-------------+
       |             |             |
       v             v             v
    API 1          API 2          API 3
       |             |             |
       +-------------+-------------+
                     |
              Shared Database
                     |
                 Redis
```

Tidak boleh menyimpan session penting hanya di memory process jika aplikasi akan dijalankan multi-instance.

---

# 50. Observability Architecture

Target:

```text
                 Application
                      |
          +-----------+-----------+
          |           |           |
          v           v           v
        Logs       Metrics      Traces
          |           |           |
          +-----------+-----------+
                      |
                Observability
                  Platform
```

Boilerplate harus dapat diintegrasikan dengan:

```text
Prometheus
Grafana
OpenTelemetry
ELK
Loki
Jaeger
```

tanpa membuat semua dependency menjadi mandatory.

---

# 51. Documentation

README harus menjelaskan:

```text
Installation
Configuration
Service Types (api, cron, worker, lambda)
Project Structure
Architecture
Database Setup
Redis Setup
Running locally
Running tests
Docker
API development
Cron job development
Worker development
Lambda adapter development
Adding repository
Adding service
Adding handler
Adding middleware
Logging
Health check
Metrics
Deployment
```

---

# 52. Developer Experience

Harus tersedia Makefile.

Contoh:

```text
make dev
make build
make build-api
make build-cron
make build-worker
make build-lambda
make run-api
make run-cron
make run-worker
make run-lambda
make test
make test-unit
make test-integration
make lint
make fmt
make vet
make docker-build
make docker-up
make docker-down
make migrate-up
make migrate-down
```

---

# 53. Code Quality

CI harus melakukan:

```text
go fmt
go vet
go test
go test -race
static analysis
lint
build
```

Pipeline:

```text
Pull Request
      |
      v
Format Check
      |
      v
Static Analysis
      |
      v
Unit Test
      |
      v
Race Test
      |
      v
Build
```

---

# 54. Definition of Done

Boilerplate dianggap selesai apabila:

```text
[ ] Go 1.27.0
[ ] Multi service type: api
[ ] Multi service type: cron
[ ] Multi service type: worker
[ ] Multi service type: lambda
[ ] Gin integrated for api
[ ] Shared application/domain lintas service type
[ ] Clean architecture
[ ] SOLID structure
[ ] Dependency injection
[ ] PostgreSQL adapter
[ ] MySQL/MariaDB adapter
[ ] MSSQL adapter
[ ] MongoDB adapter
[ ] Redis optional
[ ] Structured logger
[ ] Request ID
[ ] Correlation ID
[ ] Central error handling
[ ] Validation
[ ] Standard response
[ ] Health endpoint
[ ] Readiness endpoint
[ ] Graceful shutdown per service type
[ ] Connection pooling
[ ] Transaction abstraction
[ ] Unit tests
[ ] Integration tests
[ ] Docker
[ ] Docker Compose
[ ] Makefile targets per service type
[ ] Environment configuration
[ ] SERVICE_TYPE configuration
[ ] Metrics ready
[ ] OpenTelemetry ready
[ ] API versioning
[ ] Pagination
[ ] Security baseline
[ ] Documentation
[ ] CI pipeline
```

# 55. Contoh Flow Request

Untuk endpoint:

```text
POST /api/v1/users
```

flow:

```text
Client
  |
  v
Gin
  |
  v
Request ID Middleware
  |
  v
Recovery Middleware
  |
  v
Logger Middleware
  |
  v
Validation
  |
  v
User Handler
  |
  v
User Service
  |
  v
UserRepository interface
  |
  v
PostgreSQL Repository
  |
  v
PostgreSQL
```

Jika nanti pindah ke MongoDB:

```text
User Handler
      |
      v
User Service
      |
      v
UserRepository interface
      |
      v
MongoDB Repository
      |
      v
MongoDB
```

Handler dan business logic tidak perlu berubah.

---

# 56. Prinsip Penting untuk Boilerplate Ini

Saya akan menetapkan beberapa aturan arsitektur sebagai architectural rules.

```text
1. Handler / Job / Consumer / Lambda adapter tidak boleh mengakses database secara langsung.

2. Service tidak boleh mengetahui Gin, scheduler, queue, atau Lambda runtime.

3. Domain tidak boleh mengetahui database.

4. Domain tidak boleh mengetahui Redis.

5. Domain tidak boleh mengetahui HTTP / cron / worker / lambda transport.

6. Infrastructure tidak boleh menentukan business rule.

7. Global mutable state harus dihindari.

8. Semua external dependency harus melalui interface bila digunakan oleh application/domain.

9. Redis harus optional.

10. Database implementation harus interchangeable.

11. Application harus stateless.

12. Semua request / job / event memiliki request ID atau execution ID.

13. Semua error harus dapat dipetakan secara konsisten.

14. Semua resource harus memiliki lifecycle yang jelas.

15. Semua goroutine harus memiliki mekanisme shutdown.

16. Configuration harus divalidasi saat startup.

17. Production default harus secure.

18. Logging tidak boleh mengandung secret.

19. Retry hanya untuk operation yang aman untuk diulang.

20. Business logic harus dapat di-unit-test tanpa infrastructure.

21. Satu use case harus reusable lintas service type `api`, `cron`, `worker`, dan `lambda`.

22. Service type hanya menentukan entrypoint dan lifecycle, bukan business rule.
```

# 57. Target Architecture Final

Secara keseluruhan, blueprint-nya:

```text
        +----------+  +----------+  +----------+  +----------+
        |   API    |  |   Cron   |  |  Worker  |  |  Lambda  |
        |   Gin    |  | Scheduler|  | Consumer |  | Adapter  |
        +----+-----+  +----+-----+  +----+-----+  +----+-----+
             |             |             |             |
             +------+------+------+------+------+------+
                    |
                    v
              Transport Adapter
                    |
                    v
             Application Service
                    |
                    v
                  Domain
                    |
       +------------+------------+
       |            |            |
       v            v            v
  Repository      Cache      External API
   Interface     Interface     Interface
       |            |            |
+------+------+     |       +----+----+
|      |      |     |       |         |
v      v      v     v       v         v
PG   MySQL  MSSQL Redis   HTTP      gRPC
 |
 +----------------+
                  |
                  v
               MongoDB

Observability:

Application
     |
     +---- Logs
     +---- Metrics
     +---- Traces
```

## 58. Rekomendasi Stack Default

Untuk baseline ini saya akan memilih:

```text
Language       Go 1.27.0
Service Types  api | cron | worker | lambda
HTTP (api)     Gin
Scheduler      lightweight cron/ticker abstraction
Worker         consumer interface + adapter-ready
Lambda         function adapter (provider-agnostic)
Config         Environment based configuration
Logger         log/slog
Validation     go-playground/validator
SQL            database/sql + native drivers
PostgreSQL     pgx
MySQL/MariaDB  go-sql-driver/mysql
MSSQL          go-mssqldb
MongoDB        official MongoDB Go Driver
Cache          Redis
Metrics        Prometheus compatible
Tracing        OpenTelemetry
Testing        testing + httptest
Container      Docker
Orchestration  Docker Compose
CI             GitHub Actions
API Docs       OpenAPI
```

Yang paling penting, saya tidak akan membuat boilerplate ini bergantung pada satu ORM. Untuk base template yang dimaksudkan reusable lintas PostgreSQL, MySQL/MariaDB, MSSQL, dan MongoDB, abstraction sebaiknya berada di repository/application boundary, bukan memaksakan abstraction database yang sama untuk semua engine.

Dengan desain ini, boilerplate bisa menjadi foundation untuk REST API, cron job, background worker, serverless/lambda function, backend OTA, payment service, booking service, internal enterprise service, sampai microservice tanpa harus membongkar architecture ketika kebutuhan berkembang.
