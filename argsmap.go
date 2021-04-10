package argsmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetCommandLineArgMap(fileName string, args []string) (map[string]string, error) {
	var argHelpMap = make(map[string]oneArg)

	b := loadFile(fileName)
	jserr := json.Unmarshal(b, &argHelpMap)
	if jserr != nil {
		fmt.Printf("解析文件%s出错！\n", fileName)
	}

	var userInputArgMap = make(map[string]string)

	if len(os.Args) > 1 {
		for i := 1; i < len(args); i++ {
			flag := strings.TrimSpace(args[i])
			if usage, ok := argHelpMap[flag]; ok {
				if usage.MustHaveValue {
					v, err := getFlagValueFromArgs(usage, i, args)
					if err != nil {
						showError(usage.ArgValueErrorMsg + ",用户输入值为:'" + v + "', 期望值为:" + usage.ValueExpect)
					} else {
						fmt.Println("指定成功:", flag, v)
						userInputArgMap[flag] = v
						i++
					}
				} else {
					fmt.Println("指定成功:", flag)
					userInputArgMap[flag] = "1"
				}
			} else {
				log.Println("argsmap:", "未知参数:" + flag)
				return userInputArgMap, errors.New("未知参数:" + flag)
			}
		}
	}
	fmt.Println("------------命令行配置------------------")
	fmt.Println(userInputArgMap)
	return userInputArgMap, nil
}

type oneArg struct {
	ArgFlag      string `json:"flag"`
	ValuePattern string `json:"pattern"`
	/*是否必须有值，比如-h显示帮助就不需要，而-p 8888指定端口则必须指定，当必须指定的时候该值为true*/
	ArgValue         string `json:"value"`
	ValueExpect      string `json:"expect"`
	ArgUsage         string `json:"usage"`
	ArgValueErrorMsg string `json:"err"`
	MustHaveValue    bool   `json:"must_have_value"`
}

/**
为struct添加默认值，该方法会被自动调用
*/
func (o *oneArg) UnmarshalJSON(b []byte) error {
	type xOneArg oneArg
	/*是否必须有值，比如-h显示帮助就不需要，而-p 8888指定端口则必须指定，当必须指定的时候该值为true, 不指定则为false*/
	xo := &xOneArg{MustHaveValue: true}
	if err := json.Unmarshal(b, xo); err != nil {
		return err
	}
	*o = oneArg(*xo)
	return nil
}

func loadFile(filePath string) []byte {
	body, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("文件读取出错！" + filePath)
		return []byte{}
	}
	return body
}
func showError(msg string) {
	fmt.Println(msg)
}


func getFlagValueFromArgs(usage oneArg, i int, args []string) (string, error) {
	if i >= (len(args) - 1) {
		return "", errors.New(usage.ArgValueErrorMsg)
	} else {
		return args[i+1], nil
	}
}
