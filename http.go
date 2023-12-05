package unzip

import (
	"bytes"
	"compress/flate"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-pay/unzip/unpack"
	"github.com/go-pay/unzip/zip"
	"github.com/go-pay/xhttp"
	"github.com/go-pay/xlog"
)

var (
	httpClient = xhttp.NewClient().SetTimeout(0)
)

func httpGetRange(c context.Context, url string, start, getSize int64) (bs []byte, err error) {
	if start < 0 {
		start = 0
	}
	end := start + getSize - 1
	//xlog.Infof("getRange: %d-%d, size: %d", start, end, getSize)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", "bytes="+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10))
	res, err := httpClient.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	//xlog.Warnf("resp:%+v", res)
	//bs, err = io.ReadAll(io.LimitReader(res.Body, int64(10000<<20))) // max 10GB
	bs, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func getLocalFileHead(c context.Context, zipUrl, fileName string, fileOffset int64) (lfh *LocalFileHead, err error) {
	bs, err := httpGetRange(c, zipUrl, fileOffset, zip.FileHeaderLen)
	buf := unpack.ReadBuff(bs)
	sig := buf.Uint32()
	if sig != zip.FileHeaderSignature {
		return nil, zip.ErrFormat
	}
	// 按顺序读取
	version := buf.Uint16()
	flag := buf.Uint16()
	method := buf.Uint16()
	lastModTime := buf.Uint16()
	lastModDate := buf.Uint16()
	crc32 := buf.Uint32()
	compressedSize := buf.Uint32()
	uncompressedSize := buf.Uint32()
	fileNameLen := buf.Uint16()
	extraLen := buf.Uint16()
	lfh = &LocalFileHead{
		Signature:        sig,
		NeedVersion:      version,
		Flag:             flag,
		Method:           method,
		LastModTime:      lastModTime,
		LastModDate:      lastModDate,
		Crc32:            crc32,
		CompressedSize:   compressedSize,
		UncompressedSize: uncompressedSize,
		FileNameLen:      fileNameLen,
		ExtraLen:         extraLen,
		FileName:         fileName,
	}
	//xlog.Infof("LocalFileHead: %+v", lfh)
	return
}

func unpackBuff(c context.Context, fileName string, bs []byte) (lfh *LocalFileHead, err error) {
	buf := unpack.ReadBuff(bs)
	sig := buf.Uint32()
	if sig != zip.FileHeaderSignature {
		return nil, zip.ErrFormat
	}
	// 按顺序读取
	version := buf.Uint16()
	flag := buf.Uint16()
	method := buf.Uint16()
	lastModTime := buf.Uint16()
	lastModDate := buf.Uint16()
	crc32 := buf.Uint32()
	compressedSize := buf.Uint32()
	uncompressedSize := buf.Uint32()
	fileNameLen := buf.Uint16()
	extraLen := buf.Uint16()
	lfh = &LocalFileHead{
		Signature:        sig,
		NeedVersion:      version,
		Flag:             flag,
		Method:           method,
		LastModTime:      lastModTime,
		LastModDate:      lastModDate,
		Crc32:            crc32,
		CompressedSize:   compressedSize,
		UncompressedSize: uncompressedSize,
		FileNameLen:      fileNameLen,
		ExtraLen:         extraLen,
		FileName:         fileName,
	}
	return
}

func downLoadFile(c context.Context, zipUrl string, file *ExtractFile, saveDir ...string) (fileContent []byte, err error) {
	//xlog.Infof("downLoadFile: %+v", file.FileName)
	bs, err := httpGetRange(c, zipUrl, file.RangeStart, file.CompressedSize)
	if err != nil {
		return nil, err
	}
	decompressor := flate.NewReader(bytes.NewBuffer(bs))
	defer decompressor.Close()
	fileContent, err = io.ReadAll(decompressor)
	if err != nil {
		xlog.Errorf("io.ReadAll, err:%+v", err)
		return nil, err
	}
	//xlog.Warnf("file over :\n%s", string(fileContent))
	if len(saveDir) > 0 {
		filePath := saveDir[0] + "/" + file.FileName
		dirPath := filepath.Dir(filePath)
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			xlog.Errorf("os.MkdirAll, err:%+v", err)
			return nil, err
		}
		if err = os.WriteFile(filePath, fileContent, 0666); err != nil {
			xlog.Errorf("os.WriteFile, err:%+v", err)
			return nil, err
		}
	}
	return fileContent, nil
}

func readRemoteFile(c context.Context, zipUrl string, file *ExtractFile) (fileStream []byte, err error) {
	//xlog.Infof("ReadFile: %+v", file.FileName)
	bs, err := httpGetRange(c, zipUrl, file.RangeStart, file.CompressedSize)
	if err != nil {
		return nil, err
	}
	decompressor := flate.NewReader(bytes.NewBuffer(bs))
	defer decompressor.Close()
	fileStream, err = io.ReadAll(decompressor)
	if err != nil {
		xlog.Errorf("io.ReadAll, err:%+v", err)
		return nil, err
	}
	return fileStream, nil
}

func readZipFileHead(c context.Context, zipUrl string) (bs []byte, err error) {
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
		return nil, err
	}
	//xlog.Infof("Content-Length:%d", cl)
	// 获取最后65536字节，zip文件头信息
	bs, err = httpGetRange(c, zipUrl, cl-65536, 65536)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
