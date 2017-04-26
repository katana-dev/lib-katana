# lib-katana

C library for Boss Katana management tasks.

- SysEx message processing and generation.
- `.tsl` patch loading and generation.
- Bulk upload strategies for fast patch changes.

## Building from source

Requires

- `gcc` or other C compiler
- `automake`
- `libtool`

Building

```sh
autoreconf --install
./configure
make
```

Running tests

```
make check
```

Installing

```
sudo make install
```
