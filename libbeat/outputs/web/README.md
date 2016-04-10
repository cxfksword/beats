use https://github.com/jteeuwen/go-bindata to generate static.go

```
cd ../libbeat/outputs/web 
go-bindata -o static.go -pkg web static
```