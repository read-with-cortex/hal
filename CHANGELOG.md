# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- `hal.NewEncoder` is now called `hal.NewJSONEncoder`
- `halforms.NewEncoder` is now called `halforms.NewJSONEncoder`
- `hal.NewJSONEncoder` and `halforms.NewJSONEncoder` now have an `Encode` function instead of `toJSON`. Furthermore, the `Encoder` interface has an `Encode` function instead of `toJSON`.

## [0.1.0] - 2024-02-15

- Initial Release
