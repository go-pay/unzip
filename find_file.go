package unzip

import (
	"bytes"
	"context"

	"github.com/go-pay/unzip/zip"
	"github.com/go-pay/xlog"
)

type ExtractFile struct {
	FileName         string
	Method           uint16
	CompressedSize   int64
	UncompressedSize int64
	HeaderOffset     int64
	RangeStart       int64
	RangeEnd         int64
}

func findFiles(c context.Context, zipUrl string, bs []byte, files []string, getSize int64) (efs []*ExtractFile, err error) {
	r, err := zip.NewReader(bytes.NewReader(bs), getSize)
	if err != nil {
		return nil, err
	}
	fileNameMap := make(map[string]struct{})
	for _, v := range files {
		fileNameMap[v] = struct{}{}
	}
	for _, file := range r.File {
		xlog.Infof("fileName: %s , method: %d , size: %d , offset: %d", file.Name, file.Method, file.CompressedSize64, file.HeaderOffset) //scatter.txt
		// 收集文件
		if _, ok := fileNameMap[file.Name]; ok {
			item := &ExtractFile{
				FileName:         file.Name,
				Method:           file.Method,
				CompressedSize:   int64(file.CompressedSize64),
				UncompressedSize: int64(file.UncompressedSize64),
				HeaderOffset:     file.HeaderOffset,
			}
			// 获取下载RangeStart
			lfh, _ := getLocalFileHead(c, zipUrl, item.FileName, item.HeaderOffset)
			item.RangeStart = file.HeaderOffset + zip.FileHeaderLen + int64(lfh.FileNameLen+lfh.ExtraLen)
			item.RangeEnd = item.RangeStart + item.CompressedSize - 1
			efs = append(efs, item)
		}
	}
	return efs, nil
}

func findFilePath(c context.Context, zipUrl string, bs []byte, getSize int64) (efs []*ExtractFile, err error) {
	r, err := zip.NewReader(bytes.NewReader(bs), getSize)
	if err != nil {
		return nil, err
	}
	for _, file := range r.File {
		//xlog.Infof("fileName: %s , method: %d , size: %d , offset: %d", file.Name, file.Method, file.CompressedSize64, file.HeaderOffset) //scatter.txt
		// 收集所有文件信息
		item := &ExtractFile{
			FileName:         file.Name,
			Method:           file.Method,
			CompressedSize:   int64(file.CompressedSize64),
			UncompressedSize: int64(file.UncompressedSize64),
			HeaderOffset:     file.HeaderOffset,
		}
		// 获取下载RangeStart
		lfh, _ := getLocalFileHead(c, zipUrl, item.FileName, item.HeaderOffset)
		item.RangeStart = file.HeaderOffset + zip.FileHeaderLen + int64(lfh.FileNameLen+lfh.ExtraLen)
		item.RangeEnd = item.RangeStart + item.CompressedSize - 1
		efs = append(efs, item)
	}
	return efs, nil
}
