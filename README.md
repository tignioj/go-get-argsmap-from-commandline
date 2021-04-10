# 封装命令行参数
以前我们需要手动去提取判断用户的命令行参数，比如 `./server.exe -p 8090 -r ./` 这里有2个参数，2个值，使用本项目，经过自定义json配置封装到map中
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
编写好后，只需要简单的调用两个方法就可以获得系统命令行的格式化输出和参数map
上述配置输出如下
```
G:\goProject\go-get-argsmap-from-commandline\yourmode>go run .
2021/04/11 04:24:42 argsmap ------------Command line configuration------------------
2021/04/11 04:24:42 argsmap map[]
Usage:
|------------|----------------------------|----------------|------------------
| flag       | usage                      | expect         | default
|------------|----------------------------|----------------|------------------
| -flag1     | Teach you how to use flag1 | user expect v1 | default value v1
| -flag2     | Teach you how to use flag1 |                |
|------------|----------------------------|----------------|------------------
------------Get Map--------------
&{map[] 0xf36820}

```

例如这个配置
```
{
  "-h": {
    "usage": "show help",
    "must_have_value": false
  },
  "-p": {
    "value": "8080",
    "usage": "server port",
    "pattern": "^[0-9]+$",
    "expect": "pure number",
    "err": "invalid port"
  },
  "-r": {
    "value": "./",
    "usage": "web root",
    "err": "invalid web root"
  },
  "-a": {
    "value": "0.0.0.0",
    "usage": "listen address",
    "pattern": "(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}",
    "expect": "ipv4 address, format like 0.0.0.0",
    "err": "invalid address!"
  },
  "-c": {
    "usage": "path to server configuration",
    "err": "invalid config path"
  }
}
```

格式化输出如下

```

G:\goProject\go-get-argsmap-from-commandline\yourmode>go run .
2021/04/11 04:28:47 argsmap ------------Command line configuration------------------
2021/04/11 04:28:47 argsmap map[]
Usage:
|------------|------------------------------|-----------------------------------|------------
| flag       | usage                        | expect                            | default
|------------|------------------------------|-----------------------------------|------------
| -p         | server port                  | pure number                       | 8080
| -r         | web root                     |                                   | ./
| -a         | listen address               | ipv4 address, format like 0.0.0.0 | 0.0.0.0
| -c         | path to server configuration |                                   |
| -h         | show help                    |                                   |
|------------|------------------------------|-----------------------------------|------------

```

# 用法
 现有的DEMO
 - 简单的HTTP服务器 https://github.com/tignioj/simple-go-httpserver
 

## 方式一：使用`go get`
### 1. 初始化你的项目

 你的项目结构如下



        |-yourmodule
        |   -youcode.go
        |



```
cd yourmodule
go mod init "mymodule.com/mymodule"
```
 此时你的目录下会多一个文件`go.mod`


        |-yourmodule
        |   -youcode.go
        |   -go.mod
        |
        
        

其中go.mod如下
```
module mymodule.com/mymodul

go 1.16

```

### 2. 在执行go get获取本项目
```
go get github.com/tignioj/go-get-argsmap-from-commandline.git
```

此时你的`go.mod`
```
module mymodule.com/mymodul

go 1.16

require github.com/tignioj/go-get-argsmap-from-commandline v1.0.1-0.20210410193735-97119a2e5a7c // indirect

```
  
### 3. 在自己的模块中调用本项目的方法

#### 1. 导包

`import "github.com/tignioj/go-get-argsmap-from-commandline"`

#### 2. 调用
##### 构造方式1， 通过文件
```go
argMap, err := argsmap.NewCommandLineObj("help.json", os.Args)
```

如下
```go
package main

import (
	"fmt"
	"github.com/tignioj/go-get-argsmap-from-commandline"
	"log"
	"os"
)

func main() {
	argMap, err := argsmap.NewCommandLineObj("help.json", os.Args)
	if err != nil {
		log.Fatal(err)
	}
	m := argMap.GetCommandLineMap

	for k, v := range m {
		fmt.Println(k, v)
	}
	argMap.ShowHelp()

	fmt.Println("------------Get Map--------------")
	fmt.Println(argMap)
}
```

##### 构造方式二：通过json字符串

```go
argMap, err := argsmap.NewCommandLineObjByJSON(helpJSON, os.Args)
```
如下
```go
package main

import (
	"fmt"
	"github.com/tignioj/go-get-argsmap-from-commandline"
	"log"
	"os"
)

var helpJSON = `
{
  "-h": {
    "usage": "show help",
    "must_have_value": false
  },
  "-p": {
    "value": "8080",
    "usage": "server port",
    "pattern": "^[0-9]+$",
    "expect": "pure number",
    "err": "invalid port"
  },
  "-r": {
    "value": "./",
    "usage": "web root",
    "err": "invalid web root"
  },
  "-c": {
    "usage": "path to server configuration",
    "err": "invalid config path",
    "value": "server-config.json"
  }
}
`

func main() {
	argMap, err := argsmap.NewCommandLineObjByJSON(helpJSON, os.Args)
	if err != nil {
		log.Println("yourmode", err)
	}
	m := argMap.GetCommandLineMap

	for k, v := range m {
		fmt.Println(k, v)
	}
	argMap.ShowHelp()

	fmt.Println("------------Get Map--------------")
	fmt.Println(argMap)
}
```


#### 编写帮助文档，注意到这里的`help.json`，我们需要自己编写以便于生成帮助文档

比如

```json
{
  "-h": { 
    "usage": "show help",
    "must_have_value": false
  },
  "-p": {
    "value": "8080",
    "usage": "server port",
    "pattern": "^[0-9]+$",
    "expect": "pure number",
    "err": "invalid port"
  },
  "-r": {
    "value": "./",
    "usage": "web root",
    "err": "invalid web root"
  },
  "-a": {
    "value": "0.0.0.0",
    "usage": "listen address",
    "pattern": "(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}",
    "expect": "ipv4 address, format like 0.0.0.0",
    "err": "invalid address!"
  },
  "-c": {
    "usage": "path to server configuration",
    "err": "invalid config path"
  }
}
```

看起来就是这样的
```

G:\goProject\go-get-argsmap-from-commandline\yourmode>go run .
2021/04/11 03:40:02 argsmap ------------Command line configuration------------------
2021/04/11 03:40:02 argsmap map[]
Usage:
|------------|------------------------------|-----------------------------------|------------
| flag       | usage                        | expect                            | default
|------------|------------------------------|-----------------------------------|------------
| -h         | show help                    |                                   |
| -p         | server port                  | pure number                       | 8080
| -r         | web root                     |                                   | ./
| -a         | listen address               | ipv4 address, format like 0.0.0.0 | 0.0.0.0
| -c         | path to server configuration |                                   |
|------------|------------------------------|-----------------------------------|------------

```

## 方式二：本地用法：在其它模块上使用本模块
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

### 5. 直接你的模块代码`youcode.go`中调用同上

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

