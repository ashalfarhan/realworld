package conduit

import (
	"fmt"
	"log"

	"github.com/gosimple/slug"
	"github.com/matoous/go-nanoid/v2"
)

const (
	defaultSlugId = 8
)

func CreateSlug(title string) string {
	id, err := gonanoid.New(defaultSlugId)
	if err != nil {
		log.Printf("[Slug Generation Error] Cannot create slug for %s, %s\n", title, err.Error())
	}

	s := slug.Make(title)

	return fmt.Sprintf("%s-%s", s, id)
}
