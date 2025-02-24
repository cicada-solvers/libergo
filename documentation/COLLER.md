# Column Processing Program

## Description
This program processes a stack of data from a file, starting from a specified position and processing either by rows or columns. It can also alternate between directions.

## Usage

### Flags
- `-file`: The file containing the stack.
- `-position`: The initial position in the stack (`top_left`, `top_right`, `bottom_left`, `bottom_right`).
- `-processing`: The processing type (`row`, `column`).
- `-alternating`: Alternate between directions (default: `false`).

### Examples

#### Process stack from a file starting from the top left and processing by rows
```sh
./coller -file=stack.txt -position=top_left -processing=row