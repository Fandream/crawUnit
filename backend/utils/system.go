package utils

import (
	"crawunit/entity"
	"encoding/json"
	"github.com/apex/log"
	"io/ioutil"
	"runtime/debug"
)

func GetLangList() []entity.Lang {
	list := []entity.Lang{
		{
			Name:              "Python",
			ExecutableName:    "python",
			ExecutablePaths:   []string{"/usr/bin/python", "/usr/local/bin/python"},
			DepExecutablePath: "/usr/local/bin/pip",
			LockPath:          "/tmp/install-python.lock",
			DepFileName:       "requirements.txt",
			InstallDepArgs:    "install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt",
		},
	}
	return list
}

// 获取语言列表
func GetLangListPlain() []entity.Lang {
	list := GetLangList()
	return list
}

// 根据语言名获取语言实例，不包含状态
func GetLangFromLangNamePlain(name string) entity.Lang {
	langList := GetLangListPlain()
	for _, lang := range langList {
		if lang.ExecutableName == name {
			return lang
		}
	}
	return entity.Lang{}
}

func GetPackageJsonDeps(filepath string) (deps []string, err error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Errorf("get package.json deps error: " + err.Error())
		debug.PrintStack()
		return deps, err
	}
	var packageJson entity.PackageJson
	if err := json.Unmarshal(data, &packageJson); err != nil {
		log.Errorf("get package.json deps error: " + err.Error())
		debug.PrintStack()
		return deps, err
	}

	for d, v := range packageJson.Dependencies {
		deps = append(deps, d+"@"+v)
	}

	return deps, nil
}
