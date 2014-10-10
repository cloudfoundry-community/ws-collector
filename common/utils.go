package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func ParseInt(s string, d int) int {
	if len(s) < 1 {
		return d
	}
	//strconv.Btoi64
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		log.Fatalf("unable to parse int from %s: %v", s, err)
		return d
	}
	return int(v)
}

func GetEnvVarAsString(k, d string) string {
	if len(k) < 1 {
		return d
	}
	s := os.Getenv(k)
	if len(s) < 1 {
		return d
	}
	return s
}

func GetEnvVarAsInt(k string, d int) int {
	s := GetEnvVarAsString(k, "")
	if len(s) < 1 {
		return d
	}
	v, err := strconv.ParseInt(s, d, 8)
	if err != nil {
		log.Fatalf("unable to parse int from %s: %v", k, err)
		return d
	}
	return int(v)
}

func GetEnvVarAsBool(k string, d bool) bool {
	s := GetEnvVarAsString(k, "")
	if len(s) < 1 {
		return d
	}
	v, err := strconv.ParseBool(k)
	if err != nil {
		log.Fatalf("unable to parse bool from %s: %v", k, err)
		return d
	}
	return v
}

func Printout(o interface{}) {
	objStr, err := ToString(o)
	if err != nil {
		log.Printf("unable to marshal: %v", err.Error())
	}
	log.Print(fmt.Sprintln(string(objStr)))
}

func ToString(o interface{}) (string, error) {
	objStr, err := json.Marshal(o)
	if err != nil {
		log.Printf("unable to marshal: %v", o)
		log.Panicln(err)
		return "", err
	}
	return fmt.Sprintln(string(objStr)), nil
}

func GetNowInUtc() time.Time {
	return time.Now().UTC()
}

func GetTime(f string) string {
	if len(f) < 1 {
		f = time.RFC850
	}
	return fmt.Sprintln(GetNowInUtc().Format(f))
}
