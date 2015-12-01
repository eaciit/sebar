# sebar
A distributed in-memory key-value store. Fully written in Golang

# Usage
```
 
 server -start -port 8005
 server -start -port 8005 -master server1

```

# Start/Stop Parameters
```go
 -start 	to start
 -stop		to stop
 -port		define port
 -master	define master sebar server
 -memsize	define maximal meory size
```

# Rest 
```go
/r/list			-list all variables

/r/add 			-add variable
add parameter
- public bool qs
- cluster string qs 
- key string qs
- expiry int
- type string qs - support for only standard type, other than standard will be stored as interface{}
- value string qs or payload

/r/get			- get variable value
get parameter
- owner string qs default public
- cluster string qs default common
- key string qs

return
{
	found: bool,
	location: string, 	// will be server node location if not in the same server
	value: []byte 		// []byte value of the data 
}
```

# Auth
```go
/r/login 
- user
- pass

return {
	status: OK/NOK
	message: string
	token: string
	expiry: datetime
}

/r/logout
- user
- token
```

# Query
```go
/r/query
- op string 	save/exec
- name string qs
- query payload

return {
	status: OK/NOK
	message: string
	data: []byte
}

sample of query
From("public.table1").Group("")
``

