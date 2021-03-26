Add a score:

curl -X POST --data '{"user": 1, "total": 12}' "http://localhost:8080/scores/"
curl -X POST --data '{"user": 2, "total": 13}' "http://localhost:8080/scores/"
curl -X POST --data '{"user": 3, "total": 11}' "http://localhost:8080/scores/"
curl -X POST --data '{"user": 4, "total": 10}' "http://localhost:8080/scores/"

Update a user's score:
curl -X PUT --data '{"user": 1, "score": -1}' "http://localhost:8080/scores/"
curl -X PUT --data '{"user": 2, "score": 1}' "http://localhost:8080/scores/"
curl -X PUT --data '{"user": 4, "score": 1}' "http://localhost:8080/scores/"
curl -X PUT --data '{"user": 3, "score": -1}' "http://localhost:8080/scores/"

Top 10:
curl "http://localhost:8080/scores/top/?top=10"

Top 5 around position 10:
curl "http://localhost:8080/scores/range/?position=10&count=2"


http://buttersquid.ink/?spore=bNobwRALmBcYDYEsB2BTMAaMAPAjDAzAGyZYBMBOmAnnt
AKzXnQ4AMmKAhjG2CgEbdMcAPZYYpIcKqDIcGODhRYAHQCuLDXiEA3bgF896cEvjI0JWkRJMALD
xowGYKk1bsu0Hvxkix0CfBSMhBy0AomahosWvC6ngZGkDCmqBjYtHbWBPjUtG7OrjbuMt6ekn4B
ItJlsvLwEeqaaXBxLAnGyYipFtlZ0JnOeUUFMDgAnMU1pTy+4pLVPCF1ismRTTr6hh2wXebpvdh
MpBODo8MuoyecJQI1s-7zwaHhq43Rza3tSTtmabjiJzI4nOeROF38w2uU1uM1EBEeNSWYXqryiM
Ram0SJggfDhJAITmqwwgWBMXEwcQgACdVCgtt9ILixPjoPh7DATiSTAIKTBqbT6dimX8CGCOZgu
ckAMZpSk0ulY5I4vHYAgBao4HKQUnJAAmsr58sFSuFLPwYuYAA4JTrYHs5QLFbBlczVf5rc5Rh7
JbAAGYG6D8hXbRkqipa6qkLU+sAAcwDQb0AF0gA

                                   +-+
                                   |a|
                                   +-+
                                  /  '\
                                 /    '.
                              +-+      +-+
                              |b|      |c|
                              +-+    _/+-+
                                    //
                                +-+/
                                |d|
                                +-+.
                               |   '\
                              |"     \.
                            +-+       '+-+
                            |f|        |e|
                            +-+        +-+
                           /
                          /
                       +-+
                       |g|
  /                      +-+
 /
/

        +----+
        |    |
        +----+
          /
      23 /
        /
      +----+
digraph {null3 [shape=point]; "s:1 u:1" -> null3 "s:1 u:1" -> "s:2 u:2"[label="2"] "s:2 u:2" -> "s:1 u:3"[label="1"] null4 [shape=point]; "s:2 u:2" -> null4;}

bash -c 'echo "digraph {1->2}" | graph-easy --as_ascii'



