#engineers-notebook #golang #go-tooling

### GOOS
```sh
➜  practical-go git:(vishal-kb) ✗ ls 
go.mod hw
➜  practical-go git:(vishal-kb) ✗ 
➜  practical-go git:(vishal-kb) ✗ go run hw
package hw is not in std (/usr/local/go/src/hw)
➜  practical-go git:(vishal-kb) ✗ 
➜  practical-go git:(vishal-kb) ✗ go run hw/hw.go 
Hello world
➜  practical-go git:(vishal-kb) ✗ go run hw      
package hw is not in std (/usr/local/go/src/hw)
➜  practical-go git:(vishal-kb) ✗           
➜  practical-go git:(vishal-kb) ✗ go run hw      
package hw is not in std (/usr/local/go/src/hw)
➜  practical-go git:(vishal-kb) ✗ cd hw
➜  hw git:(vishal-kb) ✗ go run .
Hello world
➜  hw git:(vishal-kb) ✗ 
➜  hw git:(vishal-kb) ✗ go run .
Hello world
➜  hw git:(vishal-kb) ✗ go build .
➜  hw git:(vishal-kb) ✗ ./h
zsh: no such file or directory: ./h
➜  hw git:(vishal-kb) ✗ w
➜  hw git:(vishal-kb) ✗ ./hw
Hello world
➜  hw git:(vishal-kb) ✗ file hw
hw: Mach-O 64-bit executable x86_64
➜  hw git:(vishal-kb) ✗ GOOS=linux go build
➜  hw git:(vishal-kb) ✗ file hw
hw: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, BuildID[sha1]=c6a7625e4597eb1e31257845825ab81fec6c0024, with debug_info, not stripped
```
```
```

### go clean
```sh
practical-go git:(vishal-kb) ✗ ls
go.mod hw
➜  00 practical-go git:(vishal-kb) ✗ ls hw
hw    hw.go
➜  00 practical-go git:(vishal-kb) ✗ go clean
➜  00 practical-go git:(vishal-kb) ✗ ls hw
hw    hw.go
➜  00 practical-go git:(vishal-kb) ✗ cd hw
➜  hw git:(vishal-kb) ✗ ls   
hw    hw.go
➜  hw git:(vishal-kb) ✗ go clean
➜  hw git:(vishal-kb) ✗ ls
hw.go
```