"""Evil Overlord Git Support"""

import subprocess
import inspect
import os
import contextlib
import itertools
import socket
import time

from . import config
from . import log

def _module_text():
    """Return the name of the module calling git"""
    frame = inspect.stack()[2] # Use the grandparent
    module = frame[1].split('/')[-1] # Get the basename of the module
    module = module[:-3] # Remove the trailing .py
    del frame

    return module


def ignore(patterns=None, reset=False):
    """Append the ignore_pattern to the .gitignore file"""
    module = _module_text()
    path = config.data_dir()

    if not patterns and not reset:
        # Don't do anything if no pattern is specified and no need to reset
        return

    if patterns is not None and not isinstance(patterns, list):
        raise AttributeError('expected type list for patterns, got %s' %
                             type(patterns))

    moduledir = os.path.join(path, module)
    gitignore = os.path.join(moduledir, '.gitignore')
    if not os.path.exists(gitignore) or reset:
        log.debug('making directory %s', moduledir)
        os.makedirs(moduledir, exist_ok=True)

        log.debug('creating file %s', gitignore)
        with open(gitignore, 'w'):
            # Do nothing, it will be modified later
            pass

    # Write pattern to file
    with open(gitignore, 'a') as ign:
        for pattern in patterns:
            ign.write(pattern + '\n')

    # Add the .gitignore file to the index
    add([gitignore])


@contextlib.contextmanager
def chdir(path):
    """Use a context manager to run within a path"""
    old_dir = os.getcwd()
    try:
        os.chdir(path)
        yield
    finally:
        os.chdir(old_dir)

def add(files):
    """Add files to the git database"""
    if not isinstance(files, list):
        raise AttributeError('expected %s for files, got %s' %
                             (list, type(files)))

    with chdir(config.data_dir()):
        args = ['git', 'add']
        args.extend(files)

        subprocess.call(args,
                        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def delete(files):
    """Delete files from the git database"""
    if not isinstance(files, list):
        raise AttributeError('expected %s for files, got %s' %
                             (list, type(files)))

    with chdir(config.data_dir()):
        args = ['git', 'rm']
        args.extend(files)

        subprocess.call(args,
                        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def commit(message, author_date=None):
    """Commit the contents, optionally setting the authored date"""
    # Reset all Git environment variables
    git_actioners = ['GIT_AUTHOR_', 'GIT_COMMITTER_']
    git_attributes = ['NAME', 'EMAIL', 'DATE']
    for elem in itertools.product(git_actioners, git_attributes):
        try:
            del os.environ[''.join(elem)]
        except KeyError:
            pass

    if author_date is not None:
        os.environ['GIT_AUTHOR_DATE'] = time.strftime('%Y-%m-%dT%H:%M:%S%z',
                                                      time.localtime(
                                                          int(author_date)))

    # Force the committer values to Overlord
    os.environ['GIT_COMMITTER_NAME'] = 'Evil Overlord'
    os.environ['GIT_COMMITTER_EMAIL'] = 'eviloverlord@' + socket.gethostname()

    with chdir(config.data_dir()):
        message = _module_text() + ': ' + message
        subprocess.call(['git', 'commit', '-m', message],
                        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def init():
    """Initialize git in the data directory"""
    with chdir(config.data_dir()):
        subprocess.call(['git', 'init', '.'],
                        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
