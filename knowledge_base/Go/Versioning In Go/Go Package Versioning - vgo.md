#go #blog-post-notes

[Go & Versioning](https://research.swtch.com/vgo)

## Reproducible Builds
- Versioning enables reproducible builds, 
- so that if we tell to someone to try a specific or latest version of our program, 
	- we know we're going to get not just the latest version of our program, 
	- but the **EXACT** same version of all packages our program depends on
- so that we and anyone will build completely equivalent binaries

- Versioning will also let us ensure that a program builds exactly same way in future, as it does today
	- even when there are newer versions of our dependencies, the go command shouldn't start using them until asked

## 4 steps of `vgo`
### 1. *import compatibility rule* 
- establish the expectation that newer versions of a package with a given import path should be backward compatible with older versions
### 2. *minimal version selection*
- an algorithm to choose which package versions are used in a given build
### 3. Go modules
- a group of packages versioned as a single unit 
- and that declare the minimum requirements that must be satisfied by their dependencies
### 4. retrofit all this into existing go command
- so that basic workflows do not change significantly


## The Import Compatibility Rule
- All pain in package management systems is caused by trying to tame **incompatibility**
- Example:
	- if package-B declares that it requires package D-**v6 or later**
		- most package management systems allow this **later thing**
	- and if package-C declares that it requires package D-**v4 or lesser** and not v5 or later
		- most package management systems allow this **lesser thing**
	- and if we are writing package-A and we want to use both package-B and package-C, then we are out of luck
	- there is no single version of D that can be chosen to build both B and C into A
- thus, this inevitably leads to large programs not building
- `vgo` solves this by requiring the package authors to follow the *import compatibility rule*:

> If an old package and a new package have **same import path**,
> the new package **must be backwards compatible** with old package

- If a complete break is required, create a new package with a new import path
- Developers today expect to use semantic versioning to express such a break
```go
import "github.com/go-yaml/yaml/v2"
```
- Creating v2.0.0, which in semantic versioning denotes a major break, thus creates a new package with a new import path, as required by import compatibility
- Because each major version has a different import path, a given go executable might contain one of each major version
- This is expected and desirable
- It keeps programs building and allows parts of a very large program to update from v1 to v2 independently

- Expecting authors to follow import compatibility rule makes it possible. to avoid taming incompatibility, making the overall system exponentially simpler
- In practice, of course, despite the best efforts of authors, updates within same major version do occasionally break users
- Thus, it's important to use an upgrade mechanism that doesn't upgrade too quickly - **Minimal Version Selection**


## Minimal Version Selection
- Nearly all package managers today use the newest allowed version of packages involved in the build
- This is the wrong default, for two important reasons
	- The meaning of "newest allowed version" can change due to external events
		- maybe tonight someone will introduce a new version of some dependency and the tomorrow same sequence of commands u ran today would produce a different result
	- to override this default, developers spend time telling package managers "no, don't use X version"
		- then package manager spends it time searching for way not to use X version
- MVS defaults to using the oldest allowed version of every package involved in the build
	- This decision does not change from today to tomorrow, because no older version
	- even better, to override this default, developers spend their time telling the package manager, "no, use at least Y", and then the package manager can trivially decide which version to use
- MVS delivers reproducible builds by default, without a lock file
- Instead of users saying "no, that's too new", they can only say "no, that's too old".
- Import compatibility is key to MVS

## Go Modules
- A go module is a collection of packages sharing a common **import-path-prefix**, known as **module path**
- Module is the unit of versioning, and module versions are written as semantic version strings
- When developing using Git, developers will define a new semantic version of a module by adding a tag to the module's git repository
- Although semantic versions are strongly preferred, referring to specific commits is also supported
- A module is defined in a new file called `go.mod` 
- `go.mod` file is used to define the minimum version requirements of other modules it depends on
```go.mod
// My hello, world.

module "rsc.io/hello"

require (
    "golang.org/x/text" v0.0.0-20180208041248-4e4a3210bb54
    "rsc.io/quote" v1.5.2
)
```
- This file defines a module, identified by path `rsc.io/hello`
	- which itself depends on two other modules
		- "golang.org/x/text" `v0.0.0-20180208041248-4e4a3210bb54`
		- "rsc.io/quote" `v1.5.2`
- A build of a module by itself will always use specific versions of required dependencies listed in go.mod file
	- As part of a larger build, it will only use a newer version if something else in the build requires it
- To name untagged commits, the pseudo-version `v0.0.0-yyyymmddhhmmss-commitHash` identifies a specific commit made on the given date
- in semantic versioning, this string corresponds to a v0.0.0 prerelease, with prerelease identifier `yyyymmddhhmmss-commit`
- Semantic versioning precedence rules order such prereleases by string comparison
- Placing the date first in the pseudo-version syntax ensures that string comparison matches date comparison


- in addition to requirements (in the `require` block in go.mod file), go.mod files can specify exclusions and replacements mentioned in previous section