"""Commandline implementation"""
import sys
import argparse

from . import log
from . import modules

def register_cli_modules(*args, **kwargs):
    """Register CLI handlers"""
    modules.run_callback('register_command', *args, **kwargs)

def handle_commandline(parser, args):
    """Handle Overlord command line"""
    # Handle verbosity
    log.debug('Verbosity set to %s', args.verbose)
    if args.verbose is None:
        loglevel = log.Level.WARNING
    elif args.verbose == 1:
        loglevel = log.Level.INFO
    else: # Anything greater than 2
        loglevel = log.Level.DEBUG
    log.debug('Set log level to %s', loglevel)
    log.set_level(loglevel)

    if args.command is None:
        parser.print_help()
        sys.exit(1)

    module = sys.modules.get('eviloverlord.' + args.command, None)
    if module is None:
        log.error('unsupported command %s', args.command)
        parser.print_help()
        sys.exit(1)

    handler = getattr(module, 'handle_command', None)
    if handler is None:
        log.error('not implemented - %s', args.command)
        sys.exit(1)

    # Get the subparser for the individual commands
    subparsers = [
        action for action in parser._actions
        if isinstance(action, argparse._SubParsersAction)]

    try:
        subparser = subparsers[0].choices[args.command]
    except KeyError:
        log.error('unable to find subcommand %s', args.command)
        sys.exit(1)

    # Call the handler
    log.debug('Calling %s.%s', module.__name__, handler.__name__)
    handler(subparser, args)

OVERLORD_DESCRIPTION = '''
Overlord is a command-line based personal assistant. It can take notes,
keep a journal, make reminders and more.
'''

def generate_commandline(version):
    """Generate the Evil Overlord top-level command line"""
    parser = argparse.ArgumentParser(prog='overlord',
                                     description=OVERLORD_DESCRIPTION,
                                     add_help=True)

    parser.add_argument('-V', '--version',
                        action='version', version='Evil Overlord %s' % version,
                        help='display the version and exit')
    parser.add_argument('-v', '--verbose', action='count',
                        help='set verbosity level')

    commands = parser.add_subparsers(title='Commands', dest='command',
                                     metavar='')

    # Register sub-commands
    register_cli_modules(commands)

    return parser
