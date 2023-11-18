package unzip

import "context"

// DownloadRemoteFile 远程解压指定文件
// files: 需要解压的文件(完整路径)
func DownloadRemoteFile(c context.Context, zipUrl string, files []string, saveDir ...string) (err error) {
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

// DownloadRemoteFile 远程解压指定文件
func (zr *ZipReader) DownloadRemoteFile(c context.Context, files []string, saveDir ...string) (err error) {
	if zr == nil {
		return ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return ErrZipReaderDirectory
	}
	for _, f := range files {
		retFiles := zr.findFileNode(zr.directory, f)
		if len(retFiles) <= 0 {
			return NotFoundZipFile
		}
		for _, rf := range retFiles {
			err = zr.downLoadFile(c, rf.file, saveDir...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
