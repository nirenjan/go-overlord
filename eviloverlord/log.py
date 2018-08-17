"""Evil Overlord Logging Module"""

# This module helps with logging to stderr, and it is a replacement for the
# built in logging library

import sys
import inspect
import enum

# Predefined Log Levels
class Level(enum.IntEnum):
    """Log Levels"""
    DEBUG = 1
    INFO = 2
    WARNING = 3
    ERROR = 4
    CRITICAL = 5

class Log():
    """Logging class"""
    # Default Logging level
    loglevel = Level.WARNING

    @classmethod
    def write_msg(cls, level, *args):
        """Write log message to stderr"""
        if cls.loglevel <= level:
            if len(args) > 1:
                msg = args[0] % args[1:]
            else:
                msg = args[0]

            sys.stderr.write(str(level)[6:] + ': ' + str(msg) + '\n')

    @classmethod
    def set_level(cls, level):
        """Set log level"""
        if not isinstance(level, Level):
            raise IndexError("Invalid log level " + str(level))

        cls.loglevel = level

    @classmethod
    def debug(cls, *args):
        """Log Debug level message"""
        if cls.loglevel <= Level.DEBUG:
            func_caller = inspect.stack()[2] # Grand-parent
            module = func_caller[1].split('/')[-1] # Get the filename
            module = module[:-3] # Delete the trailing .py
            lineno = func_caller[2]
            del func_caller

            fmt_str = args[0]
            if not isinstance(fmt_str, str):
                fmt_str = str(fmt_str)
            fmt_str = '%s:%d - ' + fmt_str

            vargs = [fmt_str, module, lineno]
            vargs.extend(args[1:])

            cls.write_msg(Level.DEBUG, *vargs)

def set_level(level):
    """Set log level"""
    Log.set_level(level)

def debug(*args):
    """Log Debug level message"""
    Log.debug(*args)

    # if LOG_LEVEL <= Level.DEBUG:
    #     func_caller = inspect.stack()[1] # Parent
    #     module = func_caller[1].split('/')[-1] # Get the filename
    #     module = module[:-3] # Delete the trailing .py
    #     lineno = func_caller[2]
    #     del func_caller
    #
    #     fmt_str = args[0]
    #     if not isinstance(fmt_str, str):
    #         fmt_str = str(fmt_str)
    #     fmt_str = '%s:%d - ' + fmt_str
    #
    #     vargs = [fmt_str, module, lineno]
    #     vargs.extend(args[1:])
    #
    #     Log.write_msg(Level.DEBUG, *vargs)

def info(*args):
    """Log Info level message"""
    Log.write_msg(Level.INFO, *args)

def warning(*args):
    """Log Warning level message"""
    Log.write_msg(Level.WARNING, *args)

def error(*args):
    """Log Error level message"""
    Log.write_msg(Level.ERROR, *args)

def critical(*args):
    """Log Critical level message. This is also fatal"""
    Log.write_msg(Level.CRITICAL, *args)
    sys.exit(1)
