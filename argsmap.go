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

type commandLineObj struct {
	GetCommandLineMap map[string]string
	ShowHelp          func()
}

func NewCommandLineObj(fileName string, args []string) (*commandLineObj, error) {
	var argHelpMap = make(map[string]OneArg)
	b := loadFile(fileName)
	jserr := json.Unmarshal(b, &argHelpMap)
	if jserr != nil {
		fmt.Printf("An error occured when parse file':%s\n'", fileName)
	}
	m, err := GetCommandLineArgMap(argHelpMap, args)
	if err != nil {
		log.Fatal("An error occurred while parsing")
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

func GetCommandLineArgMap(argHelpMap map[string]OneArg, args []string) (map[string]string, error) {
	//var argHelpMap = make(map[string]OneArg)
	var userInputArgMap = make(map[string]string)
	if len(os.Args) > 1 {
		for i := 1; i < len(args); i++ {
			flag := strings.TrimSpace(args[i])
			if usage, ok := argHelpMap[flag]; ok {
				if usage.MustHaveValue {
					v, err := getFlagValueFromArgs(usage, i, args)
					if err != nil {
						showError(usage.ArgValueErrorMsg + ",User Input:'" + v + "', Expect for:" + usage.ValueExpect)
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
		showError("Failed to load file:" + filePath)
		return []byte{}
	}
	return body
}
func showError(msg string) {
	log.Fatal("argsmap", msg)
}

func getFlagValueFromArgs(usage OneArg, i int, args []string) (string, error) {
	if i >= (len(args) - 1) {
		return "", errors.New(usage.ArgValueErrorMsg)
	} else {
		return args[i+1], nil
	}
}
