"""Evil Overlord Init Module"""

import os
import shutil

from . import config
from . import log
from . import modules
from . import git

########################################################################
# CLI parser
########################################################################
INIT_DESC = '''
Initialize the overlord database and make it ready to use.
This must be the first command run, after that there is no
need to run init again. Running with no options specified
acts will fail if it is already initialized.
'''

def register_command(parser):
    """Register the backup command with the parent parser"""
    log.debug('Registering backup command')

    init_cmd = parser.add_parser('init', help='overlord initialization',
                                 description=INIT_DESC)

    init_cmd.add_argument('-f', '--force', action='store_true',
                          help='force reinitialization of overlord')
    init_cmd.add_argument('-w', '--wipe', action='store_true',
                          help='wipe any existing data, implies --force')

def handle_command(parser, args):
    """Handle init command line"""
    # parser is unused, delete it
    del parser

    # Check if Overlord is already initialized
    data_dir = config.data_dir()
    inited_file = os.path.join(data_dir, '.init')
    inited = os.path.exists(inited_file)
    force = args.force

    if args.wipe:
        log.debug('Deleting any existing Overlord repository')
        shutil.rmtree(data_dir, ignore_errors=True)
        inited = False
        force = True

    if inited and not force:
        log.critical('Cannot reinitialize Overlord')

    # Create the data directory, and run git initialize
    git.init()

    # Create the init marker file, we don't need any content
    with open(inited_file, 'w'):
        pass

    git.add([inited_file])
    git.commit('Evil Overlord: Surrender to the Dark Side!')

    modules.run_callback('module_init', force=force)
