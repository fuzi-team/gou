package gou

import (
	"fmt"
	"strings"

	"github.com/yaoapp/kun/any"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/str"
)

// Caller 运行器
type Caller struct {
	Name    string
	Type    string
	Class   string
	Method  string
	Args    []interface{}
	Handler func(caller *Caller) interface{}
}

// ModelHandlers 模型运行器
var ModelHandlers = map[string]func(caller *Caller) interface{}{
	"find":     callerFind,
	"get":      callerGet,
	"paginate": callerPaginate,
	"create":   callerCreate,
	"update":   callerUpdate,
	"save":     callerSave,
	"delete":   callerDelete,
	"destroy":  callerDestroy,
	"insert":   callerInsert,
}

// NewCaller 创建运行器
func NewCaller(name string, args ...interface{}) *Caller {
	caller := &Caller{Name: name, Args: args}
	caller.extraProcess()
	return caller
}

// Run 运行方法
func (caller *Caller) Run() interface{} {
	return caller.Handler(caller)
}

// extraProcess 解析执行方法  name = "models.user.Find", name = "plugins.user.Login"
// return type=models, name=login, class=user
func (caller *Caller) extraProcess() {
	namer := strings.Split(caller.Name, ".")
	last := len(namer) - 1
	if last < 2 {
		exception.New(
			fmt.Sprintf("Process:%s 格式错误", caller.Name),
			400,
		).Throw()
	}
	caller.Type = strings.ToLower(namer[0])
	caller.Class = strings.ToLower(strings.Join(namer[1:last], "."))
	caller.Method = strings.ToLower(namer[last])
	if caller.Type == "plugins" { // Plugin
		caller.Handler = callerExec
	} else if caller.Type == "models" { // Model
		handler, has := ModelHandlers[caller.Method]
		if !has {
			exception.New("%s 方法不存在", 404, caller.Method).Throw()
		}
		caller.Handler = handler
	}
}

// validateArgs( args )
func (caller *Caller) validateArgNums(length int) {
	if len(caller.Args) < length {
		exception.New(
			fmt.Sprintf("Model:%s%s(args...); 缺少查询参数", caller.Class, caller.Name),
			400,
		).Throw()
	}
}

// callerExec 运行插件中的方法
func callerExec(caller *Caller) interface{} {
	mod := SelectPluginModel(caller.Class)
	res, err := mod.Exec(caller.Method, caller.Args...)
	if err != nil {
		exception.Err(err, 500).Throw()
	}
	return res
}

// callerFind 运行模型 MustFind
func callerFind(caller *Caller) interface{} {
	caller.validateArgNums(2)
	mod := Select(caller.Class)
	params, ok := caller.Args[1].(QueryParam)
	if !ok {
		params = QueryParam{}
	}
	return mod.MustFind(caller.Args[0], params)
}

// callerGet 运行模型 MustGet
func callerGet(caller *Caller) interface{} {
	caller.validateArgNums(1)
	mod := Select(caller.Class)
	params, ok := caller.Args[0].(QueryParam)
	if !ok {
		exception.New("第1个查询参数错误 %v", 400, caller.Args[0]).Throw()
	}
	return mod.MustGet(params)
}

// callerPaginate 运行模型 MustPaginate
func callerPaginate(caller *Caller) interface{} {
	caller.validateArgNums(3)
	mod := Select(caller.Class)
	params, ok := caller.Args[0].(QueryParam)
	if !ok {
		exception.New("第1个查询参数错误 %v", 400, caller.Args[0]).Throw()
	}

	page := any.Of(caller.Args[1]).CInt()
	pagesize := any.Of(caller.Args[2]).CInt()
	return mod.MustPaginate(params, page, pagesize)
}

// callerCreate 运行模型 MustCreate
func callerCreate(caller *Caller) interface{} {
	caller.validateArgNums(1)
	mod := Select(caller.Class)
	row := any.Of(caller.Args[0]).Map().MapStrAny
	return mod.MustCreate(row)
}

// callerUpdate 运行模型 MustUpdate
func callerUpdate(caller *Caller) interface{} {
	caller.validateArgNums(2)
	mod := Select(caller.Class)
	id := caller.Args[0]
	row := any.Of(caller.Args[1]).Map().MapStrAny
	mod.MustUpdate(id, row)
	return nil
}

// callerSave 运行模型 MustSave
func callerSave(caller *Caller) interface{} {
	caller.validateArgNums(1)
	mod := Select(caller.Class)
	row := any.Of(caller.Args[0]).Map().MapStrAny
	return mod.MustSave(row)
}

// callerDelete 运行模型 MustDelete
func callerDelete(caller *Caller) interface{} {
	caller.validateArgNums(1)
	mod := Select(caller.Class)
	mod.MustDelete(caller.Args[0])
	return nil
}

// callerDestroy 运行模型 MustDestroy
func callerDestroy(caller *Caller) interface{} {
	caller.validateArgNums(1)
	mod := Select(caller.Class)
	mod.MustDestroy(caller.Args[0])
	return nil
}

// callerInsert 运行模型 MustInsert
func callerInsert(caller *Caller) interface{} {
	caller.validateArgNums(2)
	mod := Select(caller.Class)
	var colums = []string{}
	colums, ok := caller.Args[0].([]string)
	if !ok {
		anyColums, ok := caller.Args[0].([]interface{})
		if !ok {
			exception.New("第1个查询参数错误 %v", 400, caller.Args[0]).Throw()
		}
		for _, col := range anyColums {
			colums = append(colums, string(str.Of(col)))
		}
	}

	var rows = [][]interface{}{}
	rows, ok = caller.Args[1].([][]interface{})
	if !ok {
		anyRows, ok := caller.Args[1].([]interface{})
		if !ok {
			exception.New("第2个查询参数错误 %v", 400, caller.Args[1]).Throw()
		}
		for _, anyRow := range anyRows {

			row, ok := anyRow.([]interface{})
			if !ok {
				exception.New("第2个查询参数错误 %v", 400, caller.Args[1]).Throw()
			}
			rows = append(rows, row)
		}
	}

	mod.MustInsert(colums, rows)
	return nil
}