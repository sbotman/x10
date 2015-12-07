package service

import (

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
	"github.com/sbotman/x10/api"
	"net"
)

type X10DeviceResource struct {
	db gorm.DB
}

func (tr *X10DeviceResource) CreateAction(c *gin.Context) {
	var device api.X10Device

	if err := c.Bind(&device); err != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}

	conn, err := net.Dial("tcp", "localhost:1099")
	if err != nil {
		c.JSON(400, api.NewError("cannot connect to mochad"))
		return
	}
	defer conn.Close()

    cmd := "pl " + device.Code + " " + device.State + "\n"

	conn.Write([]byte(cmd))

	c.JSON(201, device)
}

func (tr *X10DeviceResource) CreateDevice(c *gin.Context) {
	var device api.X10Device

	if err := c.Bind(&device); err != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	// device.Status = api.TodoStatus
	device.Created = int32(time.Now().Unix())

	tr.db.Save(&device)

	c.JSON(201, device)
}

func (tr *X10DeviceResource) GetAllDevices(c *gin.Context) {
	var devices []api.X10Device

	tr.db.Order("created desc").Find(&devices)

	c.JSON(200, devices)
}

func (tr *X10DeviceResource) GetDevice(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var device api.X10Device

	if tr.db.First(&device, id).RecordNotFound() {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, device)
	}
}

func (tr *X10DeviceResource) UpdateDevice(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var device api.X10Device
	if err = c.Bind(&device); err != nil {
		c.JSON(400, api.NewError("problem decoding json"))
		return
	}
	device.Id = int32(id)

	var existing api.X10Device

	if tr.db.First(&existing, id).RecordNotFound() {
		c.JSON(404, api.NewError("not found"))
	} else {
		tr.db.Save(&device)
		c.JSON(200, device)
	}

}

func (tr *X10DeviceResource) PatchDevice(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	// this is a hack because Gin falsely claims my unmarshalled obj is invalid.
	// recovering from the panic and using my object that already has the json body bound to it.
	var json []api.Patch

	defer func() {
		if r := recover(); r != nil {
			c.JSON(400, api.NewError("problem decoding data"))
		}
	}()

	c.Bind(&json)

	if json[0].Op != "replace" && json[0].Path != "/status" {
		c.JSON(400, api.NewError("PATCH support is limited and can only replace the /status path"))
		return
	}

	var device api.X10Device

	if tr.db.First(&device, id).RecordNotFound() {
		c.JSON(404, api.NewError("not found"))
	} else {
		//device.Status = json[0].Value

		tr.db.Save(&device)
		c.JSON(200, device)
	}

}

func (tr *X10DeviceResource) DeleteDevice(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var device api.X10Device

	if tr.db.First(&device, id).RecordNotFound() {
		c.JSON(404, api.NewError("not found"))
	} else {
		tr.db.Delete(&device)
		c.Data(204, "application/json", make([]byte, 0))
	}
}

func (tr *X10DeviceResource) getId(c *gin.Context) (int32, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int32(id), nil
}
