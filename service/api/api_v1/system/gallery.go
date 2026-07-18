package system

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type GalleryApi struct{}

func (a *GalleryApi) GetList(c *gin.Context) {
	type Request struct {
		PageNum     int    `json:"pageNum"`
		PageSize    int    `json:"pageSize"`
		CategoryId  uint   `json:"categoryId"`
		Type        int    `json:"type"`
		Keyword     string `json:"keyword"`
	}
	req := Request{PageNum: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mGallery := models.Gallery{}
	list, total, err := mGallery.GetList(req.CategoryId, req.Type, req.Keyword, req.PageNum, req.PageSize)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func (a *GalleryApi) GetById(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mGallery := models.Gallery{}
	gallery, err := mGallery.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gallery)
}

func (a *GalleryApi) Upload(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")
	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	fileExt := strings.ToLower(path.Ext(f.Filename))
	agreeExts := []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}
	if !cmn.InArray(agreeExts, fileExt) {
		apiReturn.ErrorByCode(c, 1301)
		return
	}

	fileSize := f.Size
	if fileSize > 10*1024*1024 {
		apiReturn.ErrorByCodeAndMsg(c, 1302, "文件大小不能超过10MB")
		return
	}

	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	fileDir := fmt.Sprintf("%s/gallery/", configUpload)
	isExist, _ := cmn.PathExists(fileDir)
	if !isExist {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	filePath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	c.SaveUploadedFile(f, filePath)

	galleryType := c.PostForm("type")
	galleryTypeInt := 1
	if galleryType == "icon" {
		galleryTypeInt = 2
	}

	gallery := models.Gallery{
		Name:        f.Filename,
		Src:         filePath,
		GalleryType: galleryTypeInt,
		CategoryId:  0,
		UserId:      userInfo.ID,
	}

	created, err := gallery.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{
		"id":        created.ID,
		"name":      created.Name,
		"imageUrl":  filePath[1:],
	})
}

func (a *GalleryApi) Update(c *gin.Context) {
	type Request struct {
		Id          uint   `json:"id"`
		Name        string `json:"name"`
		CategoryId  uint   `json:"categoryId"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	gallery := models.Gallery{
		Name:        req.Name,
		CategoryId:  req.CategoryId,
	}

	err := gallery.Update(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *GalleryApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mGallery := models.Gallery{}
	gallery, err := mGallery.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	os.Remove(gallery.Src)

	err = mGallery.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *GalleryApi) BatchDelete(c *gin.Context) {
	type Request struct {
		Ids []uint `json:"ids"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mGallery := models.Gallery{}
	for _, id := range req.Ids {
		gallery, err := mGallery.GetById(id)
		if err == nil {
			os.Remove(gallery.Src)
		}
	}

	err := mGallery.BatchDelete(req.Ids)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *GalleryApi) GetCategories(c *gin.Context) {
	type Request struct {
		Type int `json:"type"`
	}
	req := Request{Type: 1}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mCategory := models.GalleryCategory{}
	list, err := mCategory.GetList(req.Type)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *GalleryApi) CreateCategory(c *gin.Context) {
	type Request struct {
		Name        string `json:"name" validate:"required"`
		Type        int    `json:"type"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}
	req := Request{Type: 1, Sort: 0}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	category := models.GalleryCategory{
		Name:        req.Name,
		GalleryType: req.Type,
		Sort:        req.Sort,
	}

	created, err := category.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *GalleryApi) UpdateCategory(c *gin.Context) {
	type Request struct {
		Id          uint   `json:"id"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	category := models.GalleryCategory{
		Name: req.Name,
		Sort: req.Sort,
	}

	err := category.Update(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *GalleryApi) DeleteCategory(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mCategory := models.GalleryCategory{}
	err := mCategory.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorByCodeAndMsg(c, 1003, err.Error())
		return
	}

	apiReturn.Success(c)
}