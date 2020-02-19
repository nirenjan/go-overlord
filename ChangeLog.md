# Changelog
All notable changes to Evil Overlord will be documented in this file.

The format is based on [Keep a Changelog] and this project adheres to [Semantic
Versioning]

## [Unreleased]
### Changed
- Evil Overlord is now written in Go, resulting in a single binary
- Backup/Restore functionality now uses gzip compression, rather than
  LZMA (XZ). This is because LZMA is not available in the Go standard
  library.
- Version is now a separate command, rather than a command line switch.
- No longer uses Git in the backend. This prefers to not call out to external
  applications, except for the editor.
- The terminal package expects this to be run in a Unix environment such as
  Linux or macOS. This has not been tested with Windows, and probably never
  will.

## [0.1.0] - 2018-08-06
### Added
- First release of Evil Overlord Personal Assistant
- Add support for journal entries
- Journal entries support multiple tags
- List all entries
- List all tags
- Show or delete a specific entry
- Display all entries, optionally filtered by a set of tags
- Support for exporting all Overlord data
- Support for importing from an existing Overlord backup

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

[Unreleased]: https://github.com/nirenjan/overlord/compare/0.1.0..HEAD
[0.1.0]: https://github.com/nirenjan/overlord/releases/tag/0.1.0
