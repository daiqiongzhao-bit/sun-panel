package models

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type Backup struct {
	BaseModel
	Name        string `gorm:"type:varchar(255)" json:"name"`
	FilePath    string `gorm:"type:varchar(500)" json:"filePath"`
	Mode        int    `gorm:"type:tinyint(1)" json:"mode"`
	Size        int64  `json:"size"`
	CreatorId   uint   `gorm:"index" json:"creatorId"`
	Status      int    `gorm:"type:tinyint(1);default:1" json:"status"`
}

type BackupTask struct {
	BaseModel
	Name          string `gorm:"type:varchar(100)" json:"name"`
	CronExpr      string `gorm:"type:varchar(100)" json:"cronExpr"`
	Mode          int    `gorm:"type:tinyint(1)" json:"mode"`
	RetentionDays int    `gorm:"default:7" json:"retentionDays"`
	Status        int    `gorm:"type:tinyint(1);default:1" json:"status"`
}

func (m *Backup) GetList(page, pageSize int, mode int) ([]Backup, int64, error) {
	var list []Backup
	var count int64
	db := Db.Model(m)
	if mode > 0 {
		db = db.Where("mode=?", mode)
	}
	offset := (page - 1) * pageSize
	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("created_at desc").Find(&list).Error
	return list, count, err
}

func (m *Backup) GetById(id uint) (Backup, error) {
	var backup Backup
	err := Db.Where("id=?", id).First(&backup).Error
	return backup, err
}

func (m *Backup) Create() (Backup, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *Backup) Delete(id uint) error {
	return Db.Delete(m, id).Error
}

func (m *Backup) UpdateStatus(id uint, status int) error {
	return Db.Model(m).Where("id=?", id).Update("status", status).Error
}

func (m *BackupTask) GetList() ([]BackupTask, error) {
	var list []BackupTask
	err := Db.Order("id asc").Find(&list).Error
	return list, err
}

func (m *BackupTask) GetById(id uint) (BackupTask, error) {
	var task BackupTask
	err := Db.Where("id=?", id).First(&task).Error
	return task, err
}

func (m *BackupTask) Create() (BackupTask, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *BackupTask) Update(id uint) error {
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":           m.Name,
		"cron_expr":      m.CronExpr,
		"mode":           m.Mode,
		"retention_days": m.RetentionDays,
		"status":         m.Status,
	}).Error
}

func (m *BackupTask) Delete(id uint) error {
	return Db.Delete(m, id).Error
}

func BackupDatabase(sourceFile, filePath string) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func BackupFullData(sourceFile, filePath string) error {
	// 全量备份：打包数据库文件 + conf 目录
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	// 备份数据库
	addFileToZip(zw, sourceFile, "database.db")

	// 备份配置文件目录
	confDir := "conf"
	if _, err := os.Stat(confDir); err == nil {
		filepath.Walk(confDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			addFileToZip(zw, path, path)
			return nil
		})
	}

	return nil
}

func addFileToZip(zw *zip.Writer, filePath, zipName string) error {
	src, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	w, err := zw.Create(zipName)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, src)
	return err
}

func RestoreDatabase(filePath string) error {
	return nil
}

func RestoreFullData(filePath string) error {
	return nil
}

func RegisterCronTask(id uint, cronExpr string) error {
	return nil
}

func UnregisterCronTask(id uint) error {
	return nil
}