package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestMakesMarkdown(t *testing.T) {
	review1 := Review{
		Title:   "Love it!",
		Content: "Great app.",
		Date:    "2022-10-26T08:44:41-07:00",
		Rating:  strconv.FormatInt(int64(rand.Intn(5)), 10),
	}

	review2 := Review{
		Title:   "Meh",
		Content: "I've seen better.",
		Date:    "2022-10-24T12:47:16-07:00",
		Rating:  strconv.FormatInt(int64(rand.Intn(3)), 10),
	}

	digest := Digest{
		Title:    "Today's digest",
		Subtitle: "Juicy reviews inside!",
		Reviews:  []Review{review1, review2},
	}

	expectedFormat := `# Today's digest
Juicy reviews inside!
## Love it!
**%s star(s)** — _Wed Oct 26 08:44:41 2022_

Great app.
## Meh
**%s star(s)** — _Mon Oct 24 12:47:16 2022_

I've seen better.
`

	expected := fmt.Sprintf(expectedFormat, review1.Rating, review2.Rating)
	actual := digest.ToMarkdown()

	if expected != actual {
		t.Fail()
	}
}
