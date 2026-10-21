## Русский

В версии v0.4.0 генерация и сохранение JSON выполняются по одной записи. Программа больше не держит весь набор данных в памяти: это позволяет создавать гораздо более крупные файлы. Максимальное количество записей в интерфейсе увеличено до 100 000 000.

- Предпросмотр показывает первые три записи. При сохранении те же схема и seed воспроизводят их в начале полного файла.
- Список полей прокручивается стрелками и показывает, какая часть схемы сейчас видна.
- Существующие файлы не перезаписываются: к имени нового файла добавляется номер.

Скачай ZIP для Windows, распакуй и запусти `TestDataGenerator.exe`. Длительность создания больших файлов зависит от количества полей и скорости диска; во время генерации интерфейс ожидает завершения записи.

## English

Version v0.4.0 generates and writes JSON one record at a time. The application no longer keeps the complete data set in memory, enabling much larger files. The UI limit is now 100,000,000 records.

- The preview shows the first three records; saving reproduces them at the start of the complete file using the same schema and seed.
- The field list scrolls with the arrow keys and shows the visible range.
- Existing output files are preserved by adding a number to new filenames.

Download the Windows ZIP, extract it, and run `TestDataGenerator.exe`. Large exports depend on field count and disk speed; the interface waits while the file is being written.
