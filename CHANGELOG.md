## v1.0.2

- fixes missing normalization of aliases when populating the alias map in the bibleref.Table struct, ensuring that all aliases are normalized for consistent lookup
- fixes missing normalized OSIS keys in the alias map of the bibleref.Table struct, ensuring that normalized OSIS keys are always present for lookup
- adds end-to-end tests for various biblical reference formats, including normalization of book names and handling of whitespace, to ensure the robustness of the parsing and normalization logic in the bibleref package

## v1.0.1

- fixes missing json tags in bibleref.Book struct

## v1.0.0

- adds `bibleref` package for parsing and normalizing biblical references
- adds `rbref` package for parsing and normalizing references to The Rule of St Benedict
