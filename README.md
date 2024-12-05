# user-auth
Is a project sampling a user authentication process. It includes support of mysql and postgres databases, user creation, login and logout and session management.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)

## Prerequisites
- [go installation](https://go.dev/doc/install)
- [mysql](https://www.mysql.com/downloads/) or [postgres](https://www.postgresql.org/download/)

## Run
Environment variables can be set in the .env file. Please find sample values below:

```env
SERVER="127.0.0.1:8089"
DATA_SOURCE="mysql"
DB_HOST="127.0.0.1"
DB_ADDR="127.0.0.1:3306"
DB_USER="user"
DB_PASS="12345678"
DB_NAME="stub"
LOG_LEVEL="DEBUG"
TLS_CERT_PATH="/tls/cert.pem"
TLS_KEY_PATH="/tls/key.pem"
ROOT_DIR="/path/to/project/"
```

## API

#### **POST /user/signup**
Sign up user.

**Request:**
```http
POST /user/login HTTP/1.1
Content-Type: application/json

{
    "username":"jd123",
    "email": "john.doe@example.com",
    "password": "password123"
}
```

**Response:**
```json
{
  "message": "User created successfully",
}
```

| Status Code | Description               |
|-------------|---------------------------|
| 200         | Login successful.         |
| 401         | Invalid credentials.      |

---

#### **POST /user/login**
Log in a user. Authentication token in saved in session cookie.

**Request:**
```http
POST /user/login HTTP/1.1
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "User logged in",
}
```

| Status Code | Description               |
|-------------|---------------------------|
| 200         | Login successful.         |
| 401         | Invalid credentials.      |

---
#### **POST /user/login**
Log in a user. Authentication token in saved in session cookie.

**Request:**
```http
POST /user/login HTTP/1.1
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "User logged in",
}
```

| Status Code | Description               |
|-------------|---------------------------|
| 200         | Login successful.         |
| 401         | Invalid credentials.      |

---

#### Protected Endpoint: **GET /user/view**
Get user account data. User needs to login in first to request his account data.

**Request:**
```http
GET /user/view?email=john.doe@example.com HTTP/1.1
```

**Response:**
```json
{
  "username": "jd123",
  "email": "john.doe@example.com",
  "created": "2024-12-05T16:06:51Z"
}
```

| Status Code | Description               |
|-------------|---------------------------|
| 200         | Login successful.         |
| 401         | Invalid credentials.      |

---