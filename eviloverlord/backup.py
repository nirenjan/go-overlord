"""Evil Overlord Backup Module"""

import json
import lzma
import sys
import time

from . import log
from . import modules

########################################################################
# CLI handlers
########################################################################
def export_handler(output_file):
    """Export data to the given output file"""
    backup_dict = {}
    display_info = True

    # Run export data in all functions
    modules.run_callback('export_data', backup_dict)

    if output_file == '-':
        output = sys.stdout
        display_info = False
    else:
        output = open(output_file, 'wb')

    output.write(lzma.compress(json.dumps(backup_dict).encode('utf-8'),
                               check=lzma.CHECK_SHA256))
    output.close()

    if display_info:
        print('Data exported to', output_file)


def import_handler(input_filename):
    """Import overlord data from an input file"""
    if input_filename == '-':
        input_filename = 'stdin'
        input_file = sys.stdin.buffer
    else:
        input_file = open(input_filename, 'rb')

    data = input_file.read()
    decompressed_data = lzma.decompress(data)
    input_file.close()

    backup_dict = json.loads(decompressed_data.decode('utf-8'))

    for module, data in backup_dict.items():
        handler = modules.find_module_function(modules.find_module(module),
                                               'import_data')

        if handler is None:
            log.warning('Unsupported data %s for import', module)
        else:
            handler(data)

    print('Finished import from', input_filename)

########################################################################
# CLI parser
########################################################################
BACKUP_DESC = '''
The Backup module allows you to backup your overlord activity and import
it on a new machine.
'''

def register_command(parser):
    """Register the backup command with the parent parser"""
    log.debug('Registering backup command')

    backup_cmd = parser.add_parser('backup', help='overlord backups',
                                   description=BACKUP_DESC)

    subcommand = backup_cmd.add_subparsers(title='Subcommands', dest='subcmd',
                                           metavar='')

    filename = time.strftime('overlord-backup-%Y-%m-%d-%H-%M', time.localtime())

    export_cmd = subcommand.add_parser('export',
                                       help='export all overlord data')
    export_cmd.add_argument('file', nargs='?', default=filename,
                            help='filename to save the exported data ' +
                            '(- for stdout)')

    import_cmd = subcommand.add_parser('import',
                                       help='import overlord data')
    import_cmd.add_argument('file', nargs='?', default='-',
                            help='filename to import data from (- for stdin)')

def handle_command(parser, args):
    """Handle backup command line"""
    if args.subcmd is None:
        parser.print_help()
    elif args.subcmd == 'export':
        log.debug('backup export ' + args.file)
        export_handler(args.file)
    elif args.subcmd == 'import':
        log.debug('backup import ' + args.file)
        import_handler(args.file)
    else:
        log.error('not implemented - backup %s', args.subcmd)
