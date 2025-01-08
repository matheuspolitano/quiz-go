# Quiz CLI Application

This **Quiz CLI** is a terminal-based application that tests your knowledge by connecting to a quiz API. You can log in, select a quiz type, answer questions, and retrieve your final score—all from the command line. It uses [Cobra](https://github.com/spf13/cobra) for structuring commands and a clean, modular Go architecture.

## Features

- **Login**: Prompt for a username and fetch an API access token.
- **Quiz Types**: Browse available quiz categories from the server.
- **Question Flow**: Fetch each question, submit your answer, and proceed.
- **Scoring**: Display final quiz results and overall accuracy rates.

## Project Structure

```
client
├── cmd          // Contains root and start commands
├── internal
│   ├── api      // API client calls (login, join quiz, next question, etc.)
│   ├── models   // Data structures (QuizType, Question, ScoreResponse, etc.)
│   └── quiz     // Quiz flow logic (prompts, question loop, score display)
├── main.go      // Entry point; runs Cobra commands
```

## Getting Started

1. **Clone** the repository:

   ```bash
   git clone https://github.com/matheuspolitano/quiz-go.git
   ```

2. **API_URL** Set the API URL:

   - Set as app.env or in variable environment

3. **Run** the binary:

   ```bash
   cd client
   go run main.go start
   ```

4. **Follow** the on-screen prompts to:
   - Enter your username
   - Choose a quiz type
   - Answer each question
   - View your final score

# Quiz Go App Backend

This repository showcases a simple quiz application written in Go, using an in‑memory (file‑based) database (`memdb`) to store users, questions, quiz types, and user quiz flows. The application is containerized using a multi-stage Docker build for optimized images. Below you will find key highlights of the app’s architecture, usage instructions, and the available API endpoints.

---

## Table of Contents

1. [Overview](#overview)
2. [Main Features](#main-features)
   - [Memdb (File‑based)](#memdb-file-based)
   - [API Endpoints](#api-endpoints)
3. [Docker and Data Persistence](#docker-and-data-persistence)
4. [Running Locally](#running-locally)
   - [Using Docker Compose](#using-docker-compose)
   - [Exposed Ports](#exposed-ports)
5. [Detailed API Endpoints](#detailed-api-endpoints)
   - [1. `POST /api/login`](#1-post-apilogin)
   - [2. `GET /api/quiz/types`](#2-get-apiquiztypes)
   - [3. `GET /api/quiz/question/:questionID`](#3-get-apiquizquestionquestionid)
   - [4. `POST /api/quiz/joinQuiz/:typeQuiz`](#4-post-apiquizjoinquiztypequiz)
   - [5. `GET /api/quiz/answer/:typeQuiz/next`](#5-get-apiquizanswertypequiznext)
   - [6. `POST /api/quiz/answer/:typeQuiz/:questionID`](#6-post-apiquizanswertypequizquestionid)
   - [7. `GET /api/quiz/answer/:typeQuiz/score`](#7-get-apiquizanswertypequizscore)
6. [Common Errors](#common-errors)
7. [Code Structure](#code-structure)
8. [License](#license)

---

## Overview

This Go application implements a quiz service where users can:

- **Create** or **reuse** an account (identified by a username),
- **Retrieve** available quiz types,
- **Join** a selected quiz flow,
- **Answer** questions,
- **See** their scores.

Data is stored in JSON files under a `data` folder, managed by a custom in‑memory (file‑based) repository.

---

## Main Features

### Memdb (File‑based)

The application uses a simple file-based repository (`memdb`) to simulate a small database. Each entity (User, Question, Quiz Type, Question Flow, History of answers) is stored as JSON in separate files:

- `users.data.json`
- `questions.data.json`
- `typesQuiz.data.json`
- `questionsFlows.data.json`
- `history.data.json`

To maintain data across container restarts, we **strongly recommend** binding or mounting a volume to the `data` folder so that newly added questions, users, and quiz flows persist.

### API Endpoints

The app exposes a RESTful API under `/api`. Authentication is done via a generated JWT token when a user logs in. Below is a high-level summary:

1. **`POST /api/login`**: Creates/reuses a user and returns a JWT.
2. **`GET /api/quiz/types`**: Lists available quiz types.
3. **`GET /api/quiz/question/:questionID`**: Retrieves a specific question.
4. **`POST /api/quiz/joinQuiz/:typeQuiz`**: Creates a new quiz flow for a user.
5. **`GET /api/quiz/answer/:typeQuiz/next`**: Fetches the next unanswered question.
6. **`POST /api/quiz/answer/:typeQuiz/:questionID`**: Submits the user’s answer for a question.
7. **`GET /api/quiz/answer/:typeQuiz/score`**: Retrieves the user’s current quiz flow score and a general accuracy rate.

All quiz endpoints (except `POST /api/login`) require a **Bearer** token in the `Authorization` header.

---

## Docker and Data Persistence

**Data Volume**: The application stores all data (users, questions, etc.) under `/app/data`. In the provided `docker-compose.yml`, the `data` directory is mounted so that the JSON files persist across restarts.

```yaml
services:
  quiz_local:
    build:
      context: './backend'
      dockerfile: Dockerfile
    environment:
      - API_PORT=80 # default value
    ports:
      - '80:80'
    volumes:
      - ./data:/app/data
```

Make sure the `data` directory is present in your host machine, as it will be mapped into the container. Otherwise, the app will create it on startup.

---

## Running Locally

### Using Docker Compose

1. **Clone** the repository.
2. Create a `data` folder in the project’s root (if not already there). This directory will hold JSON files.
3. Run `docker compose up --build -d`.

This will:

- Build the Go application using the multi-stage Dockerfile.
- Start the container, exposing port `80`.

### Exposed Ports

The application listens on port **80** within the container, which is mapped to your machine’s port **80** by default in the `docker-compose.yml`.

You can modify this in the `docker-compose.yml` file if desired, for example to `8080:80`.

---

## Detailed API Endpoints

Below are the main endpoints along with example cURL commands and typical responses.

Remember to include:

```
-H "Authorization: Bearer <your_jwt_token>"
```

on **all** endpoints requiring authentication.

---

### 1. `POST /api/login`

**Description**: Registers a new user or logs in an existing user, returning a JWT token.

- **Example Request**:

  ```bash
  curl -X POST localhost:8081/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"matheuspolitano3"}'
  ```

- **Example Response** (`201 Created`):
  ```json
  {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

---

### 2. `GET /api/quiz/types`

**Description**: Lists all quiz types available for the user.

- **Example Request**:

  ```bash
  curl -X GET localhost:8081/api/quiz/types \
  -H "Authorization: Bearer <jwt_token>"
  ```

- **Example Response** (`200 OK`):
  ```json
  [
    {
      "name": "CountryQuestions",
      "questions_id": ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10"]
    }
  ]
  ```

---

### 3. `GET /api/quiz/question/:questionID`

**Description**: Fetches a specific question by ID.

- **Example Request**:

  ```bash
  curl -X GET localhost:8081/api/quiz/question/1 \
  -H "Authorization: Bearer <jwt_token>"
  ```

- **Example Response** (`200 OK`):
  ```json
  {
    "id": "1",
    "prompt": "What is the capital of Italy?",
    "options": ["A: Rome", "B: Milan", "C: Venice", "D: Florence"],
    "answer": "A"
  }
  ```

---

### 4. `POST /api/quiz/joinQuiz/:typeQuiz`

**Description**: Initiates (joins) a quiz flow of a given quiz type for the user.

- **Example Request**:

  ```bash
  curl -X POST localhost:8081/api/quiz/joinQuiz/CountryQuestions \
  -H "Authorization: Bearer <jwt_token>"
  ```

- **Example Success Response** (`202 Accepted`):

  ```json
  {
    "user_id": "matheuspolitano3",
    "type_quiz": "CountryQuestions",
    "history": [],
    "created_at": "2025-01-08T11:45:00.6666329+01:00",
    "closed_at": "0001-01-01T00:00:00Z",
    "accuracy_rate": 1
  }
  ```

- **Example Error** (`400 Bad Request`):
  ```json
  {
    "status": "error",
    "error": "AddQuestionFlow: question flow already exists for user matheuspolitano3 and type CountryQuestions"
  }
  ```

---

### 5. `GET /api/quiz/answer/:typeQuiz/next`

**Description**: Retrieves the **next unanswered** question in an ongoing quiz flow.

- **Example Request**:

  ```bash
  curl -X GET localhost:8081/api/quiz/answer/CountryQuestions/next \
  -H "Authorization: Bearer <jwt_token>"
  ```

- **Example Success Response** (`200 OK`):

  ```json
  {
    "id": "1",
    "prompt": "What is the capital of Italy?",
    "options": ["A: Rome", "B: Milan", "C: Venice", "D: Florence"],
    "answer": "A"
  }
  ```

- **Example Error** (no more questions):
  ```json
  {
    "status": "error",
    "error": "question flow is already closed"
  }
  ```

---

### 6. `POST /api/quiz/answer/:typeQuiz/:questionID`

**Description**: Submits an answer for a question in an ongoing quiz flow.

- **Example Request**:

  ```bash
  curl -X POST localhost:8081/api/quiz/answer/CountryQuestions/1 \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{"answer":"C"}'
  ```

- **Example Success Response** (`202 Accepted`):

  ```json
  {
    "id": "ac1437ca-4e5c-4cb7-af77-82b7939dc8b1",
    "used_id": "matheuspolitano3",
    "question_id": "1",
    "answer": "C",
    "is_right": false,
    "created_at": "2025-01-08T11:55:59.371504+01:00"
  }
  ```

- **Example Error** (already answered):
  ```json
  {
    "status": "error",
    "error": "AddAnswer: question already answer. Use the next to get the question without answer"
  }
  ```

---

### 7. `GET /api/quiz/answer/:typeQuiz/score`

**Description**: Returns the current quiz flow’s score (accuracy rate) and the general accuracy rate across all flows.

- **Example Request**:

  ```bash
  curl -X GET localhost:8081/api/quiz/answer/CountryQuestions/score \
  -H "Authorization: Bearer <jwt_token>"
  ```

- **Example Success Response** (`202 Accepted`):
  ```json
  {
    "user_quiz": {
      "user_id": "matheuspolitano2",
      "type_quiz": "CountryQuestions",
      "history": ["3c78f399-062d-4602-9399-e7a8fc16ceba", "... more IDs ..."],
      "created_at": "2025-01-08T10:21:23.1649101+01:00",
      "closed_at": "2025-01-08T10:25:10.6848278+01:00",
      "accuracy_rate": 0.3
    },
    "general_accuracy_rates": 0.15
  }
  ```

---

## Common Errors

- **400 Bad Request**: Happens if you provide invalid JSON, try to rejoin a quiz flow that already exists, or answer a question more than once.
- **401 Unauthorized**: Returned if no valid `Authorization` header with JWT is provided on protected endpoints.
- **404 Not Found**: Occurs when a quiz type, question, or user does not exist (or was not found).

---

## Code Structure

Below is a simplified structure of the code:

```
.
├── cmd
│   └── quiz
│       └── main.go         # entry point
├── internal
│   ├── api
│   │   ├── server.go       # server setup, routes, startup, shutdown
│   │   ├── middleware.go   # token-based middleware
│   │   └── ...
│   ├── memdb
│   │   ├── db_manager.go   # repository aggregator (DBManager)
│   │   ├── repository.go   # generic repository logic (loads/saves JSON)
│   │   └── ...
│   ├── models
│   │   └── ...             # data structures (User, Question, etc.)
│   ├── token
│   │   └── ...             # JWT creation, parsing
│   └── ...
├── data                    # JSON files stored here (mounted as a volume in Docker)
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

- **`internal/api`**: Contains the Gin server logic, routes, and middleware for JWT auth.
- **`internal/memdb`**: Implements a simple file-based repository with concurrency safety.
- **`internal/models`**: Defines the Go structs used by the app (`User`, `Question`, `TypeQuiz`, `QuestionFlow`, etc.).
- **`internal/token`**: Provides JWT token generation and validation.

---
