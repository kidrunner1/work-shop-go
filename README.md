# Go Fiber Workshop

ตัวอย่างพื้นฐานสำหรับข้อ 5.0–6 ใช้ Fiber v2 และ go-playground/validator/v10

## โครงสร้าง

```text
workshop-go/
├── main.go                  เริ่ม server ที่ port 3000
├── routes/
│   ├── routes.go            Basic Auth, group และเส้นทาง API
│   └── routes_test.go       ทดสอบ API
├── controller/
│   ├── factorial.go         ข้อ 5.1
│   ├── ascii.go             ข้อ 5.2
│   └── register.go          ข้อ 6 อ่าน body และตรวจสอบข้อมูล
└── models/
    └── register.go          ฟิลด์สมัครสมาชิกและกฎ validate
```

ลำดับการทำงาน: request → Basic Auth → route → controller → response

## วิธีรัน

ใช้ Go 1.26 ขึ้นไปตาม dependency ที่ติดตั้ง เปิด terminal ในโฟลเดอร์โปรเจกต์ แล้วรัน:

```sh
go mod download
go run .
```

หยุด server ด้วย Ctrl+C ส่วนการทดสอบอัตโนมัติใช้ `go test ./...`

## ข้อ 5.0 Basic Auth

ทุก API รวมทั้งสมัครสมาชิก ต้องตั้งค่าใน Postman ที่แท็บ Authorization:

- Type: Basic Auth
- Username: `gofiber`
- Password: `21022566`

ถ้าไม่ส่งหรือส่งผิด จะได้ HTTP 401

## ข้อ 5.1 Factorial

`GET http://localhost:3000/api/v1/fact/5`

```json
{"number":5,"factorial":120}
```

ใช้ for คูณเลขตั้งแต่ 1 ถึงจำนวนที่รับมา รับจำนวนเต็ม 0–20 เพราะ 21! เกินขนาด int64 โดย 0! = 1 ถ้าข้อมูลผิดจะได้ HTTP 400

## ข้อ 5.2 Query Param เป็น ASCII

`GET http://localhost:3000/api/v3/joe?tax_id=Hello`

```json
{"tax_id":"Hello","ascii":"72 101 108 108 111"}
```

ใช้ `app.Group("/api/v3")` ใน routes/routes.go เปลี่ยน `/joe` เป็นชื่อเล่นของตัวเองได้ที่นั่น

รับตัวอักษรภาษาอังกฤษ ตัวเลข และอักขระ ASCII หากไม่ระบุ tax_id หรือใส่อักขระนอก ASCII เช่นภาษาไทย จะได้ HTTP 400 (ภาษาไทยมีรหัส Unicode ไม่ใช่ ASCII)

## ข้อ 5.3 ชื่อย่อ controller

ใน routes/routes.go ใช้ import alias:

```go
import c "workshop-go/controller"
```

จึงเรียก `c.TestParams` หรือ `c.Register` ได้

## ข้อ 6 สมัครสมาชิกด้วย validator

โจทย์ไม่ได้ระบุฟิลด์ จึงใช้ 4 ฟิลด์ต่อไปนี้เป็นตัวอย่าง และเน้นขั้นตอนรับ/ตรวจสอบข้อมูล ยังไม่มีฐานข้อมูลหรือการบันทึกสมาชิก

| ฟิลด์ | กฎ |
| --- | --- |
| username | จำเป็น, ภาษาอังกฤษหรือตัวเลขเท่านั้น, 3–20 ตัวอักษร |
| email | จำเป็น, รูปแบบอีเมลถูกต้อง, ไม่เกิน 254 ตัวอักษร |
| password | จำเป็น, 8–72 ตัวอักษร |
| confirm_password | จำเป็น, ตรงกับ password |

ใน Postman เลือก `POST http://localhost:3000/api/v1/register` ตั้ง Basic Auth ตามด้านบน แล้วไปที่ Body → raw → JSON ใส่:

```json
{
  "username": "student01",
  "email": "student@example.com",
  "password": "password123",
  "confirm_password": "password123"
}
```

เมื่อข้อมูลถูกต้องจะได้ HTTP 200:

```json
{
  "message": "ข้อมูลสมัครสมาชิกถูกต้อง",
  "data": {
    "username": "student01",
    "email": "student@example.com"
  }
}
```

ลองเปลี่ยนอีเมลเป็น `abc` หรือตั้ง confirm_password ไม่ตรงกับ password จะได้ HTTP 400:

```json
{"message":"ใส่ข้อมูลผิดพลาด"}
```

ให้อ่าน models/register.go ก่อน ตัวอย่าง tag `validate:"required,email"` หมายถึงต้องกรอกและต้องเป็นอีเมล ส่วน `eqfield=Password` หมายถึงต้องตรงกับฟิลด์ Password

จากนั้นอ่าน controller/register.go ซึ่งมีเพียง 3 ขั้นตอน:

1. `ctx.BodyParser(&input)` อ่าน JSON มาใส่ struct
2. `validate.Struct(input)` ตรวจทุกฟิลด์ตาม tags ถ้าผิดให้ return ข้อความผิดพลาด
3. ส่งข้อมูลที่ผ่านการตรวจสอบกลับ โดยไม่ส่งรหัสผ่าน

ใช้กฎสำเร็จรูปของ validator ตามโจทย์ จึงไม่ต้องเขียน regexp เองสำหรับกฎเหล่านี้

เอกสารอ้างอิง: [go-playground/validator](https://github.com/go-playground/validator)

## คู่มือ PDF และ Postman Collection

- [คู่มือภาษาไทย: วิธีใช้งาน หลักการทำงาน และการทดสอบ](output/pdf/go-fiber-workshop-guide-th.pdf)
- [Postman Collection พร้อมทดสอบ 33 request](postman/Go-Fiber-Workshop.postman_collection.json)
- [วิธีนำเข้าและใช้งาน Collection](postman/README.md)

รัน `go run .` แล้ว Import Collection เข้า Postman ได้เลย มี Basic Auth ตั้งไว้ให้แล้ว
