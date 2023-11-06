package unzip

import (
	"context"
	"github.com/go-pay/unzip/unpack"
	"github.com/go-pay/xlog"
	"strconv"
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
	xlog.Infof("Content-Length:%d", cl)
	// 获取最后65536字节
	bs, err := httpGetRange(c, zipUrl, 0, 30)
	if err != nil {
		return err
	}
	xlog.Infof("bs: %v", string(bs))
	// 	fmtEOCD      = "<IHHHHIIH"         // end of central directory
	//	fmtEOCD64    = "<IQHHIIQQQQ"       // end of central directory ZIP64
	//	fmtCDirEntry = "<IHHHHIIIIHHHHHII" // central directory entry
	//	fmtLocalHdr  = "<IHHHIIIIHH"       // local directory header
	p := new(unpack.Protocol)
	// 中央目录文件头（Central Directory File Header）
	///**
	// * 中央目录文件头
	// *
	// * central file header signature   4 bytes      中央文件头标识符 (0x02014b50(大端))
	// * version made by                 2 bytes      版本
	// * version needed to extract       2 bytes      提取所需的版本
	// * general purpose bit flag        2 bytes      通用位标志
	// * compression method              2 bytes      压缩方法
	// * last mod file time              2 bytes      最后修改文件时间  时分秒
	// * last mod file date              2 bytes      最后修改文件日期 年月日
	// * crc-32                          4 bytes      crc-32
	// * compressed size                 4 bytes      压缩大小
	// * uncompressed size               4 bytes      压缩前大小
	// * file name length                2 bytes      文件名长度
	// * extra field length              2 bytes      额外字段长度
	// * file comment length             2 bytes      文件注释长度
	// * disk number start               2 bytes      磁盘号
	// * internal file attributes        2 bytes      内部文件属性
	// * external file attributes        4 bytes      外部文件属性
	// * relative offset of local header 4 bytes      本地头的相对偏移量 (对应的本地文件相对于文件开始的偏移量)
	// *
	// * file name                (variable size)     文件名
	// * extra field              (variable size)     额外字段
	// * file comment             (variable size)     文件注释
	// */
	//p.Format = []string{"I", "H", "H", "H", "H", "I", "I", "I", "I", "H", "H", "H", "H", "H", "I", "I"}

	// 本地文件头组成（Local File Header）
	///**
	// * 本地文件头组成
	// *
	// * local file header signature     4 bytes      本地文件头标识符 (0x04034b50(大端))
	// * version needed to extract       2 bytes      提取需要的版本
	// * general purpose bit flag        2 bytes      通用位标志
	// * compression method              2 bytes      压缩方法
	// * last mod file time              2 bytes      最后修改文件时间 时分秒
	// * last mod file date              2 bytes      最后修改文件日期 年月日
	// * crc-32                          4 bytes      crc-32(对压缩前的文件计算)
	// * compressed size                 4 bytes      压缩大小
	// * uncompressed size               4 bytes      未压缩大小
	// * file name length                2 bytes      文件名长度
	// * extra field length              2 bytes      额外字段长度
	// * file name                 (variable size)    文件名
	// * extra field               (variable size)    额外字段
	// */
	p.Format = []string{"I", "H", "H", "H", "I", "I", "I", "I", "H", "H"}
	pack := p.UnPack(bs[:])
	xlog.Infof("bs: %v", pack)

	// findFiles
	findFiles(c, bs, files, 29)

	return nil
}
