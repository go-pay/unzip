package unzip

import (
	"context"
	"testing"

	"github.com/go-pay/xlog"
)

var (
	ctx = context.Background()
)

func TestDecompressFileFromURL(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "https://pay.wechatpay.cn/wiki/doc/apiv3/wechatpay/download/Product_5.zip"
	DownloadRemoteFile(ctx, zipUrl, []string{"Product/Qt5Core.dll", "Product/Qt5Gui.dll", "Product/Qt5Widgets.dll"}, "/Users/jerry/file")
}
