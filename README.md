# lib-katana

Shared library for Boss Katana management tasks.

- Full API in Go with C API for other languages.
- SysEx message processing and generation.
- `.tsl` patch loading and generation.
- Bulk upload strategies for fast patch changes.

## Go vs C API

There is some performance and development overhead involved when using the C API.
This means it will not be worth the effort of exposing all low-level functionality.
Such as converting types or the checksum function.
Here are the guidelines for whether a C API should expose functionality.

Functionality should be included in the C API when:
- It is non-trivial to implement.
- Compatibility is not determined by Boss/Roland specs.
- It is (part of) a higher order function exposing significant parts of the library.
- It involves data that required some research.
