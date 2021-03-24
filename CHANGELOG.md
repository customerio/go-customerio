# Changelog

## March 24, 2021
### Added
- Support for EU region
- Allow using custom `*http.Client`
### Removed
### Changed
- `customerio.NewAPIClient` and `customerio.NewTrackClient`  have a new variadic parameter for options in order to choose US/EU region and/or customer HTTP client.

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
