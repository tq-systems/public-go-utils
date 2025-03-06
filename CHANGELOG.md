## [3.1.0] - 2025-03-06
### Changed
- the Publish method of the MQTT client now uses MarshalVT for better performance, whenever possible

## [3.0.0] - 2025-03-03
### Changed
- switch protobuf library to google.golang.org/protobuf

## [2.0.6] - 2024-12-19

## [2.0.5] - 2024-11-27
### Added
- Loglevel Notice and Noticef

## [2.0.4] - 2024-11-15
### Fixed
- Prevent freezing in mqtt, when the MQTT broker fails to confirm some subscriptions and publications.

## [2.0.3] - 2024-11-04
### Fixed
- A race condition leading to deadlocks in PublishRaw and doSubscribe has been removed.

## [2.0.2] - 2024-06-14
### Fixed
- SetStatusIfIdle of status package, it always returned an error before

## [2.0.1] - 2024-06-06
### Changed
- add v2 to module in go.mod

## [2.0.0] - 2024-06-06
### Changed
- cleanly return error to caller instead of panicking / logging or just returning nil objects

## [1.3.0] - 2024-04-11
### Changed
- migrated from github.com/golang/mock to go.uber.org/mock

## [1.2.0] - 2024-02-20
### Added
- clock interface to encapsulate dependency on system time for better testability

## [1.1.1] - 2023-09-04
### Fixed
- mqtt: subscription to empty topic froze the process and now throws an error

## [1.1.0] - 2023-08-18
### Added
- auth package
- mqtt package
- rest package

### Changed
- file headers with copyright information

### Fixed
- linter issues

## [1.0.0] - 2023-03-20
