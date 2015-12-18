# Что это

**K-Switcher** — утилита для переключения раскладки набранного текста, по мотивам известного Punto Switcher-а.

## Что умеет

K-Switcher умеет перекодировать выделенный текст с русской раскладки на английскую и наоборот. После перекодировки раскладка переключается в ту, в которой набран последний фрагмент текста. Перекодировка происходит при нажатии сочетания клавиш **Shift+Break**.

Поддерживаются стандартные русская и английская раскладки Windows, а также [раскладки Бирмана](http://ilyabirman.ru/projects/typography-layout/).

## Как использовать

Скачать последнюю версию K-Switcher можно со страницы https://github.com/davidmz/k-switcher/releases (есть версии для 32- и 64-битных Windows).

Программа состоит из одного исполняемого файла, инсталляция не требуется, нужно просто запустить файл. Программа представляет собой консольное приложение.

## Как это работает

При работе программа эмулирует нажатия клавиш для работы с буфером обмена. Алгоритм работы следующий:

1. сохраняется состояние буфера обмена;
2. эмулируется нажатие `Ctrl+C` (скопировать выделенный текст);
3. текст в буфере обмена перекодируется;
4. эмулируется нажатие `Ctrl+V` (вставить выделенный текст);
5. восстанавливается состояние буфера обмена.

Таким образом, перекодировка работает только в программах, использующих сочетания клавиш `Ctrl+C` и `Ctrl+V` для работы с буфером обмена. Таких программ большинство, но могут встретиться и исключения, в которых K-Switcher будет работать некорректно.

При перекодировке K-Switcher сохраняет и восстанавливает состояние буфера, однако сохраняется только текстовое содержимое. Если в буфере обмена была картинка, ссылка на файл, фрагмент документа с форматированием и т. п., то после перекодировки это содержимое будет потеряно.
