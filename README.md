# printer-server

### 配置国内代理
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

### 获取本地打印机列表自动打印文件

### 程序运行

```
---------------------------------------------------------------
                SubERP Printer Client
                Http Listen On 0.0.0.0:8080
---------------------------------------------------------------
API DOC
1. /print       GET     Params:token    POST Params:file,printname
2. /printlist   GET     Params:token（默认值为suberp）
---------------------------------------------------------------
Design By Baike wangbaike168@qq.com @2023.03

```
