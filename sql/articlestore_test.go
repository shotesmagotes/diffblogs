package sql

import (
	"testing"
)

func TestNew(t *testing.T) {
	a := New(
		"Title1",
		"Blog",
		"what_happens_next.html",
		"2017-05-21",
		"2017-08-23",
	)

	if a.id == 0 {
		t.Error("Article was not stored in database.")
	}
}

func TestGet(t *testing.T) {
	a := New(
		"Title2",
		"Blog",
		"what_happens_next2.html",
		"2017-05-22",
		"2017-06-23",
	)

	if a.id == 0 {
		t.Error("Article was not stored in the database.")
	}

	q := map[string]interface{} {
		"title": "Title2",
		"type": "Blog",
	}

	b, _ := Get(q)

	if b.id != a.id {
		t.Error("Retrieved wrong value.")
	}
}

func TestSet(t *testing.T) {
	a := New(
		"Title3",
		"Project",
		"what_happens_next3.html",
		"2017-05-22",
		"2018-06-23",
	)

	if a.id == 0 {
		t.Error("Article was not stored in the database.")
	}

	u := map[string]interface{} {
		"title": "Title3",
		"type": "Blog",
		"filename": "an_update_happens.html",
		"modified": "2017-08-22",
	}

	b, _ := a.Set(u)
	if b.id != a.id {
		t.Error("Retrieved wrong value.")
	}
}