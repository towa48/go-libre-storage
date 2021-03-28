package assets

import (
	"encoding/json"
	"io/ioutil"
	"path"
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
	ScriptChunks      []string
	StyleChunks       []string
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

	manifest.ScriptChunks = filter(files, isScriptChunk)
	manifest.StyleChunks = filter(files, isStyleChunk)
	return *manifest, nil
}

func GetAssetContent(p string) (content string, err error) {
	if cachedContent == nil {
		cachedContent = make(map[string]string, 5)
	}

	if val, ok := cachedContent[p]; ok {
		return val, nil
	}

	manifestPath := config.Get().AssetManifestPath
	basePath := path.Dir(manifestPath)
	fullPath := path.Join(basePath, p)

	bytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	// remove source maps
	c := string(bytes)
	i := strings.Index(c, "//# sourceMappingURL")
	cachedContent[p] = c[0:i]

	return cachedContent[p], nil
}

func isScriptChunk(p string) bool {
	return strings.HasSuffix(p, ".chunk.js")
}

func isStyleChunk(p string) bool {
	return strings.HasSuffix(p, ".chunk.css")
}

func filter(values map[string]interface{}, test func(string) bool) (ret []string) {
	for k, v := range values {
		if test(k) {
			ret = append(ret, v.(string))
		}
	}
	return
}
