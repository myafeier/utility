package utility

import (
	"fmt"
	"github.com/BurntSushi/graphics-go/graphics"
	logger "github.com/myafeier/log"
	"gopkg.in/macaron.v1"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func init()  {
	logger.SetPrefix("utility")
	logger.SetLogLevel(logger.DEBUG)
}

//文件管理器结构
type FileManager struct {
	sync.Mutex
	Infile                  multipart.File
	LocalFile               *os.File
	FileName                string
	FileType                string
	FullPath                string
	ShortPath               string
	ThumbWidth              int
	ThumbHeight             int
	NormalUrl               string
	ThumbUrl                string
	UploadDir               string
	SiteRootPath            string
	SiteStaticUrl           string
	SitePublicDir           string
	SiteUploadDir           string
	SiteMaxFileUploadSizeMb int
	SiteMaxFileNumber       int
}

type UploadFile struct {
	FileName      string //文件名
	ThumbFileName string //缩略图文件名
	Url           string //文件url
	ThumbUrl      string //缩略图URL
	FilePath      string //文件存储路径
}

//返回一个文件管理器，可以
/*
 * SiteStaticUrl  静态站点url
 * SiteRootPath 运行目录绝对路径
 * SitePublicDir public目录名称
 * SiteUploadDir upload目录名称
 */
func NewFileManager(SiteStaticUrl, SiteRootPath, SitePublicDir, SiteUploadDir string, SiteMaxFileUploadSizeMb, SiteMaxFileNumber int) (fm *FileManager) {
	fm = new(FileManager)
	fm.SiteStaticUrl = SiteStaticUrl
	if SiteRootPath == "" {
		fm.SiteRootPath, _ = os.Getwd()
	} else {
		fm.SiteRootPath = SiteRootPath
	}
	if SitePublicDir == "" {
		fm.SitePublicDir = "public"
	} else {
		fm.SitePublicDir = SitePublicDir
	}
	if SiteUploadDir == "" {
		fm.SiteUploadDir = "upload"
	} else {
		fm.SiteUploadDir = SiteUploadDir
	}

	if SiteMaxFileNumber == 0 {
		fm.SiteMaxFileNumber = 9
	} else {
		fm.SiteMaxFileNumber = SiteMaxFileNumber
	}
	if SiteMaxFileUploadSizeMb == 0 {
		fm.SiteMaxFileUploadSizeMb = 5
	} else {
		fm.SiteMaxFileUploadSizeMb = SiteMaxFileUploadSizeMb
	}

	return
}

// 返回默认文件管理器
func NewDefaultFileManager(SiteStaticUrl string) (fm *FileManager) {

	return NewFileManager(SiteStaticUrl, "", "", "", 0, 0)
}

//通过multiform 上传多个文件
func (self *FileManager) UploadMultiFileFromMultiForm(ctx *macaron.Context, thumbWidth, thumbHeight int) (files []*UploadFile, err error) {
	if thumbWidth == 0 {
		self.ThumbWidth = 800
	} else {
		self.ThumbWidth = thumbWidth
	}

	if thumbHeight == 0 {
		self.ThumbHeight = 600
	} else {
		self.ThumbHeight = thumbHeight
	}

	ctx.Req.ParseMultipartForm(int64(self.SiteMaxFileUploadSizeMb*self.SiteMaxFileNumber) << 20)
	if ctx.Req.MultipartForm != nil && ctx.Req.MultipartForm.File != nil {
		fhs := ctx.Req.MultipartForm.File["image"]
		num := len(fhs) //文件数
		if num > self.SiteMaxFileNumber {
			err = fmt.Errorf("文件数超过限制：%d", self.SiteMaxFileNumber)
			logger.Error(err.Error())
			return
		}
		for _, fheader := range fhs {
			err = self.GeneralFileInfo(fheader)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			err = self.UploadPic(fheader)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			err = self.thumb()
			if err != nil {
				logger.Error(err.Error())
				return
			}
			f := new(UploadFile)
			f.FileName = self.FileName + self.FileType
			f.ThumbFileName = self.FileName + "_thumb" + self.FileType
			f.Url = self.SiteStaticUrl + string(filepath.Separator) + self.ShortPath + f.FileName
			f.ThumbUrl = self.SiteStaticUrl + string(filepath.Separator) + self.ShortPath + f.ThumbFileName
			f.FilePath = self.FullPath
			files = append(files, f)
		}
	} else {
		err = fmt.Errorf("%s", "无图片上传")
	}
	return
}

func (self *FileManager) GeneralFileInfo(fh *multipart.FileHeader) (err error) {

	originFileName := fh.Filename

	if originFileName == "" {
		err = fmt.Errorf("%s", "文件格式不被允许")
		return
	}

	fileSlice := strings.Split(originFileName, ".")
	self.FileType = "." + strings.ToLower(fileSlice[len(fileSlice)-1])

	if self.FileType != ".png" && self.FileType != ".jpg" && self.FileType != ".jpeg" && self.FileType != ".gif" && self.FileType != ".bmp" && self.FileType != ".flv" {
		err = fmt.Errorf("%s", "文件格式不被允许")
		return
	}

	common := NewCommon("")
	self.FileName = common.GetUuid()

	fileSepatator := string(filepath.Separator)
	self.UploadDir = self.SiteUploadDir
	publicDir := self.SitePublicDir

	if err != nil {
		logger.Error(err.Error())
		return
	}

	self.ShortPath = self.UploadDir + fileSepatator + time.Now().Format("2006-01-02") + fileSepatator
	self.NormalUrl = self.SiteStaticUrl + fileSepatator + self.UploadDir + fileSepatator + time.Now().Format("2006-01-02") + fileSepatator + self.FileName + self.FileType
	self.FullPath = self.SiteRootPath + fileSepatator + publicDir + fileSepatator + self.ShortPath
	self.ThumbUrl = self.NormalUrl
	return

}

//通过multipart 上传一个文件
func (self *FileManager) UploadPic(fh *multipart.FileHeader) (err error) {
	self.Lock()
	defer self.Unlock()
	err = os.MkdirAll(self.FullPath, os.ModePerm)
	if err != nil {
		logger.Error("err in create directory! err:", err)
		return
	}

	toFile, err := os.OpenFile(self.FullPath+self.FileName+self.FileType, os.O_WRONLY|os.O_CREATE, 0666)
	defer toFile.Close()

	if err != nil {
		logger.Error("Error found in create file:", err)
		return
	}
	sFile, err := fh.Open()

	if err != nil {
		return
	}
	io.Copy(toFile, sFile)
	return
}

//通过file multipart接收图片数据，并形成缩略图
func (self *FileManager) UploadSinglePicThenThumbOne(inFile multipart.File) (err error) {
	self.Lock()
	defer self.Unlock()
	err = os.MkdirAll(self.FullPath, os.ModePerm)
	if err != nil {
		logger.Error("err in create directory! err:", err)
		return
	}

	f, err := os.OpenFile(self.FullPath+self.FileName+self.FileType, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		logger.Error("Error found in create file:", err)
		return
	}
	io.Copy(f, inFile)

	//如果是图片,产生缩略图
	if self.FileType == ".png" || self.FileType == ".jpg" || self.FileType == ".jpeg" {
		err = self.thumb()
	}

	return

}

func (self *FileManager) thumb() (err error) {

	f, err := os.Open(self.FullPath + self.FileName + self.FileType)
	defer f.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	simage, name, err := image.Decode(f)

	if err != nil {
		logger.Error("1", err)
		return
	}

	imageWidth := simage.Bounds().Dx()
	imageHeight := simage.Bounds().Dy()
	var width, height int
	if imageWidth > self.ThumbWidth {
		width = self.ThumbWidth
		height = self.ThumbWidth * imageHeight / imageWidth

	} else {
		width = imageWidth
		height = imageHeight

	}
	logger.Debug("origin size,width:%d,height:%d,now width:%d,height:%d",imageWidth,imageHeight,width,height)

	dstFile := image.NewRGBA(image.Rect(0, 0, width, height))

	err = graphics.Thumbnail(dstFile, simage)
	if err != nil {
		logger.Error("2", err)
		return
	}

	tfile, err := os.Create(self.FullPath + self.FileName + "_thumb" + self.FileType)
	defer tfile.Close()
	if err != nil {
		logger.Error("3", err)
		return
	}

	switch name {
	case "png":
		err = png.Encode(tfile, dstFile)
	case "jpeg":
		err = jpeg.Encode(tfile, dstFile, &jpeg.Options{90})
	default:
		err = fmt.Errorf("File is not a image")
	}
	if err != nil {
		logger.Error("4", err)
		return
	}
	self.ThumbUrl = self.SiteStaticUrl + "/" + self.UploadDir + "/" + time.Now().Format("2006-01-02") + "/" + self.FileName + "_thumb" + self.FileType
	return
}
