package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func (c StorageClient) computeHmac256(message string) string {
	h := hmac.New(sha256.New, c.accountKey)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func currentTimeRfc1123Formatted() string {
	const dateLayout = http.TimeFormat // reuse from net/http package
	return timeRfc1123Formatted(time.Now().UTC())
}

func timeRfc1123Formatted(t time.Time) string {
	return t.Format(http.TimeFormat)
}

func mergeParams(v1, v2 url.Values) url.Values {
	out := url.Values{}
	for k, v := range v1 {
		out[k] = v
	}
	for k, v := range v2 {
		vals, ok := out[k]
		if ok {
			vals = append(vals, v...)
			out[k] = vals
		} else {
			out[k] = v
		}
	}
	return out
}

func prepareBlockListRequest(blocks []Block) string {
	s := `<?xml version="1.0" encoding="utf-8"?><BlockList>`
	for _, v := range blocks {
		s += fmt.Sprintf("<%s>%s</%s>", v.Status, v.Id, v.Status)
	}
	s += `</BlockList>`
	return s
}

func xmlUnmarshal(body io.ReadCloser, v interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	defer body.Close()
	return xml.Unmarshal(data, v)
}
