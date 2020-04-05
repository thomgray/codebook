# codebook

Command line application for traversing and searching markdown formatted documents (for now)

## Install

```
go get github/thomgray/codebook
```

## Usage

Run the application

```
codebook
```

There are 3 input modes:
* `traversal` (`>`)
* `search` (`?`)
* `command` (`:`)

### To get started

To enter either mode, type the corresponding symbol at the beginning of the command prompt (the symbol at the head of the prompt should indicate which 'mode' you are in).

To get started, you will need to configure which your `search directories`. This will tell the application where to find notes.

In `command` mode, specify directories:

```
: + <absolute path to directory>
```

List your search paths with the following command:

```
: ls
```


### Traversing documents

Provided there are some markdown documents in your search director(y/ies), you can now traverse these documents in `traversal` mode

To find the document `hello`:
```
> hello
```

__NOTE__: a document gets it's 'name' from an `h1`/`#` heading. If a document doesn't start with an h1 heading, it will take the file name as the title, and all headings within the document are considered sub-documents.

You can traverse a document tree by inputting sub sections in a single traversal command.

Tab-completion is enabled for document traversal.

#### Traversal scope

You can scope your traversal with several special characters:

* `*` means 'from the external context', i.e. the context the application starts in when no document is currently selected. Traversing in this context starts at the top level of all documents in your search paths(s)
* `/` means 'from the root of this document', i.e. traverse from the top level of the currently selected document.
* `.` means 'from the current context'. i.e. traverse from the current location in this document.

e.g.
```
> * top level document
> / from the root of here
> . sub from here
```

Note that the default is a combination of 'current' context (i.e. `.`) falling back to `external` context (`*` if there are no matches for the command in current context).

### Searching documents

TODO

### Commands

If you ever need help:
```
: help
```

Will output available command in command mode, including documentation for general usage.
