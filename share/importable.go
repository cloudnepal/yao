package share

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/yaoapp/gou"
	"github.com/yaoapp/gou/query/share"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/maps"
	"github.com/yaoapp/xiang/config"
	"github.com/yaoapp/xiang/xlog"
)

// Libs 共享库
var Libs = map[string]map[string]interface{}{}

// Load 加载共享库
func Load(cfg config.Config) {
	LoadFrom(cfg.RootLib)
}

// LoadFrom 从特定目录加载共享库
func LoadFrom(dir string) {

	if DirNotExists(dir) {
		return
	}

	// 加载共享数据
	Walk(dir, ".json", func(root, filename string) {
		name := SpecName(root, filename)
		content := ReadFile(filename)
		libs := map[string]map[string]interface{}{}
		err := jsoniter.Unmarshal(content, &libs)
		if err != nil {
			exception.New("共享数据结构异常 %s", 400, err).Throw()
		}
		for key, lib := range libs {
			key := fmt.Sprintf("%s.%s", name, key)
			Libs[key] = lib
			// 删除注释
			if _, has := lib["__comment"]; has {
				delete(lib, "__comment")
			}
		}
	})

	// 加载共享脚本
	Walk(dir, ".js", func(root, filename string) {
		name := SpecName(root, filename)
		err := gou.JavaScriptVM.Load(filename, name)
		if err != nil {
			xlog.Printf("加载脚本失败 %s", err.Error())
		}
	})
}

// UnmarshalJSON Column 字段JSON解析
func (col *Column) UnmarshalJSON(data []byte) error {
	new := ColumnImp{}
	err := jsoniter.Unmarshal(data, &new)
	if err != nil {
		return err
	}

	// 导入
	err = ImportJSON(new.Import, new.In, &new)
	if err != nil {
		return err
	}

	*col = Column(new)
	return nil
}

// UnmarshalJSON Filter 字段JSON解析
func (filter *Filter) UnmarshalJSON(data []byte) error {
	new := FilterImp{}
	err := jsoniter.Unmarshal(data, &new)
	if err != nil {
		return err
	}

	// 导入
	err = ImportJSON(new.Import, new.In, &new)
	if err != nil {
		return err
	}

	*filter = Filter(new)
	return nil
}

// UnmarshalJSON Render 字段JSON解析
func (render *Render) UnmarshalJSON(data []byte) error {
	new := RenderImp{}
	err := jsoniter.Unmarshal(data, &new)
	if err != nil {
		return err
	}

	// 导入
	err = ImportJSON(new.Import, new.In, &new)
	if err != nil {
		return err
	}

	*render = Render(new)
	return nil
}

// UnmarshalJSON Page 字段JSON解析
func (page *Page) UnmarshalJSON(data []byte) error {
	new := PageImp{}
	err := jsoniter.Unmarshal(data, &new)
	if err != nil {
		return err
	}

	// 导入
	err = ImportJSON(new.Import, new.In, &new)
	if err != nil {
		return err
	}

	*page = Page(new)
	return nil
}

// UnmarshalJSON API 字段JSON解析
func (api *API) UnmarshalJSON(data []byte) error {
	new := APIImp{}
	err := jsoniter.Unmarshal(data, &new)
	if err != nil {
		return err
	}

	// 导入
	err = ImportJSON(new.Import, new.In, &new)
	if err != nil {
		return err
	}

	*api = API(new)
	return nil
}

// ImportJSON 导入
func ImportJSON(name string, in []interface{}, v interface{}) error {
	if name == "" {
		return nil
	}

	lib, has := Libs[name]
	if !has {
		return fmt.Errorf("共享库 %s 不存在", name)
	}

	data := maps.MapStrAny{"$in": in}.Dot()
	content, err := jsoniter.Marshal(share.Bind(lib, data))
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}