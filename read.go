package unzip

import (
	"context"
)

// ReadFileByName todo:多个文件并发操作
// 通过文件名远程读取指定文件，返回key:path value:fileContent(可能存在同名，以不同key:path区分)
func (zr *ZipReader) ReadFileByName(c context.Context, files []string) (fileContentMap map[string][]byte, err error) {
	if zr == nil {
		return nil, ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return nil, ErrZipReaderDirectory
	}
	fileContentMap = make(map[string][]byte)
	for _, f := range files {
		retFiles := zr.findFileNode(zr.directory, f)
		if len(retFiles) <= 0 {
			return nil, NotFoundZipFile
		}
		for _, rf := range retFiles {
			var fileStream []byte
			if len(rf.fileContent) > 0 {
				fileContentMap[rf.filePath] = rf.fileContent
				continue
			}
			if len(zr.localZipData) <= 0 {
				fileStream, err = zr.readRemoteFile(c, rf.file)
			} else {
				fileStream, err = zr.readLocalFile(c, rf.file)
			}
			if err != nil {
				return nil, err
			}
			fileContentMap[rf.filePath] = fileStream
			rf.fileContent = fileStream
		}
	}
	return fileContentMap, nil
}

// ReadFileByPath 通过完整路径+文件名远程读取指定文件
func (zr *ZipReader) ReadFileByPath(c context.Context, filePath string) (fileContent []byte, err error) {
	if zr == nil {
		return nil, ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return nil, ErrZipReaderDirectory
	}
	rf := zr.findFileNodeByPath(zr.directory, filePath)
	if rf == nil {
		return nil, NotFoundZipFile
	}
	if len(rf.fileContent) > 0 {
		return rf.fileContent, nil
	}
	if len(zr.localZipData) <= 0 {
		fileContent, err = zr.readRemoteFile(c, rf.file)
	} else {
		fileContent, err = zr.readLocalFile(c, rf.file)
	}
	if err != nil {
		return nil, err
	}
	rf.fileContent = fileContent
	return fileContent, nil
}
