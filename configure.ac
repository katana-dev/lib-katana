# Package info (semantic versioning)
AC_INIT([lib-katana], [0.0.0], [https://github.com/katana-dev/lib-katana/issues])

# Library version (libtool versioning)
LIBKATANAVERSION=0:0:0
AC_SUBST(LIBKATANAVERSION)

# Set autotools parameters.
AC_CONFIG_AUX_DIR(autotools)
AC_CONFIG_MACRO_DIR(m4)
AM_CONFIG_HEADER(config.h)
AM_INIT_AUTOMAKE([-Wall -Werror foreign])

# Declare build dependencies
AC_PROG_CC
AM_PROG_AR

# We want libtool for the library.
LT_INIT

# What we want to generate.
AC_CONFIG_FILES([
	Makefile
	src/Makefile
])

# Do the work.
AC_OUTPUT
