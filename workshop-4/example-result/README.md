# User Management API with Transfer Feature

Go backend สำหรับจัดการข้อมูลผู้ใช้และการโอนแต้ม ด้วย Fiber framework, SQLite database และ Clean Architecture

## Features

- ✅ CRUD operations สำหรับผู้ใช้
- ✅ **Point Transfer System** - โอนแต้มระหว่างผู้ใช้
- ✅ **Point Ledger** - บันทึกประวัติการเปลี่ยนแปลงแต้ม
- ✅ **Idempotency Support** - ป้องกันการโอนซ้ำ
- ✅ **Atomic Transactions** - ความปลอดภัยของข้อมูล
- ✅ SQLite database with proper indexes
- ✅ CORS support
- ✅ JSON API
- ✅ Input validation
- ✅ Error handling
- ✅ Clean Architecture

## API Endpoints

### Health Check
```bash
curl http://localhost:3000/
```

### User Management
```bash
# Get all users
curl http://localhost:3000/users

# Get user by ID
curl http://localhost:3000/users/1

# Create user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"first_name":"สมชาย","last_name":"ใจดี","phone":"081-234-5678","email":"somchai@example.com","membership_level":"Gold","points":15420}'

# Update user
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"points":20000,"membership_level":"Platinum"}'

# Delete user
curl -X DELETE http://localhost:3000/users/1
```

### 🆕 Transfer Operations

#### Create Transfer
```bash
curl -X POST http://localhost:3000/transfers \
  -H "Content-Type: application/json" \
  -d '{
    "fromUserId": 1,
    "toUserId": 2,
    "amount": 1000,
    "note": "ขอบคุณสำหรับความช่วยเหลือ"
  }'
```

#### Get Transfer by ID (Idempotency Key)
```bash
curl http://localhost:3000/transfers/5d1f8c7a-2b5b-4b1f-9f2a-8f50b0a8d9f3
```

#### Get Transfer History
```bash
# Get transfers for specific user with pagination
curl "http://localhost:3000/transfers?userId=1&page=1&pageSize=20"
```

## Data Models

### User
```json
{
  "id": 1,
  "first_name": "สมชาย",
  "last_name": "ใจดี",
  "phone": "081-234-5678",
  "email": "somchai@example.com",
  "member_since": "21/10/2025",
  "membership_level": "Gold",
  "points": 15420,
  "created_at": "2025-10-16T21:12:13Z",
  "updated_at": "2025-10-16T21:12:13Z"
}
```

### Transfer
```json
{
  "transferId": 1,
  "idemKey": "5d1f8c7a-2b5b-4b1f-9f2a-8f50b0a8d9f3",
  "fromUserId": 1,
  "toUserId": 2,
  "amount": 1000,
  "status": "completed",
  "note": "ขอบคุณสำหรับความช่วยเหลือ",
  "createdAt": "2025-10-16T14:03:12Z",
  "updatedAt": "2025-10-16T14:03:12Z",
  "completedAt": "2025-10-16T14:03:12Z"
}
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  phone TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  member_since TEXT NOT NULL,
  membership_level TEXT NOT NULL DEFAULT 'Bronze',
  points INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
```

### Transfers Table
```sql
CREATE TABLE transfers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  from_user_id INTEGER NOT NULL,
  to_user_id INTEGER NOT NULL,
  amount INTEGER NOT NULL CHECK (amount > 0),
  status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
  note TEXT,
  idempotency_key TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  completed_at TEXT,
  fail_reason TEXT,
  FOREIGN KEY (from_user_id) REFERENCES users(id),
  FOREIGN KEY (to_user_id) REFERENCES users(id)
);
```

### Point Ledger Table
```sql
CREATE TABLE point_ledger (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  change INTEGER NOT NULL,
  balance_after INTEGER NOT NULL,
  event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
  transfer_id INTEGER,
  reference TEXT,
  metadata TEXT,
  created_at TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (transfer_id) REFERENCES transfers(id)
);
```

## Transfer Features

### 🔒 **Atomic Transactions**
- การโอนแต้มใช้ database transaction
- รับประกันความสอดคล้องของข้อมูล
- Rollback อัตโนมัติเมื่อเกิดข้อผิดพลาด

### 🔑 **Idempotency Support**
- ใช้ idempotency key เพื่อป้องกันการโอนซ้ำ
- Key จะถูกสร้างอัตโนมัติและส่งกลับใน response
- สามารถตรวจสอบสถานะการโอนด้วย key ได้

### 📊 **Point Ledger**
- บันทึกประวัติการเปลี่ยนแปลงแต้มทุกครั้ง
- Append-only สำหรับการตรวจสอบ (audit trail)
- รองรับ event types หลากหลาย

### ✅ **Business Rules**
- ไม่สามารถโอนให้ตัวเองได้
- ตรวจสอบแต้มเพียงพอก่อนโอน
- อัปเดตแต้มทั้งผู้ส่งและผู้รับพร้อมกัน

### 🛡️ **Error Handling**
- การตรวจสอบ input อย่างครอบคลัด
- Error codes ตามมาตรฐาน HTTP
- Message ภาษาไทยที่เข้าใจง่าย

## การรันโปรแกรม

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Build the application**:
   ```bash
   go build .
   ```

3. **Run the server**:
   ```bash
   go run main.go
   ```

4. **Test APIs**:
   ```bash
   # Test user APIs
   ./test_api.sh
   
   # Test transfer APIs
   ./test_transfer_api.sh
   ```

เซิร์ฟเวอร์จะทำงานที่ `http://localhost:3000`

## Architecture

โปรเจกต์ใช้ **Clean Architecture** pattern:

- **Domain Layer**: Business entities และ interfaces
- **Usecase Layer**: Application business rules
- **Repository Layer**: Data access implementations  
- **Handler Layer**: HTTP presentation layer
- **Infrastructure**: Database, web server

## Dependencies

- **github.com/gofiber/fiber/v2** - Web framework
- **github.com/mattn/go-sqlite3** - SQLite driver
- **github.com/gofiber/fiber/v2/middleware/cors** - CORS middleware

## Version History

- **v2.1.0** - เพิ่ม Transfer Feature พร้อม Point Ledger
- **v2.0.0** - Clean Architecture refactoring
- **v1.0.0** - Basic User CRUD

This is a simple Go backend application using the Fiber web framework.

## Prerequisites

- Go 1.17 or higher (currently using Go 1.23.1)

## Installation

The Fiber library is already installed in this project. If you need to install it again:

```bash
go get github.com/gofiber/fiber/v2
```

## Running the Application

To run the server:

```bash
go run main.go
```

The server will start on `http://localhost:3000`

## Available Routes

- `GET /` - Returns "Hello, World!"
- `GET /:value` - Returns the value parameter (e.g., `/hello` returns "value: hello")
- `GET /user/:name?` - Optional name parameter (e.g., `/user/john` returns "Hello john")
- `GET /api/*` - Wildcard route for API paths (e.g., `/api/user/john` returns "API path: user/john")

## Documentation

For more information, visit the official Fiber documentation: https://docs.gofiber.io/