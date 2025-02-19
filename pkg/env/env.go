package env

import (
	"fmt"
	"github.com/DenisBarabanshchikov/subscription/pkg/util"
	"net/url"
	"os"
	"strconv"
	"time"
)

func OptionalString(key string) string {
	return os.Getenv(key)
}

func OptionalStringPtr(key string) *string {
	v := OptionalString(key)
	if v == "" {
		return nil
	}
	return &v
}

func RequiredString(key string) string {
	v := OptionalString(key)
	if v == "" {
		panic(fmt.Sprintf("missing required env variable '%s'", key))
	}
	return v
}

func OptionalInt(key string) int {
	v := OptionalString(key)
	i, _ := strconv.Atoi(v)
	return i
}

func RequiredInt(key string) int {
	v := RequiredString(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("invalid int of env variable '%s': %s", key, v))
	}
	return i
}

func OptionalBool(key string) bool {
	v := OptionalString(key)
	b, _ := strconv.ParseBool(v)
	return b
}

func RequiredBool(key string) bool {
	v := RequiredString(key)
	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(fmt.Sprintf("invalid bool of env variable '%s': %s", key, v))
	}
	return b
}

func OptionalTime(key string) *time.Time {
	v := OptionalString(key)
	if v == "" {
		return nil
	}
	t, err := util.ParseTime(v)
	if err != nil {
		panic(fmt.Sprintf("invalid time of env variable '%s': %s", key, v))
	}
	return &t
}

func RequiredTime(key string) *time.Time {
	v := RequiredString(key)
	t, err := util.ParseTime(v)
	if err != nil {
		panic(fmt.Sprintf("invalid time of env variable '%s': %s", key, v))
	}
	return &t
}

func OptionalDuration(key string) time.Duration {
	v := OptionalString(key)
	d, _ := time.ParseDuration(v)
	return d
}

func RequiredDuration(key string) time.Duration {
	v := RequiredString(key)
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(fmt.Sprintf("invalid duration of env variable '%s': %s", key, v))
	}
	return d
}

func OptionalUrl(key string) *url.URL {
	v := OptionalString(key)
	if v == "" {
		return nil
	}
	u, err := url.Parse(v)
	if err != nil {
		panic(fmt.Sprintf("invalid URL of env variable '%s': %s", key, v))
	}
	return u
}

func RequiredUrl(key string) *url.URL {
	v := RequiredString(key)
	u, err := url.Parse(v)
	if err != nil {
		panic(fmt.Sprintf("invalid URL of env variable '%s': %s", key, v))
	}
	return u
}
