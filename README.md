# lib-katana

Shared library for Boss Katana management tasks.

- Full API in Go with C API for other languages.
- SysEx message processing and generation.
- `.tsl` patch loading and generation.
- Bulk upload strategies for fast patch changes.

## Roadmap

**Bleeding edge `v0.x.x`**

Go is a new language for me personally, is a recently created language and not a lot of people are
using it for shared libraries at the moment.
So odds are that I'm going to get thing wrong the first time around.

In the `0.x.x` versions the focus is on fleshing out the library and not sweating too much about
the details of the library architecture or the ideal APIs.
It's a best effort to make something fit for purpose quickly so we can find what the problems are.

In this stage I won't bother with libtool versions. Everything will break, all the time.
During this period you should commit the library in your source control and only upgrade if you have
some time to implement the API changes.

**Building v1**

Once the bleeding edge stuff has seen some adoption, the features have fleshed out for the most part,
and the main architecture and API problems are clear, it's time for v1.
We'll put some effort into fixing the major design issue and refactor to apply these lessons.
Along with that we'll start doing more strict libtool / semver release management and change logs.
If time allows it I would like to add in good tests and CI here too so the API can stabilize.

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
