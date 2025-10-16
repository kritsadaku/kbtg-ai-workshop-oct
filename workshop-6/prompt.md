Guideline
- Awesome copilot: https://github.com/github/awesome-copilot
- Spec kit: https://github.com/github/spec-kit

## 1. เพิ่ม copilot-instruction.md

ใช้ Prompt ประมาณด้านล่างนี้ได้

```
please write guideline code for system prompt for github copilot follow from current source code to #file:copilot-instructions.md 
```

เสร็จแล้วตอน Prompt ใหม่ให้สำรวจดูว่า instruction โดนอ่านทุกรอบแล้วถูกต้องหรือไม่

## 2. Chatmode

- สามารถใช้ Prompt ด้านล่างนี้เพื่อ generate shell script ก่อนได้

```shell
จงสร้าง shell script 2 อันสำหรับ
- Seed data เข้า database sqlite ตาม csv 
- นับจำนวน record users ใน database ว่ามีทั้งหมดเท่าไหร่
```

- เสร็จแล้วสร้าง Chatmode มา เพื่อให้ไปอ่าน script สำหรับ seed data และ อ่านจำนวนได้ (สามารถดู guideline ได้จาก `workshop-6/chatmodes/db-helper.chatmode.md`)
  - ในส่วน tools ขึ้นอยู่กับว่า อยากให้ chatmode ทำอะไรได้บ้าง ซึ่งถ้าแค่อยากให้จัดการตาม script ก็แค่กำหนดให้ run shell command ได้พอ ดังตัวอย่างนี้

- เมื่อเรียบร้อยให้ใช้ Prompt เช่นด้านล่างนี้ได้ ก็จะเจอว่า chatmode จะทำการ force ไป run script ตามที่เตรียมมาได้

```
check จำนวน user ให้หน่อยว่าตอนนี้มีเท่าไหร่
```