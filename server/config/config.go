package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

const (
	CNROOT = "/home/troper/Projects/Thesis/Source" // TODO: set in conf
	SERVER = CNROOT + "/server"
	CLIENT = CNROOT + "/client"
)

// FUTURE: Set via conf
// FUTURE: If defined in env use that

type dockerConf struct {
	IsDocker       bool
	ContainerName string `json:",omitempty"`
}
type cnConfConf struct {
	Path      string
	Overrides string `json:",omitempty"`
}
type cnConfIpfsJs struct {
	Conf       cnConfConf
	Docker     *dockerConf `json:",omitempty"`
	EnableLink bool        `json:",omitempty"`
}
type cnConfIpfs struct {
	Conf          cnConfConf
	Js            *cnConfIpfsJs `json:",omitempty"`
	Docker        *dockerConf `json:",omitempty"`
	InferAnnounce bool `json:",omitempty"`
}
type AppConf struct {
	Enable      bool
	AppSettings map[string]interface{}
}
type CnConf struct {
	Ipfs cnConfIpfs
	App  struct {
		Apps      map[string]AppConf
		JsAppPath string
	}
	Storage struct {
		Path string
	}
	Requests string
	Pid      string
	Root     struct {
		Tick   int
		Retain int
	}
	path string
}

func WriteUnstructured(confPath string, output map[string]interface{}) error {
	b, err := json.Marshal(output)
	if err != nil {
		return err
	}
	return os.WriteFile(confPath, b, 0600)
}

func LoadUnstructured(confPath string, output *map[string]interface{}) error {
	confFile, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(confFile, output)
	return err
}

func (conf CnConf) Write() error {
	b, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	return os.WriteFile(conf.path, b, 0600)
}

func Load(confPath string) (CnConf, error) {
	confFile, err := os.ReadFile(confPath)
	var res CnConf
	if err != nil {
		return res, err
	}
	confFile = bytes.ReplaceAll(confFile, []byte("$CNROOT"), []byte(CNROOT))
	confFile = bytes.ReplaceAll(confFile, []byte("$SERVER"), []byte(SERVER))
	confFile = bytes.ReplaceAll(confFile, []byte("$CLIENT"), []byte(CLIENT))
	err = json.Unmarshal(confFile, &res)
	if err != nil {
		return res, err
	}
	res.path = confPath
	return res, nil
}

func ReplaceOrAdd(newVals, toReplace map[string]interface{}) (map[string]interface{}, error) {
	// We take the vals in newVals and either add or replace them in toReplace
	// If the value is a JSON Object, rather than replace it wholesale we replace
	// or add by it's keys rather than the whole object
	for key, newElement := range newVals {
		oldElement, exists := toReplace[key]
		if exists {
			toReplaceSubVals, isJSONObject := oldElement.(map[string]interface{})
			if isJSONObject {
				newSubVals, isJSONObject := newElement.(map[string]interface{})
				if !isJSONObject {
					return nil, fmt.Errorf("Attempting to replace JSON object %v with a non-JSON object %v", oldElement, newElement)
				}
				updatedSubVals, err := ReplaceOrAdd(newSubVals, toReplaceSubVals)
				if err != nil {
					return nil, err
				}
				// Succesfully updated the subJSON? Replace
				toReplace[key] = updatedSubVals
			} else {
				// Not JSON? Replace
				toReplace[key] = newElement
			}
		} else {
			// Doesn't exist? Add it
			toReplace[key] = newElement
		}
	}
	return toReplace, nil
}
