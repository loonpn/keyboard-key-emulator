# keyboard-key-emulator
A keyboard key emulator. Read data from a serial port on a Windows host, and send the keyboard key values to the host based on the configuration file.

## Usage

-B specifies the Baud, - p specifies the serial port, and - f specifies the configuration file path


## Bug processing logic

- Before reading the configuration file, check if it exists, is readable, and conforms to the ini format. If these conditions are not met, the program outputs an error message and terminates.

- When parsing key value pairs in the configuration file, check if the key or value is empty or illegal. If so, the program ignores this key value pair or outputs a warning message and continues to parse other key value pairs.

- Before opening the serial port connection, check whether the Baud and serial port are valid. If not, the program outputs an error message and terminates.

- When sending keyboard key values, check if the key values exist in the list of keys supported by the robotgo library. If not, the program outputs a warning message and skips this key value.


## Configuration file format

- The [config] section is used to set the Baud and serial port.

- The [keymap] section is used to set the correspondence between keyboard key values and serial port data.


## Example config.ini

```ini
[config]
Baud=9600
Port=COM1

[keymap]
a=hello
b=world
c=123
ctrl+c=ctrl+c
```
