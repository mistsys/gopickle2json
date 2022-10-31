# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.2] - 2022-10-31
### Fixed
- Python tuples become JSON lists.

## [0.2.1] - 2022-10-27
### Fixed
- Added types.None, which was missing.

## [0.2.0] - 2022-10-27
### Changed
- Make every python object implement types.Object, and types.Object has a JSON()
  method which emits the equivalent JSON. This necessitated wrapping all the
  types which were bare Go, like bool, int, string, into named types.
  In the end, I don't emit JSON directly. I still construct the intermediate
  python objects, and then ask them to emit. This creates lots of garbage to gc,
  so I might re-explore directly emitted the JSON as as unpickle.

### forked - 2022-10-27
- forked from nlpodyssey/gopickle (or really, a child of that, nsd20463/gopickle
  b/c I wanted the improvements there)

## [0.1.0] - 2021-01-06
### Added
- More and better documentation
- `OrderedDict.MustGet()`
- `Dict.MustGet()`
- `pytorch.LoadWithUnpickler()` which allows loading PyTorch modules using a
  custom unpickler.
- Handle legacy method `torch.nn.backends.thnn._get_thnn_function_backend` when
  loading pytorch modules.

### Changed
- `FrozenSet` implementation was modified, avoiding confusion with `Set`.
- Replace build CI job with tests and coverage
- `Dict` has been reimplemented using a slice, instead of a map, because in Go
  not all types can be map's keys (e.g. slices).
- Use Go version `1.15`

### Removed
- Unused method `List.Extend`

## [0.0.1-alpha.1] - 2020-05-23
### Fixed
- Modify GitHub Action steps `Build` and `Test` including all sub-packages.

## [0.0.1-alpha.0] - 2020-05-23
### Added
- Initial implementation of `types` package
- Initial implementation of `pickle` package
- Initial implementation of `pytorch` package
