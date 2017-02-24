# Ipê

A replacement for `ls` with some special features, like tree view and Git integration. Works as terminal program or Go library.

Inspired by [jacwah/oak][1] and [ogham/exa][2].

### To-do list

- [x] List directory contents
- [ ] Identify and organize column sizes in long views
- [ ] Add (almost) all flags and args `ls` has
  - [x] `-a`, `--all`
  - [ ] `--author`
  - [ ] `-b`, `--escape`
  - [ ] `--block-size`
  - [ ] `-B`, `--ignore-backups`
  - [ ] `-C`
  - [x] `--color`
  - [ ] `-d`, `--directory`
  - [ ] `-D`, `--dired`
  - [x] `-F`, `--classify`
  - [ ] `--format`
  - [ ] `--full-time`
  - [ ] `-g`
  - [ ] `-G`, `--no-group`
  - [x] `-h`, `--human-readable`
  - [ ] `-H`, `--dereference-command-line`
  - [ ] `--dereference-command-line-symlink-to-dir`
  - [ ] `--indicator-style`
  - [ ] `-i`, `--inode`
  - [x] `-I`, `--ignore`
  - [ ] `-k`
  - [x] `-l`
  - [ ] `-L`, `--dereference`
  - [ ] `-m`
  - [ ] `-n`, `--numeric-uid-gid`
  - [ ] `-N`, `--literal`
  - [ ] `-o`
  - [ ] `-p`, `--file-type`
  - [ ] `-q`, `--hide-control-chars`
  - [ ] `--show-control-chars`
  - [ ] `-Q`, `--quote-name`
  - [ ] `--quoting-style`
  - [x] `-r`, `--reverse`
  - [x] `-R`, `--recursive`
  - [ ] `-s`, `--size`
  - [ ] `--sort`
  - [ ] `--time`
  - [ ] `--time-style`
  - [ ] `-T`, `--tabsize`
  - [ ] `-w`, `--width`
  - [x] `-x`
  - [ ] `-1`
  - [ ] `--lcontext`
  - [ ] `-Z`, `--context`
  - [ ] `--scontext`
  - [x] `--help`
  - [ ] `--version`
- [ ] Make special features
  - [x] Tree view (-t)
  - [x] Separator of columns in long view
  - [ ] Flag to specify maximum depth of recursion
  - [ ] Flag to sort by an column/field
  - [ ] Flag to get one item per line
  - [ ] Flag to filter entries (-f)
- [ ] Define colors
- [ ] Add [Git integration][3]
  - [ ] Ignore "Git ignored" files by default
  - [ ] Show files' Git status 
- [ ] Change it into a lib

### License

This project code is in the public domain. See the [LICENSE file][4].

[1]: https://github.com/jacwah/oak/
[2]: https://github.com/ogham/exa/
[3]: https://github.com/libgit2/git2go
[4]: https://github.com/Nhanderu/ipe/blob/master/LICENSE
