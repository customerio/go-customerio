# Changelog

## [v3.2.0](https://github.com/customerio/go-customerio/compare/v3.1.0...v.3.2.0) (2021-10-04)
### Added
- **client:** adds a default User-Agent header on requests and the option to set a custom User-Agent value. ([2123497](https://github.com/customerio/go-customerio/commit/212349768ba234d6c4ad3684aa6450f770f35cb8))

## [v3.1.0](https://github.com/customerio/go-customerio/compare/3.0.0...v3.1.0) (2021-09-27)
### Added
- Added new method to merge duplicate customers. More details in [API Documentation](https://customer.io/docs/api/#operation/merge)

## [v3.0.0](https://github.com/customerio/go-customerio/compare/v2.2.0...3.0.0) (2021-07-08)
### Changed
- `trackAnonymous` now requires an `anonymous_id` parameter and will no longer trigger campaigns. If you previously used anonymous events to trigger campaigns, you can still do so [directly through the API](https://customer.io/docs/api/#operation/trackAnonymous). We now refer to anonymous events that trigger campaigns as ["invite events"](https://customer.io/docs/anonymous-events/#anonymous-or-invite). 

## [v2.2.0](https://github.com/customerio/go-customerio/compare/2.1.0...v2.2.0) (2021-06-17)

### Added
- Modules support
- Updated license date
- Upgraded version support

## [v2.1.0](https://github.com/customerio/go-customerio/compare/v2.0.0...2.1.0) (2021-03-29)
### Added
- Support for EU region
- Allow using custom `*http.Client`

### Changed
- `customerio.NewAPIClient` and `customerio.NewTrackClient`  have a new variadic parameter for options in order to choose US/EU region and/or customer HTTP client.

## [v2.0.0](https://github.com/customerio/go-customerio/compare/v1.2.0...v2.0.0) (2020-12-03)
### Added
- Support for transactional api

### Removed
- Manual segmentation functions `AddCustomersToSegment` & `RemoveCustomersFromSegment`

### Changed
- ID fields in requests are url escaped
- Improved validations for required fields
- Updated `CustomerIO` struct to use a URL field instead of separate Host and SSL fields
- Improved test suite

## [v1.2.0](https://github.com/customerio/go-customerio/compare/v1.1.0...v1.2.0) (2020-09-08)
### Fixed
- Read response body regardless of content length


## [v1.1.0](https://github.com/customerio/go-customerio/compare/v1.0.0...v1.1.0) (2019-10-24)
### Added
- `TrackAnonymous` method
- `AddDevice` and `DeleteDevice` methods
- `AddCustomersToSegment` and `RemoveCustomersFromSegment` methods

### Changed
- Increase default client connection pool size to 100

## [v1.0.0](https://github.com/customerio/go-customerio/compare/4a9e70a...v1.0.0) (2017-09-12)
### Added
- Initial API client library
  - Create Identify call
  - Create Track call
  - Create Customer Delete call
  - Test suite for API calls

## [v0.0.0](https://github.com/customerio/go-customerio/commit/4a9e70a) (2015-01-12)
### Added
- Initial commit
