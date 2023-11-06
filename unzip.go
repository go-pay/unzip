package unzip

import (
	"context"
	"strconv"

	"github.com/go-pay/xlog"
)

// 解压指定文件
func DecompressFile() {

}

// 远程解压指定文件
func DecompressFileFromURL(c context.Context, zipUrl string, files []string, saveDir ...string) (err error) {
	res, err := httpClient.HttpClient.Head(zipUrl)
	if err != nil {
		xlog.Errorf("http head err:%+v", err)
		return
	}
	ar := res.Header.Get("Accept-Ranges")
	if ar != "bytes" {
		xlog.Warnf("http head err:%+v", "Accept-Ranges is not bytes")
		return
	}
	cl, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}
	//xlog.Infof("Content-Length:%d", cl)
	// 获取最后65536字节，zip文件头信息
	bs, err := httpGetRange(c, zipUrl, cl-65536, 65536)
	if err != nil {
		return err
	}
	// findFiles
	efs, err := findFiles(c, zipUrl, bs, files, 65536)
	if err != nil {
		return err
	}
	for _, v := range efs {
		//xlog.Infof("v: %#v", v)
		downLoadFile(c, zipUrl, v, saveDir...)
	}
	return nil
}
