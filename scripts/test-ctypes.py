import os
from ctypes import *

# Route to the library.
script_dir = os.path.dirname(os.path.abspath(__file__))
so_path = os.path.join(script_dir, "../src/.libs/libkatana.so")

# Load and run with ctypes.
libkatana = CDLL(so_path)
libkatana.main()
