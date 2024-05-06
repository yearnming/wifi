# WiFi-cracker

WiFi cracker by Golang

> For testing purposes only

## Environment

### Operating System

- Windows

### Language

- Simplified Chinese

You can support other languages by modifying the following variables in `pkg/setting/stat.go`:

- StatText
- SignalText
- AssociatingStatText
- AuthenticatingStatText
- DisconnectingStatText
- DisconnectedStatText
- ConnectedStatText

## Modify password generate parameter

`pkg/config/password.go`

```go
// Maximum password length
PwdMinLen   = 8
// Minimum password length
PwdMaxLen   = 10
// Password characters
PwdCharDict = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
```

## Run

```shell
go run main.go -h
wifi 自动尝试密码

Usage:
C:\Users\test\Desktop\wifi1.exe [flags]

Flags:
INPUT:
-l, -dict string  wifi密码字典


使用示例:
go run main.go -l common.txt
```
