/**
 * copy from: https://zh.mojotv.cn/go/go-range-download
 */
package download

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// FileDownloader 文件下载器
type FileDownloader struct {
	fileSize            int
	url                 string
	outputFileName      string
	totalPart           int //下载线程
	outputDir           string
	doneFilePart        []filePart
	verifyFileValue     string
	verifyFileIntegrity func([]byte) string
	header              map[string]string
	cookieKey           string // 默认为 token
	cookieValue         string
}

// filePart 文件分片
type filePart struct {
	Index int    //文件分片的序号
	From  int    //开始byte
	To    int    //解决byte
	Data  []byte //http下载得到的文件内容
}

// NewFileDownloader .
func NewFileDownloader(url, outputFileName, outputDir, verifyFileValue string,
	totalPart int,
	verifyFileIntegrity func([]byte) string,
	header map[string]string,
	cookieKey, cookieValue string) *FileDownloader {
	if outputDir == "" {
		wd, err := os.Getwd() //获取当前工作目录
		if err != nil {
			fmt.Println(err)
		}
		outputDir = wd
	}
	return &FileDownloader{
		fileSize:            0,
		url:                 url,
		outputFileName:      outputFileName,
		outputDir:           outputDir,
		totalPart:           totalPart,
		doneFilePart:        make([]filePart, totalPart),
		verifyFileValue:     verifyFileValue,
		verifyFileIntegrity: verifyFileIntegrity,
		header:              header,
		cookieKey:           cookieKey,
		cookieValue:         cookieValue,
	}
}

func parseFileInfoFrom(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)

		if err != nil {
			panic(err)
		}
		return params["filename"]
	}
	filename := filepath.Base(resp.Request.URL.Path)
	return filename
}

// head 获取要下载的文件的基本信息(header) 使用HTTP Method Head
func (d *FileDownloader) head() (int, error) {
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode > 299 {
		return 0, errors.New(fmt.Sprintf("Can't process, response is %v", resp.StatusCode))
	}
	//检查是否支持 断点续传
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return 0, errors.New("服务器不支持文件断点续传")
	}

	if len(d.outputFileName) == 0 {
		d.outputFileName = parseFileInfoFrom(resp)
	}

	if len(d.outputFileName) == 0 {
		return 0, errors.New("未知的文件名")
	}

	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	return strconv.Atoi(resp.Header.Get("Content-Length"))
}

// Run 开始下载任务
func (d *FileDownloader) Run() error {
	fileTotalSize, err := d.head()
	if err != nil {
		return err
	}
	d.fileSize = fileTotalSize

	jobs := make([]filePart, d.totalPart)
	eachSize := fileTotalSize / d.totalPart

	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].From = 0
		} else {
			jobs[i].From = jobs[i-1].To + 1
		}
		if i < d.totalPart-1 {
			jobs[i].To = jobs[i].From + eachSize
		} else {
			//the last filePart
			jobs[i].To = fileTotalSize - 1
		}
	}

	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(job filePart) {
			defer wg.Done()
			err := d.downloadPart(job)
			if err != nil {
				fmt.Println("下载文件失败:", err, job)
			}
		}(j)

	}
	wg.Wait()
	return d.mergeFileParts()
}

// 下载分片
func (d FileDownloader) downloadPart(c filePart) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	fmt.Printf("开始[%d]下载from:%d to:%d\n", c.Index, c.From, c.To)
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c.From, c.To))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("服务器错误状态码: %v", resp.StatusCode))
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(bs) != (c.To - c.From + 1) {
		return errors.New("下载文件分片长度错误")
	}
	c.Data = bs
	d.doneFilePart[c.Index] = c
	return nil

}

// getNewRequest 创建一个request
func (d FileDownloader) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.url,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if len(d.header) > 0 {
		for k, v := range d.header {
			r.Header.Set(k, v)
		}
	}
	if len(d.cookieValue) > 0 {
		r.AddCookie(&http.Cookie{Name: d.cookieKey, Value: d.cookieValue})
	}
	r.Header.Set("User-Agent", "knife client")
	return r, nil
}

// mergeFileParts 合并下载的文件
func (d FileDownloader) mergeFileParts() error {
	fmt.Println("开始合并文件")
	path := filepath.Join(d.outputDir, d.outputFileName)
	mergedFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer mergedFile.Close()

	mergeBuffer := bytes.Buffer{}
	for _, s := range d.doneFilePart {
		mergeBuffer.Write(s.Data)
	}

	if d.verifyFileIntegrity == nil {
		mergedFile.Write(mergeBuffer.Bytes())
		return nil
	}

	if len(d.verifyFileValue) != 0 && d.verifyFileIntegrity(mergeBuffer.Bytes()) == d.verifyFileValue {
		fmt.Println("文件校验成功")
		mergedFile.Write(mergeBuffer.Bytes())
	} else {
		return errors.New("文件损坏")
	}
	return nil
}
