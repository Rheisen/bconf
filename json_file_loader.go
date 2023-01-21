package bconf

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// type JSONMarshal func(v interface{}) ([]byte, error)

type JSONUnmarshal func(data []byte, v interface{}) error

func NewJSONFileLoader() *JSONFileLoader {
	return NewJSONFileLoaderWithAttributes(nil)
}

func NewJSONFileLoaderWithAttributes(decoder JSONUnmarshal, filePaths ...string) *JSONFileLoader {
	return &JSONFileLoader{
		Decoder:   decoder,
		FilePaths: filePaths,
	}
}

type JSONFileLoader struct {
	Decoder   JSONUnmarshal
	FilePaths []string
	// Encoder   JSONMarshal
}

func (l *JSONFileLoader) Clone() *JSONFileLoader {
	clone := *l

	clone.FilePaths = make([]string, len(l.FilePaths))
	copy(clone.FilePaths, l.FilePaths)

	return &clone
}

func (l *JSONFileLoader) CloneLoader() Loader {
	return l.Clone()
}

func (l *JSONFileLoader) Name() string {
	return "bconf_jsonfile"
}

func (l *JSONFileLoader) Get(fieldSetKey, fieldKey string) (string, bool) {
	maps := l.fileMaps()

	if len(maps) < 1 {
		return "", false
	}

	return l.findValueInMaps(fieldSetKey, fieldKey, &maps)
}

func (l *JSONFileLoader) GetMap(fieldSetKey string, fieldKeys []string) map[string]string {
	values := map[string]string{}

	maps := l.fileMaps()

	if len(maps) < 1 {
		return values
	}

	for _, fieldKey := range fieldKeys {
		val, found := l.findValueInMaps(fieldSetKey, fieldKey, &maps)
		if found {
			values[fieldKey] = val
		}
	}

	return values
}

func (l *JSONFileLoader) HelpString(fieldSetKey, fieldKey string) string {
	return fmt.Sprintf("JSON attribute: %s.%s", fieldSetKey, fieldKey)
}

func (l *JSONFileLoader) findValueInMaps(fieldSetKey, fieldKey string, maps *[]map[string]any) (string, bool) {
	if maps == nil {
		return "", false
	}

	for _, fileMap := range *maps {
		fieldSetAny, found := fileMap[fieldSetKey]
		if !found {
			continue
		}

		fieldSetMap, ok := fieldSetAny.(map[string]any)
		if !ok {
			continue
		}

		value, ok := fieldSetMap[fieldKey]
		if !ok {
			continue
		}

		bytes, _ := json.Marshal(value)
		valueString := string(bytes)

		if strings.HasPrefix(valueString, "[") && strings.HasSuffix(valueString, "]") {
			valueString = valueString[1 : len(valueString)-1]
			valueStringSlice := strings.Split(valueString, ",")
			for index, val := range valueStringSlice {
				valueStringSlice[index] = strings.Trim(val, "\"")
			}
			valueString = strings.Join(valueStringSlice, ",")
		} else {
			valueString = strings.Trim(valueString, "\"")
		}

		return valueString, true
	}

	return "", false
}

func (l *JSONFileLoader) fileMaps() []map[string]any {
	fileMaps := []map[string]any{}

	for _, path := range l.FilePaths {
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		fileMap := map[string]any{}
		if l.Decoder != nil {
			if err := l.Decoder(fileBytes, &fileMap); err != nil {
				continue
			}
		} else {
			if err := json.Unmarshal(fileBytes, &fileMap); err != nil {
				continue
			}
		}

		fileMaps = append(fileMaps, fileMap)
	}

	return fileMaps
}
