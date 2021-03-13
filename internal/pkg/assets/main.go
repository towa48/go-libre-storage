package assets

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/towa48/go-libre-storage/internal/pkg/config"
)

var manifest *Manifest
var cachedContent map[string]string

type Manifest struct {
	WelcomeStyleUrl   string
	WelcomeRuntimeUrl string
	WelcomeScriptUrl  string
	MainStyleUrl      string
	MainRuntimeUrl    string
	MainScriptUrl     string
	ChunkScriptUrl    string
}

func GetAssetsManifest() (value Manifest, err error) {
	if manifest != nil {
		return *manifest, nil
	}

	var m Manifest
	configPath := config.Get().AssetManifestPath
	source, err := ioutil.ReadFile(configPath)
	if err != nil {
		return m, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(source), &result)

	files := result["files"].(map[string]interface{})
	manifest = &Manifest{
		WelcomeStyleUrl:   files["welcome.css"].(string),
		WelcomeRuntimeUrl: files["runtime-welcome.js"].(string),
		WelcomeScriptUrl:  files["welcome.js"].(string),
		MainStyleUrl:      files["main.css"].(string),
		MainRuntimeUrl:    files["runtime-main.js"].(string),
		MainScriptUrl:     files["main.js"].(string),
	}

	manifest.ChunkScriptUrl = filter(files, isChunk)[0]
	return *manifest, nil
}

func GetAssetContent(path string) (content string, err error) {
	if val, ok := cachedContent[path]; ok {
		return val, nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	cachedContent[path] = string(bytes)
	return cachedContent[path], nil
}

func isChunk(path string) bool {
	return strings.HasSuffix(path, ".chunk.js")
}

func filter(values map[string]interface{}, test func(string) bool) (ret []string) {
	for k, v := range values {
		if test(k) {
			ret = append(ret, v.(string))
		}
	}
	return
}
