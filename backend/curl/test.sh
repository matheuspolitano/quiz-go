curl -X POST localhost:8081/api/login \
-H "Content-Type: application/json" \
-d '{"username":"matheuspolitano3"}'

## return the access token
## status 201 {"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImY3MmE0NDJhLWU4YmQtNDMzZS05YmRiLTQ3MmM3NDFiZDEwMCIsInVzZXJuYW1lIjoibWF0aGV1c3BvbGl0YW5vMiIsInJvbGUiOiJyZWd1bGFyIiwiaXNzdWVkX2F0IjoiMjAyNS0wMS0wOFQxMToyNzo1MC45NTM4NzYxKzAxOjAwIiwiZXhwaXJlZF9hdCI6IjIwMjYtMDEtMDhUMTE6Mjc6NTAuOTUzODc2MSswMTowMCJ9.Ae9qWJhkYyXLFv5rX9xQwJ0dKB5M__paw6PSJiLQnMs"}