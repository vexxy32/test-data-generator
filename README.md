# Test Data Generator

A small Windows terminal application that generates JSON test data from a schema you enter on screen.

## Download and run

Open the [latest release](https://github.com/vexxy32/test-data-generator/releases/latest), download `test-data-generator-windows.zip`, extract it, and run `TestDataGenerator.exe`. Go is not required.

The application creates an `output` folder in its current working directory and shows the full path to each generated file. When launched from an extracted folder in Windows Explorer, this is normally the folder beside the `.exe`.

## Features

- Generates `integer`, `float`, `string`, and `boolean` fields.
- Rounds `float` values to a precision you choose from 0 to 6 decimal places.
- Recognizes names such as `name`, `email`, `city`, `country`, and `company`, selects the `string` type automatically, and generates values from built-in lists.
- Creates an automatic `_test_id` for every record, making a problematic input easy to find in a JSON file.
- Shows a preview of the first three records; the previewed values are the exact values saved to JSON.
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

Run `TestDataGenerator.exe`. The `output` folder is created in the current working directory.

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
    "_test_id": 1,
    "name": "Anna",
    "email": "anna.smith42@example.com",
    "age": 25
  },
  {
    "_test_id": 2,
    "name": "Max",
    "email": "max.brown913@example.com",
    "age": 41
  }
]
```

## Automatic string fields

For these field names, the application selects `string` automatically and does not ask for a length range: `name`, `surname`, `email`, `city`, `country`, `company`, `department`, `job_title`, `street_address`, `phone_number`, `user_role`, `status`, `color`, `currency`, `language`, and `content_title`. Plural and underscore variants such as `cities` and `job_titles` are also supported.

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
