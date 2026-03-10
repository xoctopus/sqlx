
<a name="HEAD"></a>
## [HEAD](https://github.com/xoctopus/sqlx/compare/v0.2.1...HEAD)

> 2026-03-09

### Chore

* clean code

### Ci

* update Makefile

### Feat

* session name uniqueness


<a name="v0.2.1"></a>
## [v0.2.1](https://github.com/xoctopus/sqlx/compare/v0.2.0...v0.2.1)

> 2026-03-09

### Chore

* add drivers for pg and sqlite

### Doc

* update CHANGELOG
* update CHANGELOG

### Feat

* **session:** session and adaptor with underlying sql.DB


<a name="v0.2.0"></a>
## [v0.2.0](https://github.com/xoctopus/sqlx/compare/v0.1.27...v0.2.0)

> 2026-03-08

### Chore

* **deps:** bump xoctopus/{x,genx,pkgx,typx} to latest


<a name="v0.1.27"></a>
## [v0.1.27](https://github.com/xoctopus/sqlx/compare/v0.1.26...v0.1.27)

> 2026-03-06

### Fix

* fix confused column diff from catalog and user define


<a name="v0.1.26"></a>
## [v0.1.26](https://github.com/xoctopus/sqlx/compare/v0.1.25...v0.1.26)

> 2026-02-02

### Chore

* missing adaptor error


<a name="v0.1.25"></a>
## [v0.1.25](https://github.com/xoctopus/sqlx/compare/v0.1.24...v0.1.25)

> 2026-01-30

### Chore

* **deps:** bump github.com/xoctopus/x from 0.2.11 to 0.2.12 ([#8](https://github.com/xoctopus/sqlx/issues/8))
* **deps:** bump github.com/xoctopus/logx from 0.1.1 to 0.1.2 ([#9](https://github.com/xoctopus/sqlx/issues/9))
* **deps:** bump github.com/xoctopus/confx from 0.2.18 to 0.2.19 ([#10](https://github.com/xoctopus/sqlx/issues/10))

### Feat

* **migrator:** support migrator output to file


<a name="v0.1.24"></a>
## [v0.1.24](https://github.com/xoctopus/sqlx/compare/v0.1.23...v0.1.24)

> 2026-01-27

### Feat

* **builder:** add ForUpdate/ForShare with modifiers


<a name="v0.1.23"></a>
## [v0.1.23](https://github.com/xoctopus/sqlx/compare/v0.1.22...v0.1.23)

> 2026-01-26

### Feat

* **builder:** add SKIP_LOCKED addition


<a name="v0.1.22"></a>
## [v0.1.22](https://github.com/xoctopus/sqlx/compare/v0.1.21...v0.1.22)

> 2026-01-16

### Feat

* **session:** add endpoint option SetDefault


<a name="v0.1.21"></a>
## [v0.1.21](https://github.com/xoctopus/sqlx/compare/v0.1.20...v0.1.21)

> 2026-01-16

### Fix

* **hack:** add WithSession cleanup


<a name="v0.1.20"></a>
## [v0.1.20](https://github.com/xoctopus/sqlx/compare/v0.1.19...v0.1.20)

> 2026-01-16

### Chore

* **deps:** bump github.com/xoctopus/genx from 0.1.14 to 0.1.16 ([#6](https://github.com/xoctopus/sqlx/issues/6))

### Doc

* update CHANGELOG

### Feat

* **session:** impl endpoint LivenessCheck


<a name="v0.1.19"></a>
## [v0.1.19](https://github.com/xoctopus/sqlx/compare/v0.1.18...v0.1.19)

> 2026-01-15

### Chore

* hacking testing for session injection


<a name="v0.1.18"></a>
## [v0.1.18](https://github.com/xoctopus/sqlx/compare/v0.1.17...v0.1.18)

> 2026-01-15

### Feat

* **session:** use universal endpoint


<a name="v0.1.17"></a>
## [v0.1.17](https://github.com/xoctopus/sqlx/compare/v0.1.16...v0.1.17)

> 2026-01-11

### Chore

* format config and regen
* **deps:** bump codecov/codecov-action from 4 to 5 ([#1](https://github.com/xoctopus/sqlx/issues/1))
* **deps:** bump actions/setup-go from 5 to 6 ([#2](https://github.com/xoctopus/sqlx/issues/2))
* **deps:** bump actions/checkout from 4 to 6 ([#3](https://github.com/xoctopus/sqlx/issues/3))
* **deps:** bump github.com/xoctopus/x from 0.2.10 to 0.2.11 ([#4](https://github.com/xoctopus/sqlx/issues/4))

### Fix

* **adaptor:** check rollback error in tx failure


<a name="v0.1.16"></a>
## [v0.1.16](https://github.com/xoctopus/sqlx/compare/v0.1.15...v0.1.16)

> 2026-01-04

### Chore

* update dependencies


<a name="v0.1.15"></a>
## [v0.1.15](https://github.com/xoctopus/sqlx/compare/v0.1.14...v0.1.15)

> 2026-01-04

### Chore

* update dependencies


<a name="v0.1.14"></a>
## [v0.1.14](https://github.com/xoctopus/sqlx/compare/v0.1.13...v0.1.14)

> 2026-01-04

### Chore

* update dependencies


<a name="v0.1.13"></a>
## [v0.1.13](https://github.com/xoctopus/sqlx/compare/v0.1.12...v0.1.13)

> 2026-01-02

### Feat

* **session:** migrate endpoint


<a name="v0.1.12"></a>
## [v0.1.12](https://github.com/xoctopus/sqlx/compare/v0.0.11...v0.1.12)

> 2025-12-30

### Chore

* add linter and fixing


<a name="v0.0.11"></a>
## [v0.0.11](https://github.com/xoctopus/sqlx/compare/v0.0.10...v0.0.11)

> 2025-12-29

### Fix

* diff column datatype and timestamp default value


<a name="v0.0.10"></a>
## [v0.0.10](https://github.com/xoctopus/sqlx/compare/v0.0.9...v0.0.10)

> 2025-12-28


<a name="v0.0.9"></a>
## [v0.0.9](https://github.com/xoctopus/sqlx/compare/v0.0.8...v0.0.9)

> 2025-12-28

### Fix

* **adaptor:** mysql dsn parser


<a name="v0.0.8"></a>
## [v0.0.8](https://github.com/xoctopus/sqlx/compare/v0.0.7...v0.0.8)

> 2025-12-28


<a name="v0.0.7"></a>
## [v0.0.7](https://github.com/xoctopus/sqlx/compare/v0.0.6...v0.0.7)

> 2025-12-28

### Fix

* **adaptor:** check datatype descriptor


<a name="v0.0.6"></a>
## [v0.0.6](https://github.com/xoctopus/sqlx/compare/v0.0.5...v0.0.6)

> 2025-12-28

### Feat

* **devpkg:** use unsafe exposer for speeding generating


<a name="v0.0.5"></a>
## [v0.0.5](https://github.com/xoctopus/sqlx/compare/v0.0.4...v0.0.5)

> 2025-12-28

### Feat

* **types:** enrich basic types


<a name="v0.0.4"></a>
## [v0.0.4](https://github.com/xoctopus/sqlx/compare/v0.0.3...v0.0.4)

> 2025-12-25

### Fix

* session injections


<a name="v0.0.3"></a>
## [v0.0.3](https://github.com/xoctopus/sqlx/compare/v0.0.2...v0.0.3)

> 2025-12-24

### Feat

* curd operations generate
* **devpkg:** curd generated


<a name="v0.0.2"></a>
## [v0.0.2](https://github.com/xoctopus/sqlx/compare/v0.0.1...v0.0.2)

> 2025-12-15

### Chore

* regen
* regenerated

### Ci

* fix ci-cover with hack_dep cover filter
* fix ci-cover with hack_dep
* ignore testing for hack package

### Feat

* table for testing and hack test checker
* **adaptor:** mysql driver
* **session:** exports

### Fix

* romove column tag context once attr
* make column default option literal; KeyDefine support column name or field name
* **builder:** remove context depends on table scanning
* **devpkg:** remove context depends on table scanning

### Refact

* rename OptionsFieldNames => OptionsNames
* rename OptionsFieldNames => OptionsNames

### Test

* use fixed hack data
* **frag:** unit test


<a name="v0.0.1"></a>
## v0.0.1

> 2025-12-03

### Chore

* bump dependencies and adaption
* bump dependencies
* bump dependencies
* clearify frag APIs

### Ci

* add .github ci flow

### Feat

* table def generator
* table def generator
* sql fragments
* **builder:** mysql-on_duplicate pg-on_conflict
* **builder:** sql builder APIs
* **internal:** sql scanner,driver,adapter; scoped builder; error defines
* **sqltypes:** feat(sqltypes): add fundamental sql types
* **sqltypes:** add fundamental sql types

### Fix

* fix unit tests

