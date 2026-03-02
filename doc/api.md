# Deutsch API 接口文档

**Base URL**: `http://localhost:8888`

**通用响应格式**（除特别说明外）：

```json
{
  "code": 0,
  "msg": "success",
  "data": { ... }
}
```

错误时 `code` 非 0，`msg` 为错误描述。

**鉴权**：需登录的接口在 Header 中携带 `Authorization: bearer <jwt_token>`。

---

## 1. 配置与元数据

### 1.1 获取联邦州列表

`GET /api/v1/states`

**响应**

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 17,
    "items": [
      {
        "id": "00000000-0000-0000-0000-000000000001",
        "slug": "general",
        "name": "Allgemein",
        "nameCn": "通用"
      }
    ]
  }
}
```

### 1.2 获取应用配置

`GET /api/v1/config`

**响应**

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "totalQuestions": 460,
    "examQuestions": 33,
    "examMinutes": 60,
    "passScore": 17,
    "languageModes": [
      { "value": "de", "label": "Deutsch" },
      { "value": "cn", "label": "中文" }
    ]
  }
}
```

---

## 2. 认证

### 2.1 注册

`POST /api/v1/auth/register`

**请求体**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "invite_code": "XXXX",
  "username": "optional",
  "nickname": "可选",
  "phone": "可选"
}
```

**响应**：返回用户信息与 JWT，同登录。

### 2.2 登录

`POST /api/v1/auth/login`

**请求体**（邮箱或手机二选一）

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
或
```json
{
  "phone": "13800138000",
  "password": "password123"
}
```

**响应**

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "user": {
      "id": "xxx",
      "username": "xxx",
      "email": "xxx",
      "role": "user"
    },
    "jwt_token": "eyJhbG...",
    "expires": 1772440487,
    "max_refresh": 1772433287
  }
}
```

### 2.3 刷新 Token

`POST /api/v1/auth/jwt/refresh`（需鉴权）

**响应**：同登录，返回新的 token。

### 2.4 登出

`POST /api/v1/auth/logout`（需鉴权）

---

## 3. 题目

### 3.1 获取通用题列表

`GET /api/v1/questions/general`

**响应**

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 300,
    "items": [
      {
        "id": "xxx",
        "questionDe": "题目德文",
        "questionCn": "题目中文",
        "optionsDe": ["选项A", "选项B", "..."],
        "optionsCn": ["选项A", "选项B", "..."],
        "correctAnswer": 0,
        "explanation": "解析文本",
        "hasImage": false,
        "state": "general"
      }
    ]
  }
}
```

### 3.2 获取某州题目

`GET /api/v1/questions/state/:state_id`

**路径参数**：`state_id` 为州 ID（如 `5465e833-302f-4de6-9ca4-c082919dfa68`）或 slug（如 `baden-wuerttemberg`）。

### 3.3 模拟考试抽题

`GET /api/v1/questions/exam?state_id=xxx`

**查询参数**（可选）：`state_id`，不传则随机选一州。

**响应**：30 道通用题 + 3 道指定州题，结构与 3.1 一致。

### 3.4 获取全部题目（按州分组）

`GET /api/v1/questions/all`

**响应**

```json
{
  "code": 0,
  "data": {
    "general": [...],
    "5465e833-xxx": [...]
  }
}
```

### 3.5 获取单题详情

`GET /api/v1/questions/:question_id`

---

## 4. 学习进度（需鉴权）

### 4.1 获取用户偏好

`GET /api/v1/progress/preferences`

**响应**

```json
{
  "code": 0,
  "data": {
    "preferredExamStateId": "5465e833-xxx"
  }
}
```

### 4.2 更新偏好

`PATCH /api/v1/progress/preferences`

```json
{
  "preferredExamStateId": "5465e833-xxx"
}
```

### 4.3 获取学习进度

`GET /api/v1/progress/learning`

**响应**

```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "stateId": "xxx",
        "total": 10,
        "practicedCount": 5,
        "correctCount": 4
      }
    ]
  }
}
```

### 4.4 记录练习

`POST /api/v1/progress/learning`

```json
{
  "questionId": "xxx",
  "correct": true
}
```

### 4.5 考试记录列表

`GET /api/v1/progress/exams?pageNo=1&pageSize=10`

**查询参数**：`pageNo`（默认 1）、`pageSize`（默认 10，最大 100）。

**响应**

```json
{
  "code": 0,
  "data": {
    "total": 100,
    "items": [
      {
        "id": "xxx",
        "stateId": "xxx",
        "total": 33,
        "score": 25,
        "passed": true,
        "createdAt": 1730000000000
      }
    ]
  }
}
```

### 4.6 提交考试

`POST /api/v1/progress/exams`

```json
{
  "stateId": "5465e833-xxx",
  "answers": {
    "question-id-1": 0,
    "question-id-2": 2
  }
}
```

`answers`：`questionId -> optionIndex`（0-based）。

**响应**

```json
{
  "code": 0,
  "data": {
    "id": "记录ID",
    "total": 33,
    "score": 25,
    "passed": true,
    "details": [
      {
        "questionId": "xxx",
        "chosenAnswer": 1,
        "correct": false,
        "correctOptionIndex": 3,
        "correctOptionDe": "hier Meinungsfreiheit gilt.",
        "correctOptionCn": "这里有言论自由。",
        "explanation": "《德国基本法》第5条规定了言论自由..."
      }
    ],
    "createdAt": 1730000000000
  }
}
```

### 4.7 考试记录详情

`GET /api/v1/progress/exams/:id`

**响应**：含 `details`，结构与 4.6 中 `details` 相同，每题含正确选项与解析。

### 4.8 错题本列表

`GET /api/v1/progress/wrong-questions?pageNo=1&pageSize=10`

### 4.9 添加错题

`POST /api/v1/progress/wrong-questions`

```json
{
  "questionId": "xxx"
}
```

### 4.10 移除错题

`DELETE /api/v1/progress/wrong-questions/:question_id`

---

## 5. 用户（需鉴权）

### 5.1 获取当前用户

`GET /api/v1/user/:id`

### 5.2 更新个人资料

`PATCH /api/v1/user/profile`

```json
{
  "nickname": "新昵称",
  "description": "个人简介"
}
```

### 5.3 用户列表（Admin）

`GET /api/v1/user/list?pageNo=1&pageSize=10`

### 5.4 删除用户（Admin）

`DELETE /api/v1/user/:id`

### 5.5 启用/禁用用户（Admin）

`PATCH /api/v1/user/:id/enable`

```json
{
  "id": "xxx",
  "is_enabled": true
}
```

---

## 6. 邀请码

### 6.1 验证邀请码（无需鉴权）

`GET /api/v1/invitecode/validate/:id`

### 6.2 生成邀请码（Admin）

`POST /api/v1/invitecode/generate`

```json
{
  "count": 5
}
```

### 6.3 邀请码列表（Admin）

`GET /api/v1/invitecode/list?pageNo=1&pageSize=10&availableOnly=true`

### 6.4 启用/禁用邀请码（Admin）

`PATCH /api/v1/invitecode/:id/enable`

### 6.5 删除邀请码（Admin）

`DELETE /api/v1/invitecode/:id`

---

## 分页通用参数

使用 `Filter` 的接口支持：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| pageNo | int | 1 | 页码 |
| pageSize | int | 10 | 每页条数（最大 100） |
