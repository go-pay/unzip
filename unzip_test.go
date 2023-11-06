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
	//zipUrl := "https://test.cdn.sunmi.com/OTA/15008833235745.zip"
	zipUrl := "https://test.cdn.sunmi.com/OTA/15052081553365.zip"
	//zipUrl := "https://cdn.test.sunmi.com/temp/generalfile/hardware/47422d780a8c468aa20f86de9ee7cd59.zip"
	//zipUrl := "https://cdn.test.sunmi.com/temp/generalfile/hardware/93f292c855cb49379f9b3dcfe99c90f9.zip"
	//zipUrl := "https://test.cdn.sunmi.com/OTA%5C15331150971730.zip"
	DecompressFileFromURL(ctx, zipUrl, []string{"version.txt", "scatter.txt", "boot.img"}, "/Users/jerry/file")
}
