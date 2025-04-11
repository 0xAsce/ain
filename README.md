# ain

Append lines into a file as long as the lines don't already exist

## Usage 

a file called `list.txt` contains a list of names. `newlist.txt` contains a second
list of names, some of which appear in `list.txt` and some of which do not. `ain` is used
to append the new names in `newlist.txt` to `list.txt`.


```
▶ cat list.txt
Mark
Jake
David

▶ cat newlist.txt
James
David
Luke
Mark

▶ cat newlist.txt | ain list.txt
James
Luke

▶ cat list.txt
Mark
Jake
David
James
Luke

```

Note that the new lines added to `list.txt` are also sent to `stdout`

```
▶ cat newlist.txt | ain list.txt -d
Luke
James
```

## Flags

- To view the new lines in stdout, but not append to the file, use the dry-run option `-d`.
- To view the repeated lines in stdout, but not append to the file, use the wet-run option `-s`.
- To view the repeated lines in stdout, and append them to the file, use the repeat option `-r`.
- To append to the file, but not print anything to stdout, use quiet mode `-q`.

## Install

You can either install using go:

```
go install -v github.com/0xAsce/ain@latest
```

Or download a [binary release](https://github.com/0xAsce/ain/releases).
