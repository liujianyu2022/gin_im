## initialize project

project-root/  
├── config/  
│   ├── config.go  
│   └── config.yaml  
├── model/  
│   └── user.go  
├── repository/  
│   └── user.go  
├── service/  
│   └── user.go  
├── router/  
│   └── router.go  
├── wire/  
│   ├── wire.go  
│   └── wire_gen.go  
├── handle/  
│   └── user.go  
├── db/  
│   ├── mysql.go  
│   └── redis.go  
├── main.go  
├── go.mod  
└── go.sum  

```
go mod init gin_im
go install github.com/google/wire/cmd/wire@latest
```
