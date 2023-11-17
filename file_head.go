package unzip

// LocalFileHead 本地文件头组成（Local File Header）
// /**
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
// format: "IHHHIIIIHH", []string{"I", "H", "H", "H", "I", "I", "I", "I", "H", "H"}
type LocalFileHead struct {
	Signature        uint32 // 本地文件头标识符 (0x04034b50(大端))
	NeedVersion      uint16 // 提取需要的版本
	Flag             uint16 // 通用位标志
	Method           uint16 // 压缩方法
	LastModTime      uint16 // 最后修改文件时间 时分秒
	LastModDate      uint16 // 最后修改文件日期 年月日
	Crc32            uint32 // crc-32(对压缩前的文件计算)
	CompressedSize   uint32 // 压缩大小
	UncompressedSize uint32 // 未压缩大小
	FileNameLen      uint16 // 文件名长度
	ExtraLen         uint16 // 额外字段长度

	FileName string // 文件名
}

// CentralDirFileHead 中央目录文件头（Central Directory File Header）
// /**
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
// format "IHHHHIIIIHHHHHHII", []string{"I", "H", "H", "H", "H", "I", "I", "I", "I", "H", "H", "H", "H", "H", "I", "I"}
type CentralDirFileHead struct {
	Signature        uint32 // 头标识符 (0x04034b50(大端))
	Version          uint16 // 版本
	NeedVersion      uint16 // 提取需要的版本
	Flag             uint16 // 通用位标志
	Method           uint16 // 压缩方法
	LastModTime      uint16 // 最后修改文件时间 时分秒
	LastModDate      uint16 // 最后修改文件日期 年月日
	Crc32            uint32 // crc-32(对压缩前的文件计算)
	CompressedSize   uint32 // 压缩大小
	UncompressedSize uint32 // 压缩前大小
	FileNameLen      uint16 // 文件名长度
	ExtraLen         uint16 // 额外字段长度
	CommentLen       uint16 // 文件注释长度
	DiskNumberStart  uint16 // 磁盘号
	InternalAttr     uint16 // 内部文件属性
	ExternalAttr     uint32 // 外部文件属性
	HeaderOffset     uint32 // 本地头的相对偏移量 (对应的本地文件相对于文件开始的偏移量)

	FileName string // 文件名
}
