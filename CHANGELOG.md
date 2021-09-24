# Changelog

## July 6, 2021 v3.0.0
### Added
- Added new method to merge duplicate customers. More details in [API Documentation](https://customer.io/docs/api/#operation/merge)

## July 6, 2021 v3.0.0
### Changed
- `trackAnonymous` now requires an `anonymous_id` parameter and will no longer trigger campaigns. If you previously used anonymous events to trigger campaigns, you can still do so [directly through the API](https://customer.io/docs/api/#operation/trackAnonymous). We now refer to anonymous events that trigger campaigns as ["invite events"](https://customer.io/docs/anonymous-events/#anonymous-or-invite). 

## June 16, 2021 v2.2.0

### Added
- Modules support
- Updated license date
- Upgraded version support

## March 24, 2021 v2.1.0
### Added
- Support for EU region
- Allow using custom `*http.Client`
### Removed
### Changed
- `customerio.NewAPIClient` and `customerio.NewTrackClient`  have a new variadic parameter for options in order to choose US/EU region and/or customer HTTP client.

## December 3, 2020 v2.0.0
### Added
- Support for transactional api

### Removed
- Manual segmentation functions `AddCustomersToSegment` & `RemoveCustomersFromSegment`

### Changed
- ID fields in requests are url escaped
- Improved validations for required fields
- Updated `CustomerIO` struct to use a URL field instead of separate Host and SSL fields
- Improved test suite
