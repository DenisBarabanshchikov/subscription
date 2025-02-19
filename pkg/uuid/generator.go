package uuid

import "github.com/segmentio/ksuid"

func GenerateUUID() string {
	uuid := ksuid.New()
	return uuid.String()
}
