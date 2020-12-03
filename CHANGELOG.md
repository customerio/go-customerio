# Changelog

## December 3, 2020
### Added
- Support for transactional api

### Removed
- Manual segmentation functions `AddCustomersToSegment` & `RemoveCustomersFromSegment`

### Changed
- ID fields in requests are url escaped
- Improved validations for required fields
- Updated `CustomerIO` struct to use a URL field instead of separate Host and SSL fields
- Improved test suite
