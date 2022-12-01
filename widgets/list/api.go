package list

import (
	"fmt"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/yaoapp/gou"
	"github.com/yaoapp/yao/share"
	"github.com/yaoapp/yao/widgets/action"
)

// Guard list widget guard
func Guard(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		abort(c, 400, "the list widget id does not found")
		return
	}

	list, has := Lists[id]
	if !has {
		abort(c, 404, fmt.Sprintf("the list widget %s does not exist", id))
		return
	}

	act, err := list.getAction(c.FullPath())
	if err != nil {
		abort(c, 404, err.Error())
		return
	}

	err = act.UseGuard(c, id)
	if err != nil {
		abort(c, 400, err.Error())
		return
	}

}

func abort(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"code": code, "message": message})
	c.Abort()
}

func (list *DSL) getAction(path string) (*action.Process, error) {

	switch path {
	case "/api/__yao/list/:id/setting":
		return list.Action.Setting, nil
	case "/api/__yao/list/:id/component/:xpath/:method":
		return list.Action.Component, nil
	case "/api/__yao/list/:id/upload/:xpath/:method":
		return list.Action.Upload, nil
	case "/api/__yao/list/:id/download/:field":
		return list.Action.Download, nil
	case "/api/__yao/list/:id/save":
		return list.Action.Save, nil
	}

	return nil, fmt.Errorf("the list widget %s %s action does not exist", list.ID, path)
}

// export API
func exportAPI() error {

	http := gou.HTTP{
		Name:        "Widget List API",
		Description: "Widget List API",
		Version:     share.VERSION,
		Guard:       "widget-list",
		Group:       "__yao/list",
		Paths:       []gou.Path{},
	}

	//   GET  /api/__yao/list/:id/setting  					-> Default process: yao.list.Xgen
	path := gou.Path{
		Label:       "Setting",
		Description: "Setting",
		Path:        "/:id/setting",
		Method:      "GET",
		Process:     "yao.list.Setting",
		In:          []string{"$param.id"},
		Out:         gou.Out{Status: 200, Type: "application/json"},
	}
	http.Paths = append(http.Paths, path)

	//   GET  /api/__yao/list/:id/find  				-> Default process: yao.list.Get $param.id :query
	path = gou.Path{
		Label:       "Get",
		Description: "Get",
		Path:        "/:id/get",
		Method:      "GET",
		Process:     "yao.list.Find",
		In:          []string{"$param.id", ":query-param"},
		Out:         gou.Out{Status: 200, Type: "application/json"},
	}
	http.Paths = append(http.Paths, path)

	//   GET  /api/__yao/list/:id/component/:xpath/:method  	-> Default process: yao.list.Component $param.id $param.xpath $param.method :query
	path = gou.Path{
		Label:       "Component",
		Description: "Component",
		Path:        "/:id/component/:xpath/:method",
		Method:      "GET",
		Process:     "yao.list.Component",
		In:          []string{"$param.id", "$param.xpath", "$param.method", ":query"},
		Out:         gou.Out{Status: 200, Type: "application/json"},
	}
	http.Paths = append(http.Paths, path)

	//   POST  /api/__yao/table/:id/upload/:xpath/:method  	-> Default process: yao.list.Upload $param.id $param.xpath $param.method $file.file
	path = gou.Path{
		Label:       "Upload",
		Description: "Upload",
		Path:        "/:id/upload/:xpath/:method",
		Method:      "POST",
		Process:     "yao.list.Upload",
		In:          []string{"$param.id", "$param.xpath", "$param.method", "$file.file"},
		Out:         gou.Out{Status: 200, Type: "application/json"},
	}
	http.Paths = append(http.Paths, path)

	//   GET  /api/__yao/list/:id/download/:field  	-> Default process: yao.list.Download $param.id $param.xpath $param.field $query.name $query.token
	path = gou.Path{
		Label:       "Download",
		Description: "Download",
		Path:        "/:id/download/:field",
		Method:      "GET",
		Process:     "yao.list.Download",
		In:          []string{"$param.id", "$param.field", "$query.name", "$query.token"},
		Out: gou.Out{
			Status:  200,
			Body:    "{{content}}",
			Headers: map[string]string{"Content-Type": "{{type}}"},
		},
	}
	http.Paths = append(http.Paths, path)

	//  POST  /api/__yao/list/:id/save  						-> Default process: yao.list.Save $param.id :payload
	path = gou.Path{
		Label:       "Save",
		Description: "Save",
		Path:        "/:id/save",
		Method:      "POST",
		Process:     "yao.list.Save",
		In:          []string{"$param.id", ":payload"},
		Out:         gou.Out{Status: 200, Type: "application/json"},
	}
	http.Paths = append(http.Paths, path)

	// api source
	source, err := jsoniter.Marshal(http)
	if err != nil {
		return err
	}

	// load apis
	_, err = gou.LoadAPIReturn(string(source), "widgets.list")
	return err
}
