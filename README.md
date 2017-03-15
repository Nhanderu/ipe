# Ipê

[![GoDoc][tag1img]][tag1link]
[![Go Report Card][tag2img]][tag2link]

A replacement for `ls` with some special features, like tree view and Git integration. Works as terminal program or Go library.

Inspired by [jacwah/oak][1] and [ogham/exa][2].

### To-do list

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

### License

This project code is in the public domain. See the [LICENSE file][4].

[1]: https://github.com/jacwah/oak/
[2]: https://github.com/ogham/exa/
[3]: https://github.com/libgit2/git2go
[4]: https://github.com/Nhanderu/ipe/blob/master/LICENSE

[tag1img]: https://godoc.org/github.com/Nhanderu/ipe?status.png
[tag1link]: https://godoc.org/github.com/Nhanderu/ipe
[tag2img]: https://goreportcard.com/badge/github.com/Nhanderu/ipe
[tag2link]: https://goreportcard.com/report/github.com/Nhanderu/ipe