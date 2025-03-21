package sourcecode

import (
	"encoding/json"
	"github.com/hdget/utils/convert"
	"os"
)

func saveFile(file string, data any) error {
	var err error
	var content []byte
	switch v := data.(type) {
	case string:
		content = convert.StringToBytes(v)
	case []byte:
		content = v
	default:
		content, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
