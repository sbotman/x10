package service

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sbotman/x10/api"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
	DbPort	   string
}

type DeviceService struct {
}

func (s *DeviceService) getDb(cfg Config) (gorm.DB, error) {
	connectionString := cfg.DbUser + ":" + cfg.DbPassword + "@tcp(" + cfg.DbHost + ":" + cfg.DbPort + ")/" + cfg.DbName + "?charset=utf8&parseTime=True"

	return gorm.Open("mysql", connectionString)
}

func (s *DeviceService) Migrate(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	db.AutoMigrate(&api.X10Device{})
	return nil
}
func (s *DeviceService) Run(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	deviceResource := &X10DeviceResource{db: db}

	r := gin.Default()
	r.POST("/action", deviceResource.CreateAction)
	r.GET("/device", deviceResource.GetAllDevices)
	r.GET("/device/:id", deviceResource.GetDevice)
	r.POST("/device", deviceResource.CreateDevice)
	r.PUT("/device/:id", deviceResource.UpdateDevice)
	r.PATCH("/device/:id", deviceResource.PatchDevice)
	r.DELETE("/device/:id", deviceResource.DeleteDevice)

	r.Run(cfg.SvcHost)

	return nil
}
