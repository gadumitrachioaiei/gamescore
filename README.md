curl -X POST --data '{"user": 1, "total": 12}' "http://localhost:8080/scores/"

curl "http://localhost:8080/scores/top/?top=10"

curl "http://localhost:8080/scores/range/?position=10&count=2"