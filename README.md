# Ipê

A replacement for `ls` with some special features, like tree view and Git integration. Works as terminal program or Go library.

Inspired by [jacwah/oak][1] and [ogham/exa][2].

### To-do list

- [x] List directory contents
- [ ] Identity and organize column sizes in long views
- [ ] Add all flags and args `ls` has
  - [x] `-a`, `--all`
  - [ ] `-A`, `--almost-all` (IGNORED)
  - [ ] `--author`
  - [ ] `-b`, `--escape`
  - [ ] `--block-size`
  - [ ] `-B`, `--ignore-backups`
  - [ ] `-c`
  - [ ] `-C`
  - [x] `--color`
  - [ ] `-d`, `--directory`
  - [ ] `-D`, `--dired`
  - [ ] `-f`
  - [x] `-F`, `--classify`
  - [ ] `--format`
  - [ ] `--full-time`
  - [ ] `-g`
  - [ ] `-G`, `--no-group`
  - [x] `-h`, `--human-readable`
  - [ ] `--si` (IGNORED)
  - [ ] `-H`, `--dereference-command-line`
  - [ ] `--dereference-command-line-symlink-to-dir`
  - [ ] `--indicator-style`
  - [x] `-i`, `--inode`
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
  - [ ] `-S`
  - [ ] `--sort`
  - [ ] `--time`
  - [ ] `--time-style`
  - [ ] `-t`
  - [ ] `-T`, `--tabsize`
  - [ ] `-u`
  - [ ] `-U`
  - [ ] `-v`
  - [ ] `-w`, `--width`
  - [ ] `-x`
  - [ ] `-X`
  - [ ] `-1`
  - [ ] `--lcontext`
  - [ ] `-Z`, `--context`
  - [ ] `--scontext`
  - [x] `--help`
  - [ ] `--version`
- [ ] Make special features
  - [ ] Get all files by default, but separate dotfiles from the rest
  - [ ] `-a` to don't separate the files
  - [ ] `-A` to don't get dotfiles
  - [ ] Tree view
  - [x] Separator of columns in long view
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
[4]: https://github.com/Nhanderu/ype/blob/master/LICENSE
