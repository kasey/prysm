### Added

- Added an error field to log `Finished building block`. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14696)
- Implemented a new `EmptyExecutionPayloadHeader` function. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14713)
- Added proper gas limit check for header from the builder. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14707)
- `Finished building block`: Display error only if not nil. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14722)
- Added light client feature flag check to RPC handlers. [PR](https://github.com/prysmaticlabs/prysm/pull/14736)
- Added support to update target and max blob count to different values per hard fork config. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14678)
- Log before blob filesystem cache warm-up. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14735)


### Changed

- Process light client finality updates only for new finalized epochs instead of doing it for every block. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14713)
- Refactor subnets subscriptions. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14711)
- Refactor RPC handlers subscriptions. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14732)
- Go deps upgrade, from `ioutil` to `io`. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14737)
- Move successfully registered validator(s) on builder log to debug. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14735)

### Fixed

- Added check to prevent nil pointer deference or out of bounds array access when validating the BLSToExecutionChange on an impossibly nil validator. [[PR]](https://github.com/prysmaticlabs/prysm/pull/14705)
