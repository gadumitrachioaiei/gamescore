Add a score:
curl -X POST --data '{"user": 1, "total": 12}' "http://localhost:8080/scores/"

Update a user's score:
curl -X PUT --data '{"user": 1, "score": -1}' "http://localhost:8080/scores/"

Top 10:
curl "http://localhost:8080/scores/top/?top=10"

Top 5 around position 10:
curl "http://localhost:8080/scores/range/?position=10&count=2"