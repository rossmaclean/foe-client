package properties

import (
	"fmt"
	"github.com/magiconair/properties"
)

type Properties struct {
	RemoteUrl   string
	RunInterval string
}

var props Properties

func LoadProperties(filePath string) {
	p := properties.MustLoadFile(filePath, properties.UTF8)
	props = Properties{
		RemoteUrl:   p.MustGet("remoteUrl"),
		RunInterval: p.MustGet("runInterval"),
	}
	fmt.Printf("Using %s properties\n", filePath)
}

func GetProperties() Properties {
	return props
}
