package file

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"goweb/base"
	"goweb/entity"
	"goweb/util"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func Upload(context *gin.Context) {
	//默认最大8m
	file, header, _ := context.Request.FormFile("file")
	cover := context.Request.FormValue("cover")
	db := context.MustGet("db").(*gorm.DB)
	coverId := context.Request.FormValue("coverId")
	coverDb := entity.Cover{}
	var err error
	coverDb.Id, err = strconv.ParseInt(coverId, 10, 64)
	if "" != coverId && nil != err && coverDb.Id > 0 {
		db.First(&coverDb)
		cover = coverDb.Title
	}
	if "" == cover {
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "请选择对应目录", nil)
		return
	}
	logger := context.MustGet("log").(*log.Logger)
	id := context.MustGet("userId").(int64)
	logger.Println(header.Filename)

	if isPic, _ := CheckFileSuffix(header.Filename); isPic != true {
		logger.Printf("upload file is %s ,not support", header.Filename)
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "文件格式不支持", nil)
		return
	}

	path := strings.Join([]string{"./upload/", strconv.FormatInt(id, 10), "/", cover, "/"}, "")
	logger.Printf("upload file %s path  %s ", header.Filename, path)

	ret, e := uploadImageServer(file, header, logger)
	if nil != e {
		panic(e)
	}
	//不能使用 binary提交方式 会有异常
	go func() {
		funcName(path, header, file)
	}()

	uploadLogs := saveUploadLog(id, ret, path, header)
	if db.Create(&uploadLogs) != nil {
		base.Response(context, base.SUCCESS, uploadLogs)
	} else {
		base.Response(context, base.SUCCESS, uploadLogs)
	}
	return
}

/**
保存上傳日志
*/
func saveUploadLog(id int64, ret map[string]interface{}, path string, header *multipart.FileHeader) entity.UploadLog {
	var uploadLogs entity.UploadLog
	uploadLogs.UserId = id
	uploadLogs.CreateTime = time.Now()
	if ret["deleteUrl"] != nil {
		uploadLogs.DeleteUrl = ret["deleteUrl"].(string)
	} else {
		uploadLogs.DeleteUrl = ""
	}
	if ret["imageUrl"] != nil {
		uploadLogs.Url = ret["imageUrl"].(string)
	} else {
		uploadLogs.Url = ""
	}
	if ret["size"] != nil {
		size := ret["size"].(float64)
		uploadLogs.Length = util.Wrap(size, 0)
	} else {
		uploadLogs.Length = 0
	}
	uploadLogs.Path = path
	uploadLogs.Name = header.Filename
	return uploadLogs
}

/**
上传到图床
*/
func uploadImageServer(file multipart.File, header *multipart.FileHeader, logger *log.Logger) (map[string]interface{}, error) {
	url := "https://sm.ms/api/v2/upload"
	// 创建表单文件
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("smfile", header.Filename)
	if err != nil {
		logger.Printf("Create form file failed: %s\n", err)
	}

	// 从文件读取数据，写入表单
	_, err = io.Copy(formFile, file)
	if err != nil {
		logger.Printf("Write to form file falied: %s\n", err)
	}

	contentType := writer.FormDataContentType()
	// 发送之前必须调用Close()以写入结尾行
	err = writer.Close()
	if err != nil {
		logger.Printf("write failed: %s\n", err)
	}
	req, err := http.NewRequest("POST", url, buf)
	if nil != req {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Printf("Post failed: %s\n", err)
		return nil, err
	} else {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dat); err == nil {
		logger.Printf("upload server image response %v", dat)
		if dat["code"] == "success" {
			ret := make(map[string]interface{}, 2)
			data := dat["data"].(map[string]interface{})
			ret["imageUrl"] = data["url"]
			ret["deleteUrl"] = data["delete"]
			ret["size"] = data["size"]
			return ret, nil
		} else {
			if strings.Contains(dat["message"].(string), "Image upload repeated limit") {
				return nil, base.Err(-1, "上传图片重复")
			} else if strings.Contains(dat["message"].(string), "File is too large.") {
				return nil, base.Err(-1, "上传图片过大")
			} else if strings.Contains(dat["message"].(string), "Upload file frequency limit.") {
				return nil, base.Err(-1, "图床上传超过免费限制")
			}
			return nil, base.Err(-1, "图床异常")
		}

	} else {
		logger.Println(err)
		return nil, err
	}
}

/**
上传文件到本地 备份
*/
func funcName(path string, header *multipart.FileHeader, file multipart.File) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		_ = os.MkdirAll(path, 0777) //0777也可以os.ModePerm
		_ = os.Chmod(path, 0777)
	}
	f, _ := os.Create(path + header.Filename)
	defer f.Close()
	_, _ = io.Copy(f, file)
}

/**
校验文件格式
*/
func CheckFileSuffix(filename string) (bool, string) {
	suffix := path.Ext(filename)
	if (suffix == ".jpeg") || (suffix == ".jpg") ||
		(suffix == ".png") || (suffix == ".gif") ||
		(suffix == ".bmp") {
		return true, suffix
	}

	return false, suffix
}
