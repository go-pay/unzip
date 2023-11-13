package unzip

import (
	"fmt"
	"github.com/go-pay/xlog"
	"testing"
)

func TestReadFileFromURLByName(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "https://tangboedu-1010.oss-cn-hangzhou.aliyuncs.com/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	fileMap, err := zr.FileByName(ctx, []string{"other.txt"})
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range fileMap {
		fmt.Println("path:", k)
		fmt.Println(string(v))
	}
}

func TestReadFileFromURLByPath(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "https://tangboedu-1010.oss-cn-hangzhou.aliyuncs.com/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	fileStream, err := zr.FileByPath(ctx, "/remoteFile/level1/level2/level3/version3.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(fileStream))
}

func TestPrintDirectoryFromURL(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "https://tangboedu-1010.oss-cn-hangzhou.aliyuncs.com/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	err = zr.PrintDirectory()
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
}
