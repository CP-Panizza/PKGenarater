# PKGenarater
mysql table primary_key genarater

# usage
```go
pk := NewPKGenarater("user", 10, "id", "User")
err := pk.Init()
if err != nil {
    panic(err)
}
result, _ := pk.GetPK()
println(result)   //result: User000001
```
	
