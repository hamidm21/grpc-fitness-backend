package assert

import (
	"gitlab.com/mefit/mefit-server/utils/log"
)

func Nil(in interface{}) {
	if in != nil {
		if err, ok := in.(error); ok {
			log.Logger().Panic(err)
		}
	}
}

func True(in bool, err error) {
	if !in {
		log.Logger().Panic(err)
	}
}
