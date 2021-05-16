package apis

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/88250/gulu"
	"github.com/88250/pipe/model"
	"github.com/88250/pipe/service"
	"github.com/88250/pipe/util"
	"github.com/gin-gonic/gin"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

func Login(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)
	openID := c.Query("openid")
	user := service.User.GetUserByName(openID)
	result.Data = user
}

func Register(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		result.Code = util.CodeErr
		return
	}
	if err := service.User.AddUser(user); err != nil {
		result.Code = util.CodeErr
		return
	}
	result.Data = user
}

func GetComment(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)
	ArticleId, _ := strconv.ParseUint(c.Query("articleId"), 10, 64)
	page, _ := strconv.Atoi(c.Query("p"))

	replyComments, pageinfo := service.Comment.GetArticleComments(ArticleId, page, 1)
	var replies []*model.ThemeReply
	for _, replyComment := range replyComments {
		commentAuthor := service.User.GetUser(replyComment.AuthorID)
		if nil == commentAuthor {
			logger.Errorf("not found comment author [userID=%d]", replyComment.AuthorID)
			continue
		}
		commentAuthorBlog := service.User.GetOwnBlog(commentAuthor.ID)
		blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, commentAuthorBlog.ID)
		commentAuthorURL := blogURLSetting.Value + util.PathAuthors + "/" + commentAuthor.Name
		author := &model.ThemeAuthor{
			Name:      commentAuthor.Name,
			URL:       commentAuthorURL,
			AvatarURL: commentAuthor.AvatarURLWithSize(64),
		}

		reply := &model.ThemeReply{
			ID:        replyComment.ID,
			Content:   template.HTML(util.Markdown(replyComment.Content).ContentHTML),
			Author:    author,
			CreatedAt: replyComment.CreatedAt.Format("2006-01-02"),
		}
		replies = append(replies, reply)
	}
	data := make(map[string]interface{})
	data["comments"] = replies
	data["pagination"] = pageinfo
	result.Data = data
}

func AddComment(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)
}

func Static(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)
	ArticleId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	article := service.Article.ConsoleGetArticle(ArticleId)
	if nil == article {
		result.Code = util.CodeErr
		return
	}
	req := map[string][]map[string]interface{}{
		"data": {
			{
				"count": 1,
				"url":   "https://www.jrrm.top/blogs/zhaobingchun" + article.Path,
			},
		},
	}
	bys, _ := json.Marshal(&req)
	resp, err := http.Post("https://ld246.com/uvstat/get", "application/json", bytes.NewReader(bys))
	if err != nil {
		result.Code = util.CodeErr
		return
	}
	var respData map[string]interface{}
	bys, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(bys, &respData)
	result.Code = int(respData["code"].(float64))
	data, _ := respData["data"].(map[string]float64)
	result.Data, _ = data["https://www.jrrm.top/blogs/zhaobingchun"+article.Path]
}
