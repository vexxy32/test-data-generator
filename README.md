# Test Data Generator

A small Windows terminal application that generates JSON test data from a schema you enter on screen.

## Download and run

Open the [latest release](https://github.com/vexxy32/test-data-generator/releases/latest), download `test-data-generator-windows.zip`, extract it, and run `TestDataGenerator.exe`. Go is not required.

The application creates an `output` folder in its current working directory and shows the full path to each generated file. When launched from an extracted folder in Windows Explorer, this is normally the folder beside the `.exe`.

## Features

- Generates `integer`, `float`, `string`, and `boolean` fields.
- Validates field names, ranges, record count, and string length.
- Saves JSON files to the `output` folder.
- Lets you choose a file name and keeps previous files by adding `_2`, `_3`, and so on.
- Opens the folder with the generated file after export on Windows.

The program generates data records, not Go unit-test source files.

## Run from source

Install the Go version specified in `go.mod`, then run this command from the project folder:

```powershell
go run ./cmd/test-data-generator
```

## Build for Windows

```powershell
go build -o TestDataGenerator.exe ./cmd/test-data-generator
```

Run `testsGenerator.exe`. The `output` folder is created in the current working directory.

## Create a release ZIP

Run this command from the project folder:

```powershell
.\scripts\build-release.ps1
```

The script runs the tests, builds the application, and creates `release/test-data-generator-windows.zip`. When a new version is needed, update `VERSION` and push to `main`; GitHub Actions runs the same script and publishes the ZIP as a GitHub Release. A Windows user only needs to download the ZIP, extract it, and run `TestDataGenerator.exe`.

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

- `Tab`, `↑`, or `↓`: move between form fields; `↑` and `↓` also select a field in the schema list.
- `+`: add or save a field.
- `E`: edit the selected field.
- `Delete`: delete the selected field.
- `P` or `Ctrl+P`: prepare a preview of the first three records.
- `Enter`: save the prepared preview to a JSON file.
- `O`: open the output folder.
- `Esc`: exit.

## License

MIT
