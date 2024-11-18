# MoveMusic

## Description
An opinionated music file copy which can copy to a single file or a file with associated folders. Information about what the file gets named is extracted from the ID3 tag and properly formatted by the function.

## Installation
Instructions on how to install and set up your project.

```bash
# Clone the repository
git get github.com/punkscience/movemusic
```
# Example use case
```bash
useFolderNames := false

movemusic.CopyMusic( sourceFile, destFolder, useFolderNames )

```