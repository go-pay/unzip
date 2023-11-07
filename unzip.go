package unzip

import (
	"context"
)

// 解压指定文件
func DecompressFile() {

}

// 远程解压指定文件
func DecompressFileFromURL(c context.Context, zipUrl string, files []string, saveDir ...string) (err error) {
	// read zip file head
	bs, err := readZipFileHead(c, zipUrl)
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

// 远程读取指定文件
func ReadFileFromURL(c context.Context, zipUrl string, files []string) (fileContent map[string][]byte, err error) {
	// read zip file head
	bs, err := readZipFileHead(c, zipUrl)
	if err != nil {
		return nil, err
	}
	// findFiles
	efs, err := findFiles(c, zipUrl, bs, files, 65536)
	if err != nil {
		return nil, err
	}
	fileContent = make(map[string][]byte)
	for _, v := range efs {
		//xlog.Infof("v: %#v", v)
		fileStream, err := readFile(c, zipUrl, v)
		if err != nil {
			return nil, err
		}
		fileContent[v.FileName] = fileStream
	}
	return fileContent, nil
}
