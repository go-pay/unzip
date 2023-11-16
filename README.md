# unzip

### 支持远程指定文件解压，无需下载整个文件

1. ##### 打印远端 ZIP 文件目录

2. ##### 通过文件名远程读取指定文件

3. ##### 通过完整路径+文件名远程读取指定文件

4. todo：更多功能持续更新中

### Install
```
go get github.com/go-pay/unzip
```

### 使用示例
```golang
package main

import (
    "context"
    "fmt"

    "github.com/go-pay/unzip"
)

func main() {
    c := context.Background()
    zipUrl := "https://tangboedu-1010.oss-cn-hangzhou.aliyuncs.com/remoteFile.zip"
    // 从远端读取指定文件
    zr, err := unzip.NewZipReader(c, zipUrl)
    if err != nil {
      fmt.Println(err)
    }
    fileStream, err := zr.FileByPath(c, "/remoteFile/level1/level2/level3/version3.txt")
    if err != nil {
      fmt.Println(err)
    }
    fileContent := string(fileStream)
    fmt.Printf("fileContent: %s\n", fileContent)
}
```