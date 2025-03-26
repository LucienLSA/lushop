package global

import (
	"lushopapi/user_web/config"

	ut "github.com/go-playground/universal-translator"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
)

// var ServerConfig = new(config.ServerConfig)
