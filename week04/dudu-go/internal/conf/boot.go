package conf

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

/*
@Time : 2021/6/17 下午9:58
@Author : snaker95
@File : boot
@Software: GoLand
*/
func Scan(flagconf string, conf *Config) error {
	yamlFile, err := ioutil.ReadFile(flagconf)
	if err != nil {
		return errors.Wrap(err, "conf: ReadFile fail")
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return errors.Wrap(err, "conf: Unmarshal fail")
	}
	return nil
}
