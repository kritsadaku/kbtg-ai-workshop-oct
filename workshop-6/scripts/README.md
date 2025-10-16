# Database Scripts

Scripts สำหรับจัดการข้อมูลใน SQLite database ของ temp-backend project

## Files

- `seed_users.sh` - Script สำหรับ seed ข้อมูล users จาก CSV file เข้า database
- `count_users.sh` - Script สำหรับนับจำนวนและแสดงสถิติของ users ใน database
- `users_data.csv` - ไฟล์ CSV ตัวอย่างที่มีข้อมูล users สำหรับ testing

## การใช้งาน

### 1. Seed Data from CSV

```bash
# ใช้ไฟล์ default
cd scripts
./seed_users.sh

# กำหนดไฟล์เอง
./seed_users.sh custom_users.csv ../users.db
```

**Features:**
- ตรวจสอบว่าไฟล์ CSV และ database มีอยู่จริง
- ใช้ INSERT OR IGNORE เพื่อป้องกันข้อมูลซ้ำ
- แสดงสถิติก่อนและหลัง seed
- จัดการ transaction เพื่อความปลอดภัย
- รองรับอักขระไทย

### 2. Count Users and Statistics

```bash
# ใช้ database default
cd scripts
./count_users.sh

# กำหนด database file เอง
./count_users.sh /path/to/database.db
```

**Features:**
- นับจำนวน users ทั้งหมด
- แสดงสถิติแยกตาม membership level
- สถิติ points (รวม, เฉลี่ย, สูงสุด, ต่ำสุด)
- Timeline การสมัครสมาชิก
- Top 10 users ที่มี points มากที่สุด
- ข้อมูลสมาชิกใหม่ล่าสุด 5 คน
- ขนาดไฟล์ database และจำนวน records

## CSV Format

ไฟล์ CSV ต้องมี header และ columns ดังนี้:

```csv
first_name,last_name,phone,email,member_since,membership_level,points
สม,ใจ,0812345678,somjai@example.com,2023-01-15,Bronze,1000
```

**Columns:**
- `first_name` - ชื่อ (ไม่เกิน 3 ตัวอักษร ตาม business rules)
- `last_name` - นามสกุล (ไม่เกิน 3 ตัวอักษร ตาม business rules)
- `phone` - เบอร์โทรศัพท์
- `email` - อีเมล (ต้องไม่ซ้ำ)
- `member_since` - วันที่สมัครสมาชิก (YYYY-MM-DD)
- `membership_level` - ระดับสมาชิก (Bronze, Silver, Gold, Platinum)
- `points` - จำนวน points เริ่มต้น

## Requirements

- SQLite3 installed
- Database ต้องถูกสร้างแล้วโดยการรัน Go application
- macOS หรือ Linux environment

## Installation

### macOS
```bash
brew install sqlite
```

### Ubuntu/Debian
```bash
sudo apt-get install sqlite3
```

## Examples

### เตรียม database
```bash
# รัน Go application เพื่อสร้าง database และ tables
cd /path/to/temp-backend
go run main.go
# กด Ctrl+C เพื่อหยุด server หลังจาก database ถูกสร้างแล้ว
```

### Seed sample data
```bash
cd scripts
./seed_users.sh
```

### ดูสถิติ
```bash
./count_users.sh
```

## Error Handling

Scripts จะตรวจสอบและแสดง error message ที่เข้าใจง่าย:
- ไฟล์ CSV หรือ database ไม่พบ
- SQLite3 ไม่ได้ติดตั้ง
- ข้อผิดพลาดในการ insert ข้อมูล
- Database corruption หรือ permission issues

## Notes

- Scripts ใช้ `INSERT OR IGNORE` เพื่อป้องกันข้อมูลซ้ำ (ตาม unique email constraint)
- รองรับอักขระ Unicode/Thai
- มีการจัดการ single quotes ในข้อมูลเพื่อป้องกัน SQL injection
- ใช้ transaction เพื่อความปลอดภัยของข้อมูล