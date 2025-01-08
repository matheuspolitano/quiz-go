curl -X GET localhost:8081/api/quiz/answer/CountryQuestions/next \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImZhNjFlMDNjLTVlYjUtNGIyOC04NmFjLWE1OGNhNzNmYjhhNSIsInVzZXJuYW1lIjoibWF0aGV1c3BvbGl0YW5vMyIsInJvbGUiOiJyZWd1bGFyIiwiaXNzdWVkX2F0IjoiMjAyNS0wMS0wOFQxMTo0NDowNS4xODAzOTEzKzAxOjAwIiwiZXhwaXJlZF9hdCI6IjIwMjYtMDEtMDhUMTE6NDQ6MDUuMTgwMzkxMyswMTowMCJ9.agWAzjdJYoc4OoNyD_DMLjonxOd1vQ1CNd1VvJRS2nw" 


## Will get the next question to the user 
## 200 {"id":"1","prompt":"What is the capital of Italy?","options":["A: Rome","B: Milan","C: Venice","D: Florence"],"answer":"A"}

## Case does not exist more question will return this message
##  400 {"status":"error","error":"question flow is already closed"}