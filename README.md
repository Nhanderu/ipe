# Ipê (deprecated)

![Deprecated][badge-1-img]
[![License][badge-2-img]][badge-2-link]
[![go.dev][badge-3-img]][badge-3-link]
[![Go Report Card][badge-4-img]][badge-4-link]

A replacement for `ls` with some special features, like tree view and Git
integration. Works as terminal program or Go library.

Inspired by [jacwah/oak][1] and [ogham/exa][2].

## Deprecated

Don't use it! I made this code trying to learn Go but I never finished it and I
don't see any reason for doing it now.

## To-do list

- [x] List directory contents
- [x] Identify and organize column sizes in long views
- [ ] Make special features
  - [x] Tree view (-t)
  - [x] Separator of columns in long view (-S)
  - [x] Flag to specify maximum depth of recursion (-D)
  - [x] Flag to sort by an column/field (-s)
  - [x] Flag to filter entries (-f)
  - [x] Flag to show headers on long view (-h)
  - [ ] Differentiate files types
  - [x] Flag to show directories first (--dirs-first)
  - [x] Accept more than one value in filter and ignore flags
  - [x] Flag to show number of hard links in long view (--links)
  - [x] Flag to show number of file system blocks in long view (--blocks)
  - [x] Flag to show group in long view (--group)
- [ ] Define colors
- [ ] Add [Git integration][3]
  - [ ] Ignore "Git ignored" files by default
  - [ ] Show files' Git status
- [x] Change it into a lib
- [x] Create formatters
- [ ] Get inode, user and group in Windows
- [ ] Define column alignment in long view

## License

This project code is in the public domain. See the [LICENSE file][4].

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you shall be in the public domain, without any
additional terms or conditions.

[1]: https://github.com/jacwah/oak/
[2]: https://github.com/ogham/exa/
[3]: https://github.com/libgit2/git2go
[4]: ./LICENSE

[badge-1-img]: https://img.shields.io/badge/code-deprecated-critical?style=flat-square
[badge-2-img]: https://img.shields.io/github/license/Nhanderu/ipe?style=flat-square
[badge-2-link]: https://github.com/Nhanderu/ipe/blob/master/LICENSE
[badge-3-img]: https://img.shields.io/badge/go.dev-reference-007d9c?style=flat-square&logo=go&logoColor=white
[badge-3-link]: https://pkg.go.dev/github.com/Nhanderu/ipe
[badge-4-img]: https://goreportcard.com/badge/github.com/Nhanderu/ipe?style=flat-square
[badge-4-link]: https://goreportcard.com/report/github.com/Nhanderu/ipe
