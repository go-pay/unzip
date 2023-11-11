# unzip

支持远程制定文件解压

#### 功能持续完善中...

### Install
```
go get github.com/go-pay/unzip
```

### 使用示例
```
package main

import (
    "context"
    "fmt"

    "github.com/go-pay/unzip"
)

func main() {
    zipUrl := "https://pay.wechatpay.cn/wiki/doc/apiv3/wechatpay/download/Product_5.zip"
    err := unzip.DecompressFileFromURL(context.Background(), zipUrl, []string{"Product/Qt5Core.dll", "Product/Qt5Gui.dll", "Product/Qt5Widgets.dll"}, "/Users/jerry/file")
    if err != nil {
        fmt.Println(err)
        return
    }
}
```