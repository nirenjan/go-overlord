"""Terminal escape sequence handling"""

import shutil
import enum

class Color(enum.IntEnum):
    """Colors used by the terminal"""
    Black = 0
    Red = 1
    Green = 2
    Yellow = 3
    Blue = 4
    Magenta = 5
    Cyan = 6
    White = 7

def _csi(command):
    """Command Sequence Indicator"""
    return '\033[' + command

def foreground(color):
    """Make text foreground color"""
    return _csi("3%dm" % color)

def background(color):
    """Make text background color"""
    return _csi("4%dm" % color)

def reset():
    """Reset the terminal text"""
    return _csi('m')

def clear():
    """Clear the screen"""
    return _csi('2J') + _csi('H')

def bold():
    """Make text bold"""
    return _csi('1m')

def italic():
    """Make text italic"""
    return _csi('2m')

def underline():
    """Make text underlined"""
    return _csi('3m')

def horizontal_line(char='-'):
    """Print a horizontal line spanning the width of the terminal"""
    cols, _ = shutil.get_terminal_size()

    # Validate char
    if not isinstance(char, str) or not char:
        char = '-'
    elif len(char) > 1:
        char = char[0]

    return '-' * cols
