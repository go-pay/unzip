package unzip

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/go-pay/unzip/zip"
	"github.com/go-pay/xhttp"
	"github.com/go-pay/xlog"
)

func httpGetRange(url string, fileSize, start, getSize int64) {
	if start < 0 {
		start = 0
	}
	xlog.Infof("getRange:%d-%d", start, fileSize-1)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", "bytes="+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(fileSize-1, 10))
	res, err := xhttp.NewClient().HttpClient.Do(req)
	if err != nil {
		xlog.Errorf("http Do err:%+v", err)
		return
	}
	defer res.Body.Close()
	//xlog.Warnf("resp:%+v", res)
	bs, err := io.ReadAll(io.LimitReader(res.Body, int64(1000<<20))) // default 10MB change the size you want
	if err != nil {
		xlog.Errorf("http head err:%+v", err)
		return
	}
	//p := new(Protocol)
	//p.Format = []string{"I", "H", "H", "H", "H", "I", "I", "I", "I", "H", "H", "H", "H", "H", "I", "I"}
	//pack := p.UnPack(bs[62713:])
	//
	//xlog.Infof("bs: %v", pack)

	//xlog.Warnf("bs: %v",string(bs))
	r, err := zip.NewReader(bytes.NewReader(bs), getSize)
	if err != nil {
		xlog.Error(err)
		return
	}
	for _, file := range r.File {
		xlog.Infof("headerOfHfset: %d", file.HeaderOffset)
		xlog.Infof("fileName: %s, size: %d", file.Name, file.CompressedSize64) //scatter.txt
		if file.Name == "system.transfer.list123" {
			//xlog.Infof("fileName: %s, size: %d", file.Name, file.CompressedSize64)
			open, err := file.Open()
			if err != nil {
				xlog.Errorf("file.Open err:%+v", err)
				//if strings.Contains(err.Error(), "negative offset") {
				//	getRange(url, fileSize, fileSize-getSize-65536, getSize)
				//}
				return
			}
			xlog.Infof("open: %v", open)

			bs, err := io.ReadAll(open)
			if err != nil {
				xlog.Error(err)
				return
			}
			xlog.Infof("find over :\n%s", string(bs))

		}
	}
}
