# Test Data Generator

[English](README.md) | [Русский](README.ru.md)

A Windows terminal application that generates JSON test data from a schema you enter on screen. The source can also be built with Go.

## Download and run

Open the [latest release](https://github.com/vexxy32/test-data-generator/releases/latest), download `test-data-generator-windows.zip`, extract it, and run `TestDataGenerator.exe`. Go is not required.

The application creates an `output` folder in its current working directory and shows the full path to each generated file. When launched from an extracted folder in Windows Explorer, this is normally the folder beside the `.exe`.

## Quick start

1. Enter a field name and type, then press `+` to add it to the schema. Names such as `name` and `cities` select `string` automatically.
2. Set the number of records and the output filename. Leave Seed empty for a new random seed, or enter one to reproduce a result.
3. Press `Ctrl+P` to preview up to three records, then `Enter` to write the complete JSON file. If you press `Enter` before preparing a preview, the first press prepares the preview and the next saves the file.

The application shows the generated seed and opens the output folder after saving.

## Features

- Generates `integer`, `float`, `string`, and `boolean` fields.
- Lets you set numeric ranges, ordinary string length ranges, and `float` precision from 0 to 6 decimal places. Empty bounds default to `0..100` for numbers and `1..20` for ordinary strings.
- Recognizes names such as `name`, `email`, `city`, `country`, and `company`, selects the `string` type automatically, and generates values from built-in lists.
- Generates `name` and `surname` together with an `email` based on the same pair within each record.
- Offers a per-field **Empty value** flag. In at least one automatically selected record, enabled string fields become `""`, while enabled integer, float, and boolean fields become `null`. All enabled fields use the same selected record; it may be outside the first three shown in the preview.
- Creates an automatic `_test_id` for every record, making a problematic input easy to find in a JSON file.
- Supports a generation seed: reuse the same seed, schema, field order, and record count to reproduce the same data.
- Shows a preview of the first three records; the previewed values are the exact values saved to JSON.
- Saves and loads named JSON templates containing the schema, record count, output filename, and field flags. The seed is not saved.
- Generates and writes pretty-printed JSON one record at a time, without keeping the complete data set in memory. The current UI limit is 100,000,000 records; large exports can take time and disk space, and the interface waits until writing finishes.
- Validates field names, ranges, record count, and string length.
- Saves JSON files to the `output` folder in the current working directory.
- Lets you choose a file name and keeps previous files by adding `_2`, `_3`, and so on.
- Opens the folder with the generated file after export on Windows. `Ctrl+N` starts a new data set without restarting.

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

## Field settings and automatic string fields

For these field names, the application selects `string` automatically and does not ask for a length range: `name`, `surname`, `email`, `city`, `country`, `company`, `department`, `job_title`, `street_address`, `phone_number`, `user_role`, `status`, `color`, `currency`, `language`, and `content_title`. Plural and underscore variants such as `cities` and `job_titles` are also supported.

Other string field names produce random text within the chosen length range. Boolean fields generate `true` or `false`. The **FIELD FLAGS** panel belongs to the field being added or edited: select **Empty value**, press `Space` or `Enter` to switch from off to on, then press `+` to save the field. The Russian UI displays these states as `[нет]` and `[да]` and marks enabled fields with `[пустые]` in the schema list.

The guaranteed empty values are deliberate edge cases. `null` is different from zero or `false`; `""` is different from a missing JSON key. All fields remain present in the record.

## Templates

Enter a template name and press `Ctrl+S` to save it under `templates/<name>.json` in the current working directory. `Ctrl+L` opens the list; choose with the arrow keys and press `Enter` to load. Saving the same template name replaces that template. Templates store the field definitions and flags, count, and output filename. They do not store generated records or the seed. Template validation currently limits the stored count to 10,000 even though direct generation allows more.

## Controls

- `Tab`/`Shift+Tab`, `↑`, or `↓`: move between form fields; `↑` and `↓` also select a field in the schema list.
- `Space` or `Enter` on **Empty value**: toggle the flag for the field currently being edited. `Enter` elsewhere handles preview and saving.
- `+`: add or save a field.
- `E` on the schema list: edit the selected field.
- `Delete` on the schema list: delete the selected field.
- `Ctrl+P` (or `P` on the schema list): prepare a preview of the first three records.
- `Enter`: first prepare a preview if needed, then generate and save the complete JSON file on the next press.
- `Seed`: leave it empty to create a new random seed; keep or enter a seed to reproduce a generated data set.
- `Ctrl+S`/`Ctrl+L`: save a template or choose one to load. `S`/`L` also work on the schema list.
- `Ctrl+O` (or `O` on the schema list): open the output folder after saving.
- `Ctrl+N`: start a new test without restarting the application.
- `Esc`: cancel field editing, close the template picker, or exit.

## Feedback

Found a bug or have an idea? Open an [issue](../../issues), visit [@vexxy32](https://github.com/vexxy32), or email [egorkuznecov254@gmail.com](mailto:egorkuznecov254@gmail.com).

## License

MIT
