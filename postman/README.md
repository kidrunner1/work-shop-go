# Postman Collection สำหรับ Go Fiber Workshop

## เริ่มใช้งาน

1. เปิด terminal ใน `C:\src\workshop-go` แล้วรัน `go run .` เปิด server ค้างไว้
2. เปิด Postman เลือก **Import** แล้วเลือกไฟล์ `Go-Fiber-Workshop.postman_collection.json` ในโฟลเดอร์นี้
3. เปิด Collection **Go Fiber Workshop - ข้อ 5.0 ถึง 6**
4. ลองโฟลเดอร์ **5.1 Factorial - รับ Path** เลือก **Factorial 5 (200)** แล้วกด **Send**
5. ดู response: `{"number":5,"factorial":120}` และดู **Test Results** ว่าผ่าน

Collection มี Basic Auth `gofiber` / `21022566` พร้อมใช้ ไม่ต้องสร้าง Environment เพิ่ม

## มีอะไรให้ลอง

| โฟลเดอร์ | จำนวน request | สิ่งที่ทดสอบ |
| --- | ---: | --- |
| 5.1 Factorial | 7 | เลขปกติ, 0, 20, ค่าติดลบ, เกินขนาด, ข้อความ, ทศนิยม |
| 5.2 ASCII | 5 | Hello, ตัวอักษรผสมตัวเลข, ไม่ส่ง query, ค่าว่าง, ภาษาไทย |
| 6 Register | 12 | ข้อมูลถูก, ฟิลด์หาย, JSON พัง, ชนิดข้อมูลผิด และกฎ validator |
| 5.0 Basic Auth | 9 | ไม่ส่ง Auth, รหัสผิด, ชื่อผู้ใช้ผิด สำหรับทั้ง 3 API |

รวม 33 request พร้อมสคริปต์ตรวจ status และผลตอบกลับใน **Scripts > Post-response** (บางเวอร์ชันชื่อ **Tests**)

## โฟกัสข้อ 6

เปิด **Register - ข้อมูลถูกต้อง (200)** แล้วดู **Body > raw > JSON** จะเห็นตัวอย่าง username, email, password และ confirm_password แก้แล้วกด Send ได้เลย

ถ้าต้องการลองข้อมูลผิด ให้เลือก request ที่เตรียมไว้ เช่น **Register - email ผิดรูปแบบ (400)** ซึ่งคาดหวังข้อความ `ใส่ข้อมูลผิดพลาด`

แต่ละ request มีผลที่คาดหวังกำหนดไว้ หากเปลี่ยนข้อมูลของ request สำเร็จเป็นข้อมูลผิด สคริปต์เดิมที่คาดหวัง 200 จะขึ้น FAIL ตามปกติ ใช้ request กรณีผิดที่เตรียมไว้เพื่อเทียบผลได้ง่าย

API นี้ตรวจข้อมูลเท่านั้น ยังไม่มีฐานข้อมูลหรือการบันทึกสมาชิก

## ทดสอบทั้งหมด

เลือก Collection แล้วใช้ **Run collection / Run** เลือกทุก request ตั้ง Iterations เป็น 1 แล้วเริ่มรัน

- กรณีข้อมูลถูกต้องคาดหวัง 200
- กรณีข้อมูลผิดคาดหวัง 400
- กรณี Auth ผิดคาดหวัง 401

ดังนั้น response 400 หรือ 401 ในกรณีที่ตั้งใจทดสอบ ต้องแสดงผลทดสอบเป็น **PASS**

## ตัวแปรใน Collection

| ตัวแปร | ค่าเริ่มต้น | วิธีใช้ |
| --- | --- | --- |
| base_url | http://localhost:3000 | URL server ไม่ใส่ / ท้ายค่า |
| nickname | joe | ชื่อ route ข้อ 5.2 ไม่ใส่ / |
| basic_user | gofiber | ชื่อผู้ใช้สำหรับ Basic Auth |
| basic_pass | 21022566 | รหัสผ่านสำหรับ Basic Auth |

แก้ค่าที่ Collection > Variables ถ้าเปลี่ยน nickname ต้องแก้ route ใน `routes/routes.go` ให้ตรงกัน แล้วเริ่ม server ใหม่ด้วย

ถ้าเลือก Environment ที่มีชื่อตัวแปรซ้ำ ค่าของ Environment อาจแทนค่าของ Collection ให้เลือก No environment สำหรับตัวอย่างนี้

ถ้าเชื่อมต่อไม่ได้ ตรวจว่าเปิด `go run .` อยู่และใช้ port 3000 หากใช้ Postman บนเว็บ ให้ใช้ Desktop Agent เพื่อเรียก server ในเครื่อง

อ้างอิงการเขียนสคริปต์: https://learning.postman.com/docs/tests-and-scripts/write-scripts/test-scripts/
