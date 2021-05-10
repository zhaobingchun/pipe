package apis

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetOpenID(c *gin.Context) {
	var (
		appid   string
		secret  string
		address string
	)
	platform := c.Query("platform")
	switch platform {
	case "wx":
		appid = os.Getenv("WX_APPID")
		secret = os.Getenv("WX_SECRET")
		address = "https://api.weixin.qq.com/sns/oauth2/access_token"
	}
	code := c.Query("code")
	url := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=%s",
		address, appid, secret, code, "authorization_code")
	bys, err := httpGet(url)
	if err != nil {
		logger.Error("get wx open id", err)
	}
	c.Data(http.StatusOK, "application/json", bys)
}

func httpGet(url string) (bys []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("http code [%d]", resp.StatusCode)
		return
	}
	bys, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	logger.Debugf("get url[%s], resp[%s]", url, string(bys))
	return
}
