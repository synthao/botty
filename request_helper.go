package botty

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

func addEntitiesToRequest(v url.Values, name, text string, entities []MessageEntity) error {
	if len(entities) > 0 {
		for i := range entities {
			if entities[i].Offset == 0 {
				entities[i].Length = utf8.RuneCountInString(strings.TrimSpace(text))
			}
		}

		serializedEntities, err := json.Marshal(entities)
		if err != nil {
			return fmt.Errorf("can't sent message, %w", err)
		}

		v.Add(name, string(serializedEntities))
	}

	return nil
}

func addEntitiesToForm(f MultipartForm, name, text string, entities []MessageEntity) error {
	if len(entities) > 0 {
		for i := range entities {
			if entities[i].Offset == 0 {
				entities[i].Length = utf8.RuneCountInString(strings.TrimSpace(text))
			}
		}

		serializedEntities, err := json.Marshal(entities)
		if err != nil {
			return fmt.Errorf("can't sent message, %w", err)
		}

		if err := f.AddField(name, string(serializedEntities)); err != nil {
			return err
		}
	}

	return nil
}
