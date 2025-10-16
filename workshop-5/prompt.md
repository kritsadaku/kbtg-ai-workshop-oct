
## Step

MCP installation
https://code.visualstudio.com/docs/copilot/customization/mcp-servers

### 1. GitMCP
1. เข้า https://gitmcp.io/ แล้วใส่ repo นี้เข้าไป > กด "To MCP"
2. เลือก VS Code > สร้าง file `.vscode/mcp.json` แล้วเพิ่ม config ตามด้านล่างนี้เข้าไป พร้อมกด start mcp
```json
{
  "servers": {
    "kbtg-ai-workshop-oct Docs": {
      "type": "sse",
      "url": "https://gitmcp.io/mikelopster/kbtg-ai-workshop-oct"
    }
  }
}
```
3. กลับมาที่ Github Copilot ทดสอบ Prompt 

```
จง follow style guideline ตาม repo kbtg-ai-workshop-oct Docs
```

### 2. Playwright MCP

1. เข้า https://github.com/microsoft/playwright-mcp แล้วเพิ่ม config ลงใน `.vscode/mcp.json` ตาม config ด้านล่าง (follow ตาม github ได้) แล้วกด Start MCP

```json
{
  "servers": {
    "...": {},
    "playwright": {
      "command": "npx",
      "args": [
        "@playwright/mcp@latest",
        "--isolated",
      ]
    }
  }
}
```

2. ทดสอบ Prompt

```
ใช้ Playwright MCP ตรวจสอบ https://www.kbtg.tech/th/home ว่าเว็บไซต์ใส่ Contact อยู่ที่ Nonthaburi ใช่ไหม
```