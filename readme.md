# sebar
A distributed in-memory key-value store. Fully written in Golang

# Usage
```
server -start -port 8005
server -start -port 8005 -master server1
```

## Start/Stop Parameters
 -start 	to start
 -stop		to stop
 -port		define port
 -master	define master sebar server
 -memsize	define maximal meory size

## Rest 
__/r/list			-list all variables__

__/r/add 			-add variable__
add parameter
- public bool qs
- cluster string qs 
- key string qs
- expiry int
- type string qs - support for only standard type, other than standard will be stored as interface{}
- value string qs or payload

__/r/get			- get variable value__
get parameter
- owner string qs default public
- cluster string qs default common
- key string qs

```go
return
{
	found: bool,
	location: string, 	// will be server node location if not in the same server
	value: []byte 		// []byte value of the data 
}
```

## Auth
__/r/login__ 
- user
- pass

```go
return {
	status: OK/NOK
	message: string
	token: string
	expiry: datetime
}
```

__/r/logout__
- user
- token


## Query
__/r/query__
- op string 	save/exec
- name string qs
- query payload

```go
return {
	status: OK/NOK
	message: string
	data: []byte
}
```

Sample of query
```sql
From("public.table1").Group("")
```