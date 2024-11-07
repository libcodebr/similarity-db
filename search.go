package similarity

import (
	"errors"
	"github.com/xrash/smetrics"
	"sort"
	"strings"
	"sync"
)

var (
	ErrTitleIsEmpty  = errors.New("title is empty")
	ErrQueryIsEmpty  = errors.New("query is empty")
	ErrDocumentIsNil = errors.New("document is nil")
	ErrNotFound      = errors.New("not found")
)

type Document struct {
	Title string
	Value any
}

type result struct {
	Value      any
	Similarity float64
}

type DB interface {
	Add(document *Document) error
	Batch(documents []*Document) (err error)
	Search(query string, lenght int) ([]any, error)
	Length() int
}

func New() DB {
	return &db{
		values: map[string]any{},
	}
}

type db struct {
	sync.RWMutex

	values map[string]any // map[title]any
}

func (d *db) Add(document *Document) error {
	if document == nil {
		return ErrDocumentIsNil
	}

	d.Lock()
	defer d.Unlock()
	d.values[document.Title] = document.Value

	return nil
}

func (d *db) Batch(documents []*Document) (err error) {
	if len(documents) == 0 {
		return ErrDocumentIsNil
	}

	for _, v := range documents {
		if err := d.Add(v); err != nil {
			err = errors.Join(err)
		}
	}

	return err
}

// DB searches movies and series
func (d *db) Search(query string, lenght int) ([]any, error) {
	if query == "" {
		return nil, ErrQueryIsEmpty
	}

	d.RLock()
	defer d.RUnlock()

	var results []*result
	for k, v := range d.values {
		if boyerMoore(k, query) != -1 {
			similarity, err := calculateSimilarity(k, query)
			if err != nil {
				continue
			}

			results = append(results, &result{
				Value:      v,
				Similarity: similarity,
			})
		}
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	response := make([]any, lenght)
	for i, v := range results {
		if len(response) > lenght {
			break
		}

		response[i] = v
	}

	return response, nil
}

func (d *db) Length() int {
	return len(d.values)
}

func calculateSimilarity(title string, query string) (float64, error) {
	if title == "" {
		return 0, ErrTitleIsEmpty
	}

	var maxSimilarity float64
	similarity := smetrics.JaroWinkler(query, title, 0.7, 4)
	if similarity > maxSimilarity {
		maxSimilarity = similarity
	}

	return maxSimilarity, nil
}

// boyerMoore is a string search algorithm that searches for a pattern in a text way faster than the regex search
func boyerMoore(text, pattern string) int {
	text = strings.ToLower(text)

	if pattern == "" || strings.TrimSpace(pattern) == "" {
		return 0
	}
	pattern = strings.ToLower(pattern)

	n := len(text)
	m := len(pattern)
	if m == 0 {
		return 0
	}

	bc := make([]int, 256)
	for i := range bc {
		bc[i] = m
	}
	for i := 0; i < m-1; i++ {
		bc[pattern[i]] = m - i - 1
	}

	s := 0
	for s <= n-m {
		j := m - 1
		for j >= 0 && pattern[j] == text[s+j] {
			j--
		}
		if j < 0 {
			return s
		}
		s += bc[text[s+m-1]]
	}
	return -1
}
