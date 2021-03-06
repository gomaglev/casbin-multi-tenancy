package api

import (
	"encoding/json"
	"fmt"
	"gin-casbin/internal/app/ginplus"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ResourceSet
var ResourceSet = wire.NewSet(wire.Struct(new(Resource), "*"))

// Resource
type Resource struct {
}

// Get
func (a *Resource) Get(c *gin.Context) {
	ID := c.Param("id")
	lang := c.Query("lang")

	fileName := fmt.Sprintf("data/json/%s/%s.json", ID, lang)
	content, _ := ioutil.ReadFile(fileName)
	if len(content) == 0 {
		ginplus.ResJSON(c, 200, "")
	}
	var data interface{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		ginplus.ResError(c, err)
	}
	ginplus.ResJSON(c, 200, data)
}
