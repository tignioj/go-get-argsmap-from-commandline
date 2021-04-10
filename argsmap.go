package argsmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type commandLineObj struct {
	GetCommandLineMap map[string]string
	ShowHelp          func()
}

func NewCommandLineObj(fileName string, args []string) (*commandLineObj, error) {
	var argHelpMap = make(map[string]OneArg)
	b := loadFile(fileName)
	jserr := json.Unmarshal(b, &argHelpMap)
	if jserr != nil {
		return nil, errors.New("An error occurred while parsing file:" + fileName)
	}
	m, err := GetCommandLineArgMap(argHelpMap, args)
	if err != nil {
		return nil, errors.New("An error occurred while parsing from map:" + fmt.Sprint(argHelpMap))
	}
	c := commandLineObj{
		GetCommandLineMap: m,
		ShowHelp: func() {
			fmt.Println("Usage:")
			f := "    %-10s%-20s%-20s%-20s\n"
			fmt.Printf(f, "flag", "usage", "expect", "default")
			for k, v := range argHelpMap {
				fmt.Printf(f, k, v.ArgUsage, v.ValueExpect, v.ArgValue)
			}
		},
	}
	return &c, nil
}

/**
get args map from commandline
*/
func GetCommandLineArgMap(argHelpMap map[string]OneArg, args []string) (map[string]string, error) {
	//var argHelpMap = make(map[string]OneArg)
	var userInputArgMap = make(map[string]string)
	if len(os.Args) > 1 {
		for i := 1; i < len(args); i++ {
			// help flag, such as "-p", "-h"
			flag := strings.TrimSpace(args[i])
			// means has a key call `flag` in help.json file.
			if usage, ok := argHelpMap[flag]; ok {
				// must provide value, you can config it in help.json file
				if usage.MustHaveValue {
					v, err := getFlagValueFromArgs(usage, i, args)
					if err != nil {
						//showError(usage.ArgValueErrorMsg + ",User Input:'" + v + "', Expect for:" + usage.ValueExpect)
						showError(err)
					} else {
						log.Println("argsmap", "Binding success:", flag, v)
						userInputArgMap[flag] = v
						i++
					}
				} else {
					log.Println("argsmap", "Bidding success:", flag)
					userInputArgMap[flag] = "1"
				}
			} else {
				log.Println("argsmap:", "Unknown param:"+flag)
				return userInputArgMap, errors.New("Unknown param:" + flag)
			}
		}
	}
	log.Println("argsmap", "------------Command line configuration------------------")
	log.Println("argsmap", userInputArgMap)
	return userInputArgMap, nil
}

type OneArg struct {
	/* -p -h 等*/
	ArgFlag      string `json:"flag"`
	/* 正则匹配参数 */
	ValuePattern string `json:"pattern"`
	/* 封装用户输入 */
	ArgValue         string `json:"value"`
	/* 期望输入 */
	ValueExpect      string `json:"expect"`
	/* 显示该flag的用法 */
	ArgUsage         string `json:"usage"`
	/* 当不匹配时候，显示的错误信息 */
	ArgValueErrorMsg string `json:"err"`
	/*是否必须有值，比如-h显示帮助就不需要，而-p 8888指定端口则必须指定，当必须指定的时候该值为true*/
	MustHaveValue    bool   `json:"must_have_value"`
}

/**
为struct添加默认值，该方法会被自动调用
*/
func (o *OneArg) UnmarshalJSON(b []byte) error {
	type xOneArg OneArg
	/*是否必须有值，比如-h显示帮助就不需要，而-p 8888指定端口则必须指定，当必须指定的时候该值为true, 不指定则为false*/
	xo := &xOneArg{MustHaveValue: true}
	if err := json.Unmarshal(b, xo); err != nil {
		return err
	}
	*o = OneArg(*xo)
	return nil
}

func loadFile(filePath string) []byte {
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		showError(errors.New("Failed to load file:" + filePath))
		return []byte{}
	}
	return body
}
func showError(msg error) {
	log.Fatal("argsmap:", msg)
}

func getFlagValueFromArgs(usage OneArg, i int, args []string) (string, error) {
	if i >= (len(args) - 1) {
		return "", errors.New("you have not provide value, info:" +  usage.ArgValueErrorMsg)
	}

	userIn := args[i+1]

	if usage.ValuePattern != "" {
		r := regexp.MustCompile(usage.ValuePattern)
		if !r.MatchString(userIn) {
			return "", errors.New( "your input:" + userIn + "' not match pattern:'" + usage.ValuePattern + "', info:" + usage.ArgValueErrorMsg )
		}
	}

	return args[i+1], nil
}
