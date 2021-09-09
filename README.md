# Http Multiplexer


---
```
LISTEN_ADDR=:8080
CONNECTIONS_LIMIT=100
```


---
```
Active routes
    ----
    /process/urls [POST]
    url: []string
    !(len(url) > 20)

    request-example:
        {
            "url": [
                "https://jsonplaceholder.typicode.com/todos/1"
            ]
        }
    
    response-example:
        [
              {
                    "content_type": "application/json; charset=utf-8",
                    "counter": 0,
                    "data": "{\n  \"userId\": 1,\n  \"id\": 1   ,\n  \"title\": \"fugiat veniam minus\",\n  \"completed\": false\n}",
                    "url": "https://jsonplaceholder.typicode.com/todos/3"
              },
        ]
    
    ---
    /process [GET]
    response:
    Testify
```

