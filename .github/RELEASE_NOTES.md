## Русский

В версии v0.5.0 добавлен блок «ФЛАГИ ПОЛЯ». Флаг «Пустое значение» позволяет намеренно добавить нестандартные данные: `""` для строковых полей и `null` для числовых и логических. Хотя бы одна запись в файле получит эти значения для всех полей с включённым флагом; выбор воспроизводится по seed.

- Флаг сохраняется и загружается вместе с шаблоном.
- Предпросмотр первых трёх записей снова точно совпадает с началом полного файла даже при включённом флаге.
- JSON снова форматируется с отступами и строками для каждого поля; потоковая запись сохранена.
- Связи внутри программы для экспорта и хранения шаблонов оформлены через небольшие интерфейсы без изменения пользовательских команд.

Скачай ZIP для Windows, распакуй и запусти `TestDataGenerator.exe`. Новые JSON-файлы будут форматированными; уже созданные файлы не меняются.

## English

Version v0.5.0 adds a **FIELD FLAGS** panel. The **Empty value** flag deliberately inserts edge cases: `""` for strings and `null` for numeric and boolean fields. At least one record in each file gets these values for every enabled field; the choice is reproducible with the seed.

- The flag is stored in and restored from templates.
- The first three previewed records match the beginning of the complete file even with the flag enabled.
- JSON records are indented across multiple lines again while export remains streaming.
- Export and template storage are now connected through small interfaces, without changing the user controls.

Download the Windows ZIP, extract it, and run `TestDataGenerator.exe`. Newly generated JSON files use the readable format; existing files are unchanged.
