## 0.4.4 (2017-07-07)

Bug fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Fixed a bug that syntax error occurred when SQL contains partition queries

## 0.4.4 (2017-03-29)

Bug fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Fixed a bug that will get a difference when contains some partition information table columns

## 0.4.3 (2017-03-13)

Support partition table

### Added

- support partition table
 - only "LINEAR KEY", "LINEAR HASH" and "RANGE COLUMNS" are supported
 - supported alter only. drop and remove is not supported

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Fixed a bug that will get a difference when table names contain uppercase letters 

## 0.4.2 (2017-02-10)

Minor feature released

### Added

- resolve column position modification 

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Fix SEGV error 


## 0.4.1 (2017-02-08)

Add migrate table collation

### Added

- build
    - sync table and column character set

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Nothing


## 0.4.0 (2017-01-27)

Revert version

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- revert 0.2.4

### Fixed

- Nothing


## 0.3.1 (2016-12-12)

Bug fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- design command: STDOUT output now removed

### Fixed

- Enables specification of design option `-s`

## 0.3.0 (2016-12-12)

- design command: Change export format 
- design command: Add `-s` option
- Fix test
- Fix bug

### Added

- Add separate option to design command 
- Change STDOUT output to files for each tables (if `-s` option specified)

### Deprecated

- Nothing

### Removed

- design command: STDOUT output now removed

### Fixed

- Type translation bug fixed
- Test


## 0.2.6 (2016-12-02)

- Bug Fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Support bigint export


## 0.2.5 (2016-12-01)

- Bug Fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Support index comment


## 0.2.4 (2016-11-29)

- Bug Fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Unescape double escape string


## 0.2.3 (2016-11-28)

- Bug Fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Do not compare privileges


## 0.2.2 (2016-11-28)

- Bug Fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Fix error when using MySQL5.5


## 0.2.1 (2016-11-26)

- Bug Fix

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Always a difference between schemas with different names 


## 0.2.0 (2016-11-21)

- Add export subcommand
- Rename subcommand export to design 

### Added

- Added subcommand to export table data as CSV

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Adjust the maximum number of MySQL connections


## 0.1.1 (2016-11-15)

Bug fixed

### Added

- Nothing

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Quote the default value that needed to be quoted

## 0.1.0 (2016-11-15)

Initial release

### Added

- Add Fundamental features

### Deprecated

- Nothing

### Removed

- Nothing

### Fixed

- Nothing
