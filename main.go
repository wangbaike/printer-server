package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"path"

	"github.com/AllenDang/w32"
	"golang.org/x/sys/windows/registry"
)


func main() {
	//启动提示信息
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("		SubERP Printer Client						")
	fmt.Println("		Http Listen On 0.0.0.0:8080					")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("API DOC")
	fmt.Println("1. /print	GET	Params:token	POST Params:file,printname")
	fmt.Println("2. /printlist	GET	Params:token")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Design By Baike wangbaike168@qq.com @2023.03")

	//路由
  http.HandleFunc("/print", printHandler)
	http.HandleFunc("/printlist", printerList)
	//http服务开启
    log.Fatal(http.ListenAndServe(":8080", nil))
}

//打印上传的文件
func printHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        file, fh, err := r.FormFile("file")
		printName := r.FormValue("printname")
		
		//token验证
		query := r.URL.Query()
		token := query["token"][0]
		if token != "suberp" {
			fmt.Fprintf(w, "access deny token error")
			return
		}
		
        if err != nil {
            fmt.Println(err)
            return
        }

		defer file.Close()

        // 创建临时文件
        tempFile, err := ioutil.TempFile("", "print*" + path.Ext(fh.Filename))

        if err != nil {
            fmt.Println(err)
            return
        }

		defer tempFile.Close()

        // 将上传的文件内容写入临时文件
        _, err = io.Copy(tempFile, file)

        if err != nil {
            fmt.Println(err)
            return
        }

		file.Close()
		tempFile.Close()
		
        // 调用本地打印机打印文件
		err = w32.ShellExecute(0, "print", tempFile.Name(), printName, "", 0)
		if err != nil {
			fmt.Println(err)
			return
        }

		fmt.Fprintf(w, "ok")

		// 等待一段时间，确保文件已经使用再删除
		time.AfterFunc(time.Second*10, func() {
			err = os.Remove(tempFile.Name());
			if err != nil {
				fmt.Println(err)
				return
       		}
		})

		return
    } else {
        http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    }
}

//获取本地打印机列表
func printerList(w http.ResponseWriter, r *http.Request) {
	//token验证
	query := r.URL.Query()
	token := query["token"][0]
	if token != "suberp" {
		fmt.Fprintf(w, "access deny token error")
		return
	}

    const printers = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Print\Printers`
    k, err := registry.OpenKey(registry.LOCAL_MACHINE, printers, registry.ENUMERATE_SUB_KEYS)
    if err != nil {
        panic(err)
    }
    defer k.Close()
    subKeys, err := k.ReadSubKeyNames(-1)
    if err != nil {
        panic(err)
    }
	var result string
    for _, subKey := range subKeys {
        result +=subKey + "\n"
    }
	fmt.Fprintf(w, result)
}
