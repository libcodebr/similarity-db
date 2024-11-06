# Search Package

The `search` package is a Go library designed for storing and searching documents efficiently. It offers a straightforward API to add, batch insert, and search for documents based on their titles using optimized string search algorithms and similarity scoring.

## Features

- **Add individual documents** with a title and value.
- **Batch insertion** for multiple documents.
- **Search functionality** using a combination of the Boyer-Moore string search algorithm and Jaro-Winkler similarity for relevant matches.
- **Thread-safe operations** ensured through the use of `sync.RWMutex`.

## Installation

Ensure you have Go installed and set up in your environment. Then, add this package to your project:

```bash
go get github.com/yourusername/yourproject/search
```

Import it in your Go file:

```go
import search "github.com/libcodebr/memory-search"
```

## Usage

### Creating a New Search Database

```go
db := search.New()
```

### Adding a Document

```go
doc := &search.Document{
    Title: "Example Title",
    Value: "Example Content",
}
err := db.Add(doc)
if err != nil {
    log.Fatalf("Error adding document: %v", err)
}
```

### Batch Adding Documents

```go
docs := []*search.Document{
    {Title: "Doc 1", Value: "Content 1"},
    {Title: "Doc 2", Value: "Content 2"},
}
err := db.Batch(docs)
if err != nil {
    log.Fatalf("Error adding documents in batch: %v", err)
}
```

### Searching for Documents

```go
results, err := db.Search("query", 3)
if err != nil {
    if errors.Is(err, search.ErrNotFound) {
        log.Println("No matching documents found.")
    } else {
        log.Fatalf("Error searching documents: %v", err)
    }
}

for _, result := range results {
    fmt.Println(result)
}
```

### Checking Database Length

```go
fmt.Printf("Number of documents in database: %d\n", db.Length())
```

## Error Handling

The `search` package comes with several predefined errors:
- `ErrTitleIsEmpty`: Returned when a document with an empty title is processed.
- `ErrQueryIsEmpty`: Returned when an empty query is passed to the `Search` function.
- `ErrDocumentIsNil`: Returned when attempting to add a `nil` document.
- `ErrNotFound`: Returned when no matching documents are found.