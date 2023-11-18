package unzip

import (
	"fmt"
	"testing"

	"github.com/go-pay/xlog"
)

func TestReadFileFromURLByName(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "https://tangboedu-1010.oss-cn-hangzhou.aliyuncs.com/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	fileMap, err := zr.ReadFileByName(ctx, []string{"other.txt"})
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
	fileStream, err := zr.ReadFileByPath(ctx, "/remoteFile/level1/level2/level3/version3.txt")
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

func TestDownZipFileFromURL(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "https://tangboedu-1010.oss-cn-hangzhou.aliyuncs.com/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	err = zr.DownloadRemoteFile(ctx, []string{"version1.txt"}, "/Users/sm3245/Downloads")
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
}

func TestPrintDirectoryFromLocal(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "/Users/sm3245/Downloads/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	err = zr.PrintDirectory()
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
}

func TestReadFileFromLocalByName(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "/Users/sm3245/Downloads/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	fileMap, err := zr.ReadFileByName(ctx, []string{"other.txt"})
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range fileMap {
		fmt.Println("path:", k)
		fmt.Println(string(v))
	}
}

func TestReadFileFromLocalByPath(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	zipUrl := "/Users/sm3245/Downloads/remoteFile.zip"
	zr, err := NewZipReader(ctx, zipUrl)
	if err != nil {
		xlog.Errorf("err:%v", err)
	}
	fileStream, err := zr.ReadFileByPath(ctx, "/remoteFile/level1/level2/level3/version3.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(fileStream))
}
