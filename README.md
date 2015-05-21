fidu
====

A 2D fiducial marker generator for MultiTaction displays.

The purpose of this tool is to replace the MarkerFactory application that is shipped
with the Cornerstone SDK and depends on Qt to generate the fiducial markers.

`fidu` is capable of generating the same fiducial markers using only the standard libraries
provided by Go.

## Usage

Generate the default fiducial marker using the command:
```
fidu
```
The result is a 288px by 288px PNG image with the representation of the marker.

The user can override the default options, using the same flags that are present in MarkerFactory. For example:
```
fidu --code 20 --division 3 --blocksize 16 --filename code.png
```

There is an additional `--size` flag which can be used to define the size, in pixels, of each side of
the resulting PNG image. However, `--size` and `--division` must comply in order to get an integer value from the
expression `size = blocksize * (division + 4)`.
