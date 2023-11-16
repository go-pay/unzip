package unzip

import (
	"bytes"
	"compress/flate"
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/go-pay/unzip/zip"
	"github.com/go-pay/xlog"
)

var (
	ErrHead               = errors.New("Accept-Ranges is not bytes")
	ErrZipReader          = errors.New("zip reader is nil")
	ErrZipReaderDirectory = errors.New("zip reader directory is nil")
)

type ZipReader struct {
	r         *zip.Reader
	zipUrl    string
	bs        []byte    // Zip package central directory information
	directory *FileNode // file directory information
}

type FileNode struct {
	file     *ExtractFile // file Details
	filePath string
	isFile   bool
	children []*FileNode
}

func (zr *ZipReader) buildFileNode(parent *FileNode, file *ExtractFile, filePath string) {
	parts := strings.Split(file.FileName, "/")
	if len(parts) == 1 {
		// 文件名没有"/"，直接添加到父节点的子节点列表中
		node := &FileNode{
			file:     file,
			filePath: filePath + "/" + file.FileName, // 更新节点的文件路径
			isFile:   true,
			children: []*FileNode{},
		}
		parent.children = append(parent.children, node)
		return
	}
	// 文件名包含"/"，需要递归处理目录结构
	dirName := parts[0]
	// 只截取parts[0]取文件目录名，保留完整的子文件路径parts[1:]
	childFileName := strings.Join(parts[1:], "/")
	// 后续没有子文件了，返回
	if childFileName == "" {
		return
	}
	childNode := zr.findChildNode(parent.children, dirName)
	if childNode == nil {
		// 目录节点不存在，创建新的目录节点
		childNode = &FileNode{
			file:     &ExtractFile{FileName: dirName},
			filePath: filePath + "/" + dirName, // 更新子节点的文件路径
			isFile:   false,
			children: []*FileNode{},
		}
		parent.children = append(parent.children, childNode)
	}
	// 递归处理子目录和文件
	// todo:test optimize
	extractFile := &ExtractFile{
		FileName:         childFileName,
		Method:           file.Method,
		CompressedSize:   file.CompressedSize,
		UncompressedSize: file.UncompressedSize,
		HeaderOffset:     file.HeaderOffset,
		RangeStart:       file.RangeStart,
		RangeEnd:         file.RangeEnd,
	}
	zr.buildFileNode(childNode, extractFile, filePath+"/"+dirName)
}

func (zr *ZipReader) findChildNode(children []*FileNode, name string) *FileNode {
	for _, child := range children {
		if child.file.FileName == name && !child.isFile {
			return child
		}
	}
	return nil
}

func (zr *ZipReader) findFileNode(node *FileNode, fileName string) []*FileNode {
	var result []*FileNode
	if node.file != nil && node.file.FileName == fileName {
		result = append(result, node)
	}
	children := node.children
	for _, child := range children {
		childResult := zr.findFileNode(child, fileName)
		result = append(result, childResult...)
	}
	return result
}

func (zr *ZipReader) findFileNodeByPath(node *FileNode, filePath string) *FileNode {
	if node.filePath == filePath {
		return node
	}
	for _, child := range node.children {
		if strings.HasPrefix(filePath, child.filePath) {
			return zr.findFileNodeByPath(child, filePath)
		}
	}
	return nil
}

func (zr *ZipReader) printFileNode(node *FileNode, indent string, isLast bool) {
	if node.file == nil {
		return
	}
	// 打印线条指示
	line := "├── "
	if isLast {
		line = "└── "
	}
	// todo:标准化打印
	fmt.Println(indent + line + node.file.FileName) // 打印文件名
	// 为了文件树顺序正确，将文件放在文件夹之前打印
	children := node.children
	sort.Slice(children, func(i, j int) bool {
		return children[i].isFile
	})
	// 计算子节点的缩进和是否为最后一个节点
	childIndent := indent + "│   "
	for i, child := range children {
		isLastChild := i == len(children)-1
		zr.printFileNode(child, childIndent, isLastChild)
	}
}

func (zr *ZipReader) readFile(c context.Context, file *ExtractFile) (fileStream []byte, err error) {
	//xlog.Infof("ReadFile: %+v", file.FileName)
	bs, err := httpGetRange(c, zr.zipUrl, file.RangeStart, file.CompressedSize)
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

func (zr *ZipReader) init(c context.Context, zipUrl string) error {
	// Obtain the central directory section of the zip package
	res, err := httpClient.HttpClient.Head(zipUrl)
	if err != nil {
		return err
	}
	zr.zipUrl = zipUrl
	ar := res.Header.Get("Accept-Ranges")
	if ar != "bytes" {
		return ErrHead
	}
	cl, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}
	// Obtain the last 65536 bytes, zip file header information
	zr.bs, err = httpGetRange(c, zipUrl, cl-65536, 65536)
	if err != nil {
		return err
	}
	// Initialize remote zip file directory
	r, err := zip.NewReader(bytes.NewReader(zr.bs), 65536)
	if err != nil {
		return err
	}
	// 构建文件树头节点
	zr.directory = &FileNode{
		file:     &ExtractFile{FileName: ""},
		filePath: "",
		isFile:   false,
		children: []*FileNode{},
	}
	for _, file := range r.File {
		// xlog.Infof("fileName: %s , method: %d , size: %d , offset: %d", file.Name, file.Method, file.CompressedSize64, file.HeaderOffset) //scatter.txt
		// 收集文件
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

		// todo:test
		// 将item以树形结构存储到zr.directory
		zr.buildFileNode(zr.directory, item, "")
	}
	return nil
}

func NewZipReader(c context.Context, zipUrl string) (zr *ZipReader, err error) {
	zr = new(ZipReader)
	if err = zr.init(c, zipUrl); err != nil {
		return nil, err
	}
	return zr, nil
}

// 打印远端zip文件目录
func (zr *ZipReader) PrintDirectory() error {
	if zr == nil {
		return ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return ErrZipReaderDirectory
	}
	zr.printFileNode(zr.directory.children[0], "", false)
	return nil
}

// todo:多个文件并发操作
// 通过文件名远程读取指定文件，返回key:path value:fileContent(可能存在同名，以不同key:path区分)
func (zr *ZipReader) ReadFileByName(c context.Context, files []string) (fileContent map[string][]byte, err error) {
	if zr == nil {
		return nil, ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return nil, ErrZipReaderDirectory
	}
	fileContent = make(map[string][]byte)
	for _, f := range files {
		retFiles := zr.findFileNode(zr.directory.children[0], f)
		for _, rf := range retFiles {
			fileStream, err := zr.readFile(c, rf.file)
			if err != nil {
				return nil, err
			}
			fileContent[rf.filePath] = fileStream
		}
	}
	return fileContent, nil
}

// 通过完整路径+文件名远程读取指定文件
func (zr *ZipReader) ReadFileByPath(c context.Context, filePath string) (fileContent []byte, err error) {
	if zr == nil {
		return nil, ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return nil, ErrZipReaderDirectory
	}
	retFiles := zr.findFileNodeByPath(zr.directory.children[0], filePath)
	if retFiles == nil {
		return nil, errors.New("file not found")
	}
	fileStream, err := zr.readFile(c, retFiles.file)
	if err != nil {
		return nil, err
	}
	return fileStream, nil
}
