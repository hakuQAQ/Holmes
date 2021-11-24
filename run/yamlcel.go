package run

import (
	"Holmes/utils"
	"encoding/json"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/spf13/viper"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"log"
	"os"
	"strings"
)

type YamlRule struct {
	Infos   string
	Matches []string
}

type Info struct {
	Company 	string `json:"company"`
	Product     string `json:"product"`
	Softhard    string `json:"softhard"`
}

func ParseYaml(filename string) YamlRule{
	tmprule := YamlRule{}
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := viper.New()
	config.AddConfigPath(path + "\\" + filename[:strings.LastIndex(filename, "\\")])  //设置读取的文件路径
	config.SetConfigName(filename[strings.LastIndex(filename, "\\"):]) //设置读取的文件名
	config.SetConfigType("yaml") //设置文件的类型
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		utils.OptionsError("read yamlconfig err", 2)
	}
	//打印文件读取出来的内容:
	mathrules := config.Get("matches").(string)
	slice := utils.LinestoSlice(mathrules)

	for _, line := range slice{
		if line != "" {
			tmprule.Matches = append(tmprule.Matches, line)
		}
	}
	data, err:=json.Marshal(config.Get("info"))
	if err != nil {
		utils.OptionsError("yamlfile have no info", 2)
	}
	var info Info
	err = json.Unmarshal(data, &info)
	tmprule.Infos = info.Product
	return tmprule
}

func LoadYaml(fileName string) ([]YamlRule, error) {
	if !utils.Exists(fileName) {
		utils.OptionsError("rule file not exists!", 2)
	}
	var rules []YamlRule
	if utils.IsFile(fileName) {
		rules = append(rules, ParseYaml(fileName))
	}
	if utils.IsDir(fileName) {
		dirfiles := utils.ReadDir(fileName)
		for _, file := range dirfiles {
			rules = append(rules, ParseYaml(file))
		}
		return rules, nil
	}
	return nil, nil
}


type customLib struct {}

func (e *Engine) ExecuteCel(responinfo Responseinfo, rules []YamlRule) (string, error) {
	env, err := cel.NewEnv(cel.Lib(customLib{}))
	if err != nil {
		log.Fatal(err)
	}
	var product []string

	for _, r := range rules {
		var matches string
		for i, match := range r.Matches{
			if i < len(r.Matches)-1 {
				matches = matches + "(" + match + ") || "
			} else {
				matches = matches + "(" + match + ")"
			}
		}
		ast, iss := env.Compile(matches)
		if iss.Err() != nil {
			log.Fatal(iss.Err())
		}
		prg, err := env.Program(ast)
		if err != nil {
			log.Fatal(err)
		}
		out, _, err := prg.Eval(map[string]interface{}{
			"body":				responinfo.Content,
			"title":			responinfo.Title,
			"header":			responinfo.Headers,
			"server":			responinfo.Server,
			"cert":				responinfo.Cert,
			"banner":			"",
			"protocol":			"",
		})

		if out.(types.Bool) {
			product = append(product, r.Infos)
		}
	}
	return utils.SlicetoSting(product), nil
}

func (customLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Declarations(
			decls.NewVar("body", decls.String),
			decls.NewVar("title", decls.String),
			decls.NewVar("header", decls.String),
			decls.NewVar("server", decls.String),
			decls.NewVar("cert", decls.String),
			decls.NewVar("banner", decls.String),
			decls.NewVar("protocol", decls.String),
			decls.NewFunction("icontains", decls.NewInstanceOverload("icontains_func",
				[]*exprpb.Type{decls.String, decls.String},
				decls.Bool,
			)),
		),
	}
}
func (customLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption {
		cel.Functions(
			&functions.Overload{
				Operator: "icontains_func",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bmatch", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bmatch", lhs.Type())
					}
					if strings.Contains(strings.ToLower(string(v1)), strings.ToLower(string(v2))) {
						return types.Bool(true)
					}
					return types.Bool(false)
				},
			},

		),
	}
}