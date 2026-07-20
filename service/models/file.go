package models

type File struct {
	BaseModel
	Src      string `json:"src"`
	UserId   uint   `json:"userId"`
	FileName string `json:"fileName" gorm:"varchar(255)"` // 文件名
	Method   int    `gorm:"int(5)" json:"method"`         // 上传方式
	Ext      string `gorm:"varchar(255)" json:"ext"`      // 扩展名
	Category string `gorm:"varchar(50);default:'all'" json:"category"` // 分类：icon / wallpaper / all
}

// 添加一个文件记录
func (m *File) AddFile(userId uint, fileName, ext, src string, category ...string) (File, error) {
	cat := "all"
	if len(category) > 0 && category[0] != "" {
		cat = category[0]
	}
	file := File{
		UserId:   userId,
		FileName: fileName,
		Src:      src,
		Ext:      ext,
		Category: cat,
	}
	err := Db.Create(&file).Error

	return file, err
}
