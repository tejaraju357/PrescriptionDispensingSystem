Absolutely! Here’s a **complete, ready-to-paste `README.md` file** for your GitHub repository. It’s fully humanized, includes setup instructions, API documentation with examples, and the database schema—all in one file:

---

````markdown
# Prescription Dispensing System

A simple, fast, and reliable system for managing medicines, prescriptions, and dispensing stock.  
Built with **Go**, **Fiber**, and **PostgreSQL**, with optional **Redis** for concurrency-safe stock management.

---

## Table of Contents
1. [Overview](#overview)  
2. [Setup Instructions](#setup-instructions)  
3. [API Documentation](#api-documentation)  
4. [Database Schema](#database-schema)  

---

## Overview

This project allows you to:

- Manage users (admin, pharmacists).  
- Add and view medicines.  
- Dispense stock safely even under concurrent requests.  
- Create and track prescriptions.  

All APIs return JSON responses and support **JWT authentication** for protected routes.

---

## Setup Instructions

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/prescription-dispensing-system.git
cd prescription-dispensing-system
````

### 2. Install Go dependencies

```bash
go mod tidy
```

### 3. Setup PostgreSQL

* Create the database:

```sql
CREATE DATABASE prescription_db;
```

* Update `.env` or `config` with your DB credentials:

```
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=prescription_db
DB_HOST=localhost
DB_PORT=5432
```

### 4. Run Database Migrations

Create tables manually or via migration scripts:

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255),
    role VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Medicines table
CREATE TABLE medicines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    dosage_form VARCHAR(50),
    stock_quantity INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Prescriptions table
CREATE TABLE prescriptions (
    id SERIAL PRIMARY KEY,
    patient_name VARCHAR(100),
    medicine_name VARCHAR(100),
    dosage VARCHAR(50),
    quantity INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### 5. Run the Server

```bash
go run main.go
```

Server runs at: `http://localhost:9000`

---

## API Documentation

Base URL:

```
http://localhost:9000
```

All requests and responses are JSON.

---

### 1. Register a User

**POST** `/register`

**Request:**

```json
{
  "name": "Admin",
  "email": "admin@test.com",
  "password": "pass123",
  "role": "admin"
}
```

**Response:**

```json
{
  "id": 1,
  "name": "Admin",
  "email": "admin@test.com",
  "role": "admin"
}
```

---

### 2. Login

**POST** `/login`

**Request:**

```json
{
  "email": "admin@test.com",
  "password": "pass123"
}
```

**Response:**

```json
{
  "message": "Login successful",
  "token": "<jwt_token>"
}
```

---

### 3. Add Medicine

**POST** `/addMedicine`

**Request:**

```json
{
  "name": "Paracetamol",
  "dosage_form": "Tablet",
  "stock_quantity": 100
}
```

**Response:**

```json
{
  "id": 1,
  "name": "Paracetamol",
  "dosage_form": "Tablet",
  "stock_quantity": 100
}
```

---

### 4. Get All Medicines

**GET** `/medicines`

**Response:**

```json
[
  {
    "id": 1,
    "name": "Paracetamol",
    "dosage_form": "Tablet",
    "stock_quantity": 100
  }
]
```

---

### 5. Dispense Stock

**POST** `/dispenseStock`

**Request:**

```json
{
  "name": "Paracetamol",
  "dosage_form": "Tablet",
  "stock_quantity": 5
}
```

**Response:**

```json
{
  "name": "Paracetamol",
  "dosage_form": "Tablet",
  "stock_quantity": 5
}
```

> Note: Stock updates are **atomic**. Redis locks prevent negative stock during concurrent requests.

---

### 6. Create Prescription

**POST** `/presc`

**Request:**

```json
{
  "patient_name": "John Doe",
  "medicine_name": "Paracetamol",
  "dosage": "500mg",
  "quantity": 2
}
```

**Response:**

```json
{
  "id": 1,
  "patient_name": "John Doe",
  "medicine_name": "Paracetamol",
  "dosage": "500mg",
  "quantity": 2
}
```

---

## Database Schema

### ERD Diagram

```
+-----------+       +-----------+       +--------------+
|   users   |       | medicines |       | prescriptions|
+-----------+       +-----------+       +--------------+
| id        |<----->| id        |       | id           |
| name      |       | name      |       | patient_name |
| email     |       | dosage_form|      | medicine_name|
| password  |       | stock_qty |       | dosage       |
| role      |       | created_at|       | quantity     |
| created_at|       | updated_at|       | created_at   |
| updated_at|                           | updated_at   |
+-----------+                           +--------------+
```

---

### Notes

* JWT tokens are required for protected endpoints.
* Redis is used for concurrency-safe stock dispensing.
* All endpoints return clear JSON messages for easy integration.

```

---

✅ This is **all-in-one**, ready for GitHub. You just need to:

1. Copy it into `README.md`.  
2. Update repository URLs and `.env` info if needed.  
3. Optionally add ERD images/screenshots to make it visual.  

---

If you want, I can also **add a Quick Start with Docker + Redis section** so anyone can run it with one command locally.  

Do you want me to do that?
```
