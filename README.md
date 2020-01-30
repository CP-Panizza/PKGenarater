# PKGenarater
mysql table primary_key genarater

# usage
```go
pk := NewPKGenarater("user", 10, "id", "user")
err := pk.Init()
if err != nil {
	panic(err)
}
```
	
