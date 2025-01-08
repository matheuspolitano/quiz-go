curl -X GET localhost:8081/api/quiz/types \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg5N2FjZmRlLTFiYTgtNDE0ZC1iNWM0LTY5MzZkOTlkYTg3MiIsInVzZXJuYW1lIjoibWF0aGV1c3BvbGl0YW5vMiIsInJvbGUiOiJyZWd1bGFyIiwiaXNzdWVkX2F0IjoiMjAyNS0wMS0wOFQwODozMDoxOS44ODkxMTQxKzAxOjAwIiwiZXhwaXJlZF9hdCI6IjIwMjYtMDEtMDdUMDg6MzA6MTkuODg5MTE0MSswMTowMCJ9.EQuR4oHOSfoxFka5xiqQFwoHo8XT-8YPxUx14hdRQnE" 

# return all quiz types availables for the user
# status 200 [{"name":"CountryQuestions","questions_id":["1","2","3","4","5","6","7","8","9","10"]}]