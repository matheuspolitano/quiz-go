curl -X POST localhost:8081/api/quiz/joinQuiz/CountryQuestions \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImZhNjFlMDNjLTVlYjUtNGIyOC04NmFjLWE1OGNhNzNmYjhhNSIsInVzZXJuYW1lIjoibWF0aGV1c3BvbGl0YW5vMyIsInJvbGUiOiJyZWd1bGFyIiwiaXNzdWVkX2F0IjoiMjAyNS0wMS0wOFQxMTo0NDowNS4xODAzOTEzKzAxOjAwIiwiZXhwaXJlZF9hdCI6IjIwMjYtMDEtMDhUMTE6NDQ6MDUuMTgwMzkxMyswMTowMCJ9.agWAzjdJYoc4OoNyD_DMLjonxOd1vQ1CNd1VvJRS2nw" 


## return a questionFlow who is the responsable to manage the type of quiz to the user, init with accuracy 1 and empty
## 202 {"user_id":"matheuspolitano3","type_quiz":"CountryQuestions","history":[],"created_at":"2025-01-08T11:45:00.6666329+01:00","closed_at":"0001-01-01T00:00:00Z","accuracy_rate":1}

## case the user already joined have this quiz will return this error
## 400 {"status":"error","error":"AddQuestionFlow: question flow already exists for user matheuspolitano3 and type CountryQuestions"}