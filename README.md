# Nicotine-String-Sorter
Sorter for your Url:Log:Pass strings in Go


![image](https://raw.githubusercontent.com/Underneach/Nicotine-String-Sorter/String-Sorter-regexp/image_1.png)
![image](https://raw.githubusercontent.com/Underneach/Nicotine-String-Sorter/String-Sorter-regexp/image_2.png)


## What the sorter can do

    Getting strings from a file or files in a folder
    Saving as Log:Pass or Url:Log:Pass
    Sorting by site request (google.com) or keyword in the link (google)
    Multithreaded sorting and simultaneous writing to files while skipping duplicate strings - reading a database of any size

## What the cleaner can do

    Cleaning a database of any size - strings are processed immediately upon reading, without loading the list into RAM
    Cleaning multiple databases individually or all databases into one file
    Removing invalid strings (A-z / 0-9 / Special characters | 10-256 characters | UNKNOWN)
    Duplicate removal implemented via xxh3 hash



## Stack
+  Multithreading - github.com/panjf2000/ants
+  Colored output - github.com/fatih/color
+  Processor specs - github.com/klauspost/cpuid
+  Getting available memory - github.com/pbnjay/memory
+  File encoding detection - github.com/saintfish/chardet
+  Progress bar - github.com/schollz/progressbar
