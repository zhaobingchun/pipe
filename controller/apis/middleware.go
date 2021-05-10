package apis

import (
	"github.com/88250/pipe/util"
	"github.com/gin-gonic/gin"
)

func HandlerAPI(c *gin.Context) {
	my := &util.SessionData{
		UID:   1,
		UName: "zhaobingchun",
		URole: 1,
		BID:   1,
	}
	my.Save(c)
	c.Next()
}
