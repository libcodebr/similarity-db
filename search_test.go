package similarity

import (
	"testing"
)

func TestAdd(t *testing.T) {
	db := New()

	doc := &Document{
		Title: "Test Movie",
		Value: "Movie Details",
	}

	err := db.Add(doc)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if db.Length() != 1 {
		t.Fatalf("Expected length to be 1, got %d", db.Length())
	}
}

func TestAddNilDocument(t *testing.T) {
	db := New()

	err := db.Add(nil)
	if err == nil {
		t.Fatalf("Expected error when adding nil document, got nil")
	}
}

func TestSearch(t *testing.T) {
	db := New()

	docs := []*Document{
		{Title: "Test Movie 1", Value: "Details of Movie 1"},
		{Title: "Another Movie", Value: "Details of Another Movie"},
		{Title: "Documentary", Value: "Details of Documentary"},
	}

	err := db.Batch(docs)
	if err != nil {
		t.Fatalf("Expected no error on batch add, got %v", err)
	}

	results, err := db.Search("Movie", 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	db := New()

	_, err := db.Search("", 2)
	if err == nil {
		t.Fatalf("Expected error ng with empty query, got nil")
	}
}

func TestBatchAdd(t *testing.T) {
	db := New()

	docs := []*Document{
		{Title: "First Doc", Value: "Value 1"},
		{Title: "Second Doc", Value: "Value 2"},
	}

	err := db.Batch(docs)
	if err != nil {
		t.Fatalf("Expected no error on batch add, got %v", err)
	}

	if db.Length() != 2 {
		t.Fatalf("Expected length to be 2, got %d", db.Length())
	}
}

func TestSearchWithNoResults(t *testing.T) {
	db := New()

	db.Add(&Document{Title: "Only Movie", Value: "Some Details"})

	_, err := db.Search("Nonexistent", 1)
	if err == nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
