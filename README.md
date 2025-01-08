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
