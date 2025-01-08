curl -X POST localhost:8081/api/quiz/answer/CountryQuestions/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImZhNjFlMDNjLTVlYjUtNGIyOC04NmFjLWE1OGNhNzNmYjhhNSIsInVzZXJuYW1lIjoibWF0aGV1c3BvbGl0YW5vMyIsInJvbGUiOiJyZWd1bGFyIiwiaXNzdWVkX2F0IjoiMjAyNS0wMS0wOFQxMTo0NDowNS4xODAzOTEzKzAxOjAwIiwiZXhwaXJlZF9hdCI6IjIwMjYtMDEtMDhUMTE6NDQ6MDUuMTgwMzkxMyswMTowMCJ9.agWAzjdJYoc4OoNyD_DMLjonxOd1vQ1CNd1VvJRS2nw"  \
-H "Content-Type: application/json" \
-d '{"answer":"C"}'

## This enpoint will send the ansewer 
## 202 {"id":"ac1437ca-4e5c-4cb7-af77-82b7939dc8b1","used_id":"matheuspolitano3","question_id":"1","answer":"C","is_right":false,"created_at":"2025-01-08T11:55:59.371504+01:00"}

## If you trie answer two time will ave a error
## 400 {"status":"error","error":"AddAnswer: question already answer. Use the next to get the question without answer"}