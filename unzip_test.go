package unzip

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-pay/xlog"
)

var (
	ctx = context.Background()
)

func TestDecompressFileFromURL(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	//zipUrl := "https://test.cdn.sunmi.com/OTA/15008833235745.zip"
	//zipUrl := "https://test.cdn.sunmi.com/OTA/15052081553365.zip"
	zipUrl := "https://cdn.test.sunmi.com/temp/generalfile/hardware/ded3f85f5362401a8bba9807685bd5f1.zip"
	//zipUrl := "https://cdn.test.sunmi.com/temp/generalfile/hardware/93f292c855cb49379f9b3dcfe99c90f9.zip"
	//zipUrl := "https://test.cdn.sunmi.com/OTA%5C15331150971730.zip"
	//zipUrl := "https://cdn.test.sunmi.com/temp/generalfile/hardware/75f01ce5197e455aba98ef24aaab280d.zip"
	//zipUrl := "https://cdn.test.sunmi.com/temp/generalfile/hardware/4d15613e62774708b4db3bd72bd9b02f.zip"
	//err := DecompressFileFromURL(ctx, zipUrl, []string{"version.txt"}, "/Users/jerry/file")
	//if err != nil {
	//	fmt.Println(err)
	//}
	file, _ := ReadFileFromURL(ctx, zipUrl, []string{"75f01ce5197e455aba98ef24aaab2803/version.txt"})
	fmt.Println(string(file["75f01ce5197e455aba98ef24aaab2803/version.txt"]))
}
