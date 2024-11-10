# **ep**: *"Emacs passwords"*

**ep** is a simple password management tool for the command-line that is backed by Emacs.

It has no dependencies other than a copy of Emacs that is of at least version `23.1`. Since **ep** can be built into a statically-linked Go executable, it should "just work" on Linux, MacOS, and Windows.

## Supported password management back-ends

Behind the scenes, your passwords can be stored in a variety of different formats. For example:

- An encrypted `~/.authinfo.gpg` or `~/.netrc.gpg` file.
- Inside the encrypted [`~/.password-store`](https://www.passwordstore.org).

You can even use multiple of these back-ends simultaneously, and **ep** will search through all of your Emacs `auth-sources` until it finds a matching password. The **ep** tool intentionally does not have a way to configure any of this. Instead, it will just respect your Emacs `auth-source` configuration. Please see the official Emacs [`auth-source`](https://www.gnu.org/software/emacs/manual/html_mono/auth.html) documentation for details.

## Installation

For now, you must clone the repository and build the executable yourself. To do this, you will need to install [Go](https://go.dev) and run `go build`. The resulting binary located at `./ep` can then be copied into a folder which is in your `$PATH`.

**ep** assumes that you have [Emacs](https://www.gnu.org/software/emacs/) installed and that it is available on your `$PATH`.

## Usage

One goal of **ep** is to be a drop-in replacement for the [`pass`](https://www.passwordstore.org) command, but this is not the case yet. For now, use `ep --help` to view the full documentation of its currently-supported functionality.

### Get a password for user `me` at domain `example.com`

```
ep me@example.com
```

Or, alternatively,

```
ep show me@example.com
```

### Add/edit/delete passwords

TODO: I have not yet added this functionality to **ep**.

However, if your passwords are stored in an encrypted `~/.authinfo.gpg` or `~/.netrc.gpg` file, then there is an easy way to edit your passwords from within Emacs:

- Include [`(epa-file-enable)`](https://www.gnu.org/software/emacs/manual/html_node/epa/Encrypting_002fdecrypting-gpg-files.html#index-epa_002dfile_002denable) somewhere in your Emacs configuration file.
- Use `C-x C-f` (`find-file`) to open your `~/.authinfo.gpg` or `~/.netrc.gpg` file. Emacs will transparently decrypt the file upon opening it, and you can then edit it just like any other text file. Please see [this](https://www.gnu.org/software/emacs/manual/html_node/emacs/Authentication.html) Emacs documentation page to learn about the syntax of this file. Emacs will encrypt the file for you whenever you save your changes.
