# Test Data Generator

A terminal-based Go application for generating JSON test data from a user-defined schema.

## Features

- Generates `integer`, `float`, `string`, and `boolean` fields.
- Validates field names, ranges, record count, and string length.
- Saves JSON files to the `output` folder.
- Lets you choose a file name and keeps previous files by adding `_2`, `_3`, and so on.
- Opens the folder with the generated file after export on Windows.

## Run from source

Install Go, then run this command from the project folder:

```powershell
go run .
```

## Build for Windows

```powershell
go build -o testsGenerator.exe .
```

Run `testsGenerator.exe`. Generated files are stored in the `output` folder.

## Create a release ZIP

Run this command from the project folder:

```powershell
.\build-release.ps1
```

The script creates `release/test-data-generator-windows.zip`. Attach that ZIP file to a GitHub Release. A Windows user only needs to download the ZIP, extract it, and run `TestDataGenerator.exe`.

## Example output

```json
[
  {
    "age": 25,
    "active": true
  },
  {
    "age": 41,
    "active": false
  }
]
```

## Controls

- `Tab` or arrow keys: move between fields.
- `+`: add or save a field.
- `E`: edit the selected field.
- `Delete`: delete the selected field.
- `Enter`: generate and export records.
- `O`: open the output folder.

## License

MIT
