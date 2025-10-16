

## 1. Setup project

- setup React project + tailwind ผ่าน AI (โดยใช้ url กำกับ)

```
Setup project React ด้วย Vite โดย follow ตาม Document url นี้ https://vite.dev/guide/ (โปรดอ่าน url ด้วย) โดยให้สร้างภายใน folder นี้ ห้ามสร้าง folder ใหม่
```

- setup tailwind

```
Setup tailwind ด้วยวิธีของ Vite + Tailwind (tab: Using Vite) โดย follow ตาม Document url นี้ https://tailwindcss.com/docs/installation/using-vite (โปรดอ่าน url ด้วย)
```

## 2. convert component จาก vue component มาเป็น React component

```
จงแปลง Vue Component ด้านล่างนี้ เป็น React Component ชื่อ `components/Payment.jsx` และ import มาใช้ในหน้า `App.jsx`

...<code vue component>...
```

## 3. ทำหน้าใหม่โดยใช้ design token ที่กำหนดไว้ (เพื่อให้ได้ style เหมือนกัน)

- copy `workshop-3/specs/index.css` ไปไว้ตำแหน่ง `src/index.css`

```
เปลี่ยน style หน้า `App.jsx` ให้ follow ตาม design token ไฟล์ `src/index.css`
```

- ปรับ design ตาม UI (ใช้ file `workshop-3/specs/design/example.png`)

```
ปรับ design หน้า App.jsx follow ตาม UI file ที่ upload design นี้
```

## 4. เพิ่ม Storybook เพื่อกำหนด Specs ของ Frontend Component

- setup storybook

```
Setup storybook ด้วยวิธีของ Vite โดย follow ตาม Document url นี้ https://storybook.js.org/docs/get-started/frameworks/react-vite (โปรดอ่าน url ด้วย)

และเพิ่ม `Payment.stories.js` สำหรับ specs ของ Payment Component 
```

