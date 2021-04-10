# 封装命令行参数
以前我们需要手动去提取判断用户的命令行参数，比如 `./server -p 8090 -r --noauth./` 这里有三个参数，2个值，使用本项目，经过自定义json配置封装到map中
json配置格式如下
```json
{
    "-flag1": {
        "usage":"Teach you how to use flag1",
        "value":"default value v1",
        "expect": "user expect v1",
        "err":"show error message when not match pattern1"
    },
    "-flag2": {
        "usage":"Teach you how to use flag1",
        "must_have_value": false
    }
}
```
例如
```
{
  "-h": {
    "usage": "显示帮助",
    "must_have_value": false
  },
  "-p": {
    "value": "8080",
    "usage": "展示端口",
    "pattern": "\\d+",
    "expect": "纯数字",
    "err": "端口指定错误！"
  },
  "-r": {
    "value": "./",
    "usage": "服务器根目录",
    "expect": "正确的地址格式",
    "err": "该地址有误！"
  },
  "-c": {
    "usage": "配置文件目录",
    "err": "配置地址有误"
  }
}
```

# 用法
## 在其它模块上使用本模块

### 1. 先下载本模块到本地，项目结构如下

    

    |-argsmap
    |   - argsmap.go
    |   - go.mod
    |-youmode
    |   -youcode.go
    |   -go.mod
    |



### 2. 进入你的模块
```
cd youmode
```


### 3. go模块指向本地
```shell
go mod edit -replace="tignioj.io/argsmap"="../argsmap"
go mod tidy
```
输出
PS G:\goProject\ArgsMapWithHelp\yourmode> go mod tidy
go: found tignioj.io/argsmap in tignioj.io/argsmap v0.0.0-00010101000000-000000000000


### 4. 此时你的项目`go.mod`类似如下
```
module tignioj.io/gohttpserver

go 1.16

replace tignioj.io/argsmap => ../argsmap

require tignioj.io/argsmap v0.0.0-00010101000000-000000000000
```

### 5. 直接你的模块代码`youcode.go`中调用

```go
package main

import (
    "log"
	"tignioj.io/argsmap"
)
func main() {
	argMap, argerr:= argsmap.GetCommandLineArgMap("help.json", os.Args)
	if argerr != nil{
		log.Fatal(argerr)
	}
}
```

### 6. 编辑命令行配置文件

添加一个`help.json`



    |-argsmap
    |   - argsmap.go
    |   - go.mod
    |-youmode
    |   -youcode.go
    |   -go.mod
    |   -help.json
    |



- help.json文件举例

```json
{
  "-h": {
    "usage": "显示帮助",
    "must_have_value": false
  },
  "-p": {
    "value": "8080",
    "usage": "展示端口",
    "pattern": "\\d+",
    "expect": "纯数字",
    "err": "端口指定错误！"
  },
  "-r": {
    "value": "./",
    "usage": "服务器根目录",
    "expect": "正确的地址格式",
    "err": "该地址有误！"
  },
  "-c": {
    "usage": "配置文件目录",
    "err": "配置地址有误"
  }
}
```

### 7. 运行
```
cd ../yourmode/
go run youmode.go -p 8080 -r ./
```

