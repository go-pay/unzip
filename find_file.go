package unzip

import (
	"bytes"
	"context"
	"github.com/go-pay/unzip/zip"
	"github.com/go-pay/xlog"
)

type InnerFile struct {
	FileName       string
	CompressedSize int64
	HeaderOffset   int64
	RangeStart     int64
	RangeEnd       int64
}

func findFiles(c context.Context, bs []byte, files []string, getSize int64) (ifs []*InnerFile, err error) {
	r, err := zip.NewReader(bytes.NewReader(bs), getSize)
	if err != nil {
		return nil, err
	}
	fileNameMap := make(map[string]struct{})
	for _, v := range files {
		fileNameMap[v] = struct{}{}
	}
	for _, file := range r.File {
		xlog.Infof("fileName: %s , method: %d , size: %d , offset: %d", file.Name, file.Method, file.CompressedSize64, file.HeaderOffset)                               //scatter.txt
		xlog.Infof("fileName: %s,isDir: %v, fileHeaderLen: 30, extra: %v, len(comment): %d", file.Name, file.FileInfo().IsDir(), string(file.Extra), len(file.Comment)) //scatter.txt
		// 收集文件
		if _, ok := fileNameMap[file.Name]; ok {
			ifs = append(ifs, &InnerFile{
				FileName:       file.Name,
				CompressedSize: int64(file.CompressedSize64),
				HeaderOffset:   file.HeaderOffset,
			})
		}

		//if file.Name == "system.transfer.list123" {
		//	//xlog.Infof("fileName: %s, size: %d", file.Name, file.CompressedSize64)
		//	open, err := file.Open()
		//	if err != nil {
		//		xlog.Errorf("file.Open err:%+v", err)
		//		//if strings.Contains(err.Error(), "negative offset") {
		//		//	getRange(url, fileSize, fileSize-getSize-65536, getSize)
		//		//}
		//		return
		//	}
		//	xlog.Infof("open: %v", open)
		//
		//	bs, err := io.ReadAll(open)
		//	if err != nil {
		//		xlog.Error(err)
		//		return
		//	}
		//	xlog.Infof("find over :\n%s", string(bs))
		//}
	}
	return nil, nil
}
