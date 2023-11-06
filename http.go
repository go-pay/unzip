package unzip

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/go-pay/xhttp"
	"github.com/go-pay/xlog"
)

var (
	httpClient = xhttp.NewClient()
)

func httpGetRange(c context.Context, url string, start, getSize int64) (bs []byte, err error) {
	if start < 0 {
		start = 0
	}
	end := start + getSize - 1
	xlog.Infof("getRange: %d-%d", start, end)
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
