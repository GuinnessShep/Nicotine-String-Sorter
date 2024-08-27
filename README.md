# Nicotine-String-Sorter

A tool to sort your Url:Log:Pass strings using Go


![image](https://raw.githubusercontent.com/Underneach/Nicotine-String-Sorter/String-Sorter-regexp/image_1.png)
![image](https://raw.githubusercontent.com/Underneach/Nicotine-String-Sorter/String-Sorter-regexp/image_2.png)


## Features of the Sorter

- Extracts strings from a file or files within a folder
- Saves strings in the Log:Pass or Url:Log:Pass format
- Sorts based on site requests (e.g., google.com) or keywords within URLs (e.g., google)
- Supports multithreaded sorting and concurrent file writing while skipping duplicate strings, enabling it to handle databases of any size

## Features of the Cleaner

- Cleans databases of any size by processing strings immediately upon reading, without loading the entire list into RAM
- Cleans multiple databases either separately or merges them into one file
- Removes invalid strings (A-z / 0-9 / Special characters | 10-256 characters | UNKNOWN)
- Removes duplicates using xxh3 hashing

## Technology Stack
- **Multithreading**: [github.com/panjf2000/ants](https://github.com/panjf2000/ants)
- **Colored Output**: [github.com/fatih/color](https://github.com/fatih/color)
- **Processor Specs**: [github.com/klauspost/cpuid](https://github.com/klauspost/cpuid)
- **Available Memory Detection**: [github.com/pbnjay/memory](https://github.com/pbnjay/memory)
- **File Encoding Detection**: [github.com/saintfish/chardet](https://github.com/saintfish/chardet)
- **Progress Bar**: [github.com/schollz/progressbar](https://github.com/schollz/progressbar)
