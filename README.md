# canonref

Canonical reference parsing, normalization, and validation for structured religious and liturgical texts.

This repository provides **small, dependency-light Go packages** for working with
canonical citations (e.g. _Rule of St Benedict_ references) in a precise, auditable way.
It is intended to be used by multiple tools and CLIs, not as an application itself.

---

## Purpose

Many texts central to Christian tradition are cited using compact, conventional forms:

- `RB Prol. 1–7`
- `RB 48.1–9`
- `RB 4`

These references carry **implicit structure** (section, chapter, verse, range) that is
often left unmodeled or treated as opaque strings.

`canonref` exists to:

- parse canonical citations into structured data
- validate them against known textual constraints
- normalize them into a single, stable string form
- make downstream tooling precise rather than heuristic

---

## Packages

### `rbref`

Canonical references for **The Rule of St Benedict**.

Supported forms include:

- Prologue references
  - `RB Prol. 1`
  - `RB Prol. 1–7`
- Chapter + verse references
  - `RB 48.1`
  - `RB 48.1–9`
- Chapter-only references
  - `RB 4`

The parser enforces:

- valid chapter ranges (1–73)
- positive verse numbers
- well-formed ranges
- explicit distinction between Prologue and chapters

All references are normalized (e.g., hyphen → en-dash) and rendered in a single
canonical format.

Example:

```go
ref, err := rbref.Parse("RB 48.1-9")
if err != nil {
    // handle error
}

fmt.Println(ref.String())
// Output: RB 48.1–9
```

### `bibleref`

Canonical references for **the Bible**.

Supported forms include:

- Book abbreviations and full names (case-insensitive)
  - `Prov 31:10–31`
  - `Proverbs 31:10–31`
  - `PRO 31:10-31`
- Multi-testament coverage (Old Testament, New Testament, Apocrypha)
  - `1 Samuel 17:4-11`
  - `Wis 1:1` (Wisdom of Solomon)
- Chapter-only references
  - `Prov 31`
- Single verses and ranges
  - `Matt 5:3`
  - `John 1:1–14`

The parser enforces:

- valid book names via an alias table
- valid chapter ranges for each book
- positive verse numbers
- well-formed ranges
- proper normalization of punctuation (hyphens → en-dashes)

All references are normalized and rendered in canonical format (e.g., `Prov 31:10–31`).

Example:

```go
books := []bibleref.Book{
    {OSIS: "Prov", Name: "Proverbs", Aliases: []string{"prov", "pro"}, Testament: "OT", Order: 20, Chapters: 31},
    // ... more books
}
tbl, err := bibleref.NewTable(books)

ref, err := bibleref.Parse("Proverbs 31:10-31", tbl)
if err != nil {
    // handle error
}

fmt.Println(ref.String())
// Output: Prov 31:10–31
```

---

## Design Principles

- **Library-first**: packages must remain useful without any CLI.
- **Minimal dependencies**: standard library preferred.
- **Explicit structure**: avoid stringly-typed representations.
- **Canonical output**: one true string form per reference.
- **Clear failure modes**: parsing vs validation errors are distinct.

---

## Relationship to Other Repositories

Typical usage within the ecosystem:

- **`liturgical-time-index`**
  Consumes `rbref` to attach structured RB references to calendar days.

- **Citation resolver CLIs**
  Parse user input and resolve it against corpora using `canonref` types.
