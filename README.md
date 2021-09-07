# Http Multiplexer


---
```
Enviroment variables

LISTEN_ADDR=:8080
DEBUG=true // not implemented :|
CONNECTIONS_LIMIT=100
```


---
```
May be open api will better solution, but i think is overhead

Active routes
    ----
    /process/urls [POST]

    {
        "url": [
            "https://jsonplaceholder.typicode.com/todos/1"
        ]
    }
    
    url: []string
    
    !(len(url) > 20)
    response:
    [
          {
                "content_type": "application/json; charset=utf-8",
                "counter": 0,
                "data": "{\n  \"userId\": 1,\n  \"id\": 3,\n  \"title\": \"fugiat veniam minus\",\n  \"completed\": false\n}",
                "url": "https://jsonplaceholder.typicode.com/todos/3"
          },
    ]
    
    ---
    /process [GET]
    response:
    Testify
```

