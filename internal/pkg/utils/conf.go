package utils

import (
	"github.com/zeromicro/go-zero/core/conf"
)

const (
	ValueConfFile = "etc/deutsch.yaml"
)

func LoadConfig(out any) error {
	return UnmarshalConfig(ValueConfFile, out)
}

func UnmarshalConfig(path string, out any) error {
	conf.MustLoad(path, out)

	err := SetDefaults(out)
	if err != nil {
		return err
	}
	err = GetValidator().Struct(out)
	if err != nil {
		return err
	}
	return nil
}
