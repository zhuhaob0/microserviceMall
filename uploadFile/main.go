package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/upload", UploadFile)
	panic(r.Run(":12345"))
}

type Response struct {
	Data    map[string]interface{} `json:"data"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
}

func UploadFile(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil || len(form.File) == 0 {
		return
	}
	var fileName string
	if form, err := ctx.MultipartForm(); err == nil {
		// 获取文件
		files := form.File["file"]
		for _, file := range files {
			// 根据时间戳生成文件名
			fileNameInt := time.Now().Unix()
			idx := strings.LastIndex(file.Filename, ".")
			fileName = file.Filename[:idx] + "-" + strconv.FormatInt(fileNameInt, 10) + file.Filename[idx:]

			filePath := "upload/" + fileName
			err := ctx.SaveUploadedFile(file, filePath)
			if err != nil {
				return
			}
		}
		ctx.JSON(http.StatusOK, Response{gin.H{"filename": fileName}, 200, "文件上传成功！"})
		return
	} else {
		ctx.JSON(http.StatusOK, Response{gin.H{"filename": fileName}, 400, "文件上传失败！"})
		return
	}
}
