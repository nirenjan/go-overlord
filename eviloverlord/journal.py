"""Journal Logging"""

# For pylint processing
from __future__ import print_function

import argparse
import collections
import errno
import fnmatch
import hashlib
import json
import math
import os
import subprocess
import time

from . import terminal
from . import config
from . import git
from . import log

########################################################################
# Journal handling
########################################################################
class JournalError(Exception):
    """Exception raised during Journal processing"""
    pass

DBEntry = collections.namedtuple('DBEntry',
                                 ['entry_id', 'date', 'path', 'tags', 'title'])

class Journal:
    """Journal module - operations are at the class level"""
    _data = os.path.join(config.data_dir(), 'journal')
    _db = os.path.join(_data, '.db')
    _checksum = os.path.join(_data, '.db_checksum')
    _version = os.path.join(_data, '.db_version')

    @classmethod
    def path(cls, date):
        """Return the path to the journal entry for the given date/time"""
        return os.path.join(cls._data,
                            '%04d' % date.tm_year,
                            '%02d' % date.tm_mon,
                            '%02d' % date.tm_mday,
                            '%02d%02d.log' % (date.tm_hour, date.tm_min))

    @staticmethod
    def db_entry(entry):
        """Compute the database entry line from the Entry object"""
        return '%s:%s:%s:%s\n' % (entry.entry_id, entry.path,
                                  ' '.join(entry.tags), entry.title)

    @classmethod
    def module_init(cls):
        """Initialize the journal module"""
        # Create the data directory
        os.makedirs(cls._data, exist_ok=True)

        # Create the database file, no issues if it already exists
        with open(cls._db, 'a'):
            pass

        cls.update_checksum()

        git.ignore(reset=True,
                   patterns=[
                        '# Ignore all hidden files, except .gitignore',
                        '.*',
                        '!.gitignore'])

        git.commit('init')

    @classmethod
    def regenerate_database(cls):
        """Regenerate the Journal database"""
        log.info('Regenerating journal database')
        # Walk the journal directory for entries, and generate the database
        files = []

        # recursive argument for glob.glob is only available from
        # Python 3.5 onwards, this has to work with older versions
        for root, _, filenames in os.walk(cls._data):
            for filename in fnmatch.filter(filenames, '*.log'):
                files.append(os.path.join(root, filename))

        # Generate the database
        with open(cls._db, 'w') as database:
            for efile in sorted(files):
                entry = Entry.from_path(efile)
                db_line = cls.db_entry(entry)
                database.write(db_line)

        cls.update_checksum()


    @classmethod
    def add_entry(cls, entry):
        """Add an Entry object to the database"""
        with open(cls._db, 'a') as database:
            database.write(cls.db_entry(entry))

        # Update the database checksum
        cls.update_checksum()

    @classmethod
    def delete_entry(cls, entry):
        """Delete an entry object from the database"""
        entries = cls().entries
        entries = [Entry.from_path(db_entry.path) for db_entry in entries
                   if db_entry.entry_id != entry.entry_id]

        with open(cls._db, 'w') as database:
            for journal_entry in entries:
                database.write(cls.db_entry(journal_entry))

        # Update the database checksum
        cls.update_checksum()


    @classmethod
    def calculate_checksum(cls):
        """Calculate the checksum after saving the database"""
        checksum = hashlib.sha1()
        with open(cls._db, 'r') as database:
            for db_line in database:
                checksum.update(db_line.encode('utf-8'))

        return checksum.hexdigest()

    @classmethod
    def update_checksum(cls):
        """Update the checksum after saving the database"""
        with open(cls._checksum, 'w') as csum:
            csum.write(cls.calculate_checksum() + '\n')


    @classmethod
    def verify_database(cls):
        """Verify the checksum matches the database"""
        if not os.path.exists(cls._db):
            log.warning('Database missing')
            return False

        if not os.path.exists(cls._checksum):
            log.warning('Database checksum missing')
            return False

        with open(cls._checksum, 'r') as csum:
            saved_checksum = csum.read().rstrip()

        calculated_checksum = cls.calculate_checksum()

        result = (calculated_checksum == saved_checksum)
        if not result:
            log.warning('Database checksum mismatch')

        return result

    @classmethod
    def verify_or_regenerate_database(cls):
        """Verify the database is sane, otherwise, regenerate it"""
        if not cls.verify_database():
            cls.regenerate_database()

    def __init__(self):
        """Initialize the Journal database - read from the file"""
        self.entries = []

        self.verify_or_regenerate_database()

        with open(self._db, 'r') as database:
            for entryline in database:
                entry_id, path, tags, title = entryline.rstrip('\n').split(':')
                # Get the date from the path
                date = '-'.join(path.split(os.path.sep)[-4:-1])

                entry = DBEntry(entry_id, date, path, tags.split(), title)
                self.entries.append(entry)


    def filter(self, tags=None):
        """Filter journal entries based on the tags provided"""
        tags = tags or []
        if not isinstance(tags, list):
            raise AttributeError('expected %s for tags, got %s' %
                                 (list, type(tags)))

        entries = self.entries
        if tags:
            # Filter based on tags
            entries = []
            for entry in self.entries:
                if any([tag for tag in entry.tags if tag in tags]):
                    entries.append(entry)

        return entries

    ###################################################################
    # Backup and Restore functionality
    ###################################################################
    @classmethod
    def from_dict(cls, journal):
        """Import the Journal from the given JSON dict object"""
        if not isinstance(journal, list):
            log.critical('Import expected %s, got %s', list, type(journal))

        entries = [Entry.from_dict(entry) for entry in journal]

        commit_msg = 'import-journal-entries\n\n'

        for entry in entries:
            commit_msg += '- ' + entry.title + '\n'
            entry.save()
            git.add([entry.path])

        git.commit(commit_msg)

        cls.regenerate_database()

        return True

    def to_dict(self):
        """Return a dict containing all Journal entries"""
        entries = [Entry.from_path(entry.path).to_dict()
                   for entry in self.entries]
        return {'journal': entries}

    def to_json(self):
        """Return a JSON string containing all Journal entries"""
        return json.dumps(self.to_dict())

class Entry:
    """Entry handling"""
    _IDENT_ID = '@ID\t'
    _IDENT_DATE = '@Date\t'
    _IDENT_TITLE = '@Title\t'
    _IDENT_TAGS = '@Tags\t'

    def __init__(self, title=None, timestamp=None, tags=None, text=''):
        self.text = text

        if timestamp is None:
            timestamp = time.time()
        self.timestamp = timestamp

        if tags is None:
            tags = []
        self.tags = tags

        if title is None:
            title = text.split('\n', 1)[0]
        self.title = title

        # This will be overwritten later
        self.fulltext = self.file_text()
        self.entry_id = self.compute_entry_id()

        self.path = Journal.path(time.localtime(timestamp))

    def display(self):
        """Display the journal entry"""
        def stardate(timestamp):
            """Compute the Stardate, given the timestamp"""
            jd = (timestamp / 86400.0 + 40587.5)
            return ' (Stardate ' + ('%05.9f' % jd)[:-7] + ')'

        output = ''
        # Display the date
        output += terminal.reset() + \
                  terminal.foreground(terminal.Color.Yellow) + \
                  time.strftime('%a %b %d %H:%M:%S %Z %Y',
                                time.localtime(self.timestamp)) + \
                  terminal.reset() + \
                  terminal.foreground(terminal.Color.Red) + \
                  stardate(self.timestamp) + terminal.reset() + '\n'

        # Display the title
        output += terminal.bold() + \
                  self.title + \
                  terminal.reset() + '\n'
        output += terminal.bold() + \
                  '=' * len(self.title) + \
                  terminal.reset() + '\n'

        # Display the text
        output += self.text.rstrip() + '\n' + '\n'

        # Display tags, if any
        if self.tags:
            output += terminal.bold() + \
                      'Tags:\t' + \
                      terminal.reset() + \
                      terminal.foreground(terminal.Color.Cyan) + \
                      ' '.join(self.tags) + \
                      terminal.reset() + '\n'

        # Display a horizontal line to separate the entries
        output += terminal.horizontal_line() + '\n'

        return output

    @staticmethod
    def strip_ident(ident, text):
        """Strip an identifier tag and return the rest of the line"""
        if text.startswith(ident):
            text = text[len(ident):]

        return text.rstrip('\n')

    def file_text(self):
        """Compute the file text contents"""
        text = self.title + '\n' + self.text + '\n'
        text += self._IDENT_DATE + time.strftime('%Y-%m-%dT%H:%M:%S%z\n',
                                                 time.localtime(self.timestamp))
        if self.tags:
            text += self._IDENT_TAGS + ' '.join(self.tags) + '\n'

        text += self._IDENT_TITLE + self.title + '\n'

        return text

    def save(self):
        """Save the entry to disk"""
        # Make the directory for the journal entry
        os.makedirs(os.path.dirname(self.path), exist_ok=True)

        # Write the entry to disk
        with open(self.path, 'w') as entry:
            entry.write(self.file_text())

    def compute_entry_id(self):
        """Compute and save the entry ID for the given entry"""
        return hashlib.sha1(self.fulltext.encode('utf-8')).hexdigest()[:10]

    @staticmethod
    def split_text(fulltext):
        """Get the title and text from the full text"""
        title, *text = fulltext.split('\n', 1)

        # Remove trailing whitespace from the text lines
        text = ''.join([line.rstrip() for line in text])

        return title, text

    @classmethod
    def parse_text(cls, entry):
        """Parse the file at the given path and return the Entry fields"""
        text = ''
        timestamp = None
        title = None
        tags = []
        saved_id = None

        computed_id = hashlib.md5()

        for line in entry:
            if not line.startswith('@'):
                text = text + line
            elif line.startswith(cls._IDENT_DATE):
                timestamp = math.floor(time.mktime(time.strptime(
                    cls.strip_ident(cls._IDENT_DATE, line),
                    '%Y-%m-%dT%H:%M:%S%z')))
            elif line.startswith(cls._IDENT_TITLE):
                title = cls.strip_ident(cls._IDENT_TITLE, line)
            elif line.startswith(cls._IDENT_TAGS):
                tags = cls.strip_ident(cls._IDENT_TAGS, line).split()
            elif line.startswith(cls._IDENT_ID):
                # This is an old legacy field, no longer needed with the
                # new method
                saved_id = cls.strip_ident(cls._IDENT_ID, line)
                log.debug('Saved ID is %s', saved_id)

                computed_id_hex = computed_id.hexdigest()[:10]
                if computed_id.hexdigest()[:10] != saved_id:
                    log.warning('Possibly corrupted entry')
                    log.warning('Expected id %s, got %s', saved_id,
                                computed_id_hex)
            else:
                log.critical('Unrecognized tag %s', line.split()[0])

            # Update the computed ID
            computed_id.update(line.encode('utf-8'))

        # Compute the text title
        texttitle, text = cls.split_text(text)
        if title is not None and texttitle != title:
            log.warning('Entry may be corrupted. Mismatched titles')
            log.warning('Expected "%s", got "%s"', title, texttitle)
        title = texttitle

        # Make sure that there is both title and some text
        if not title:
            log.critical('Cannot write an empty journal entry')

        if not text:
            log.critical('Journal entry must have some text')

        return title, text, timestamp, tags, saved_id


    @classmethod
    def from_path(cls, path):
        """Load entry given the path"""
        if not os.path.exists(path):
            log.critical('Unable to load %s: %s', path,
                         errno.errorcode[errno.ENOENT])

        with open(path) as entryfile:
            title, text, timestamp, tags, _ = cls.parse_text(entryfile)
        return cls(title=title, tags=tags, timestamp=timestamp, text=text)

    ###################################################################
    # Backup and Restore functionality
    ###################################################################
    @classmethod
    def from_dict(cls, entrydict):
        """Create an Entry object from the entrydict"""
        try:
            title = entrydict['title']
        except KeyError:
            log.critical('import: missing required title')

        try:
            timestamp = entrydict['timestamp']
        except KeyError:
            log.critical('import: missing required timestamp')

        try:
            tags = entrydict['tags']
        except KeyError:
            log.critical('import: missing required tags')

        try:
            text = entrydict['text']
        except KeyError:
            log.critical('import: missing required text')

        return cls(title=title, tags=tags, timestamp=timestamp, text=text)

    def to_dict(self):
        """Dump the Entry object as a dict"""
        entry = collections.OrderedDict()

        entry['timestamp'] = self.timestamp
        entry['title'] = self.title
        entry['text'] = self.text
        entry['tags'] = self.tags

        return entry

########################################################################
# CLI functions
########################################################################
def editor(filename):
    """Get the editor for the system"""
    # Use $EDITOR in preference to git config core.editor in preference to vim
    _editor = os.environ.get('EDITOR', None)

    if _editor is None:
        try:
            _editor = subprocess.check_output(['git', 'config', 'core.editor'])
        except subprocess.CalledProcessError:
            pass

    if _editor is None:
        _editor = 'vim'
    else:
        _editor = _editor.decode('utf-8').rstrip()

    # Ideally, we'd use the subprocess module to communicate with the editor
    # However, I haven't figured out how to get it to work with a proper TTY
    # For now, I'm using os.system, although it's not the best approach.
    cmd = _editor + ' ' + filename
    log.debug('Calling editor - "%s"', cmd)
    os.system(cmd)


DEFAULT_NEW_ENTRY = """
# Enter your journal message here. Lines beginning with # are deleted
# from the journal
"""
def create_entry(tags=None):
    """Create a new journal entry"""
    with git.chdir(os.path.join(config.data_dir(), '.git')):
        with open('COMMIT_EDITMSG', 'w') as entry:
            entry.write(DEFAULT_NEW_ENTRY.lstrip())

        editor('COMMIT_EDITMSG')

        entrytext = []
        tags = tags or []
        timestamp = math.floor(time.time())
        with (open('COMMIT_EDITMSG', 'r')) as entry:
            # Delete all comment lines
            entrytext = [line for line in entry
                         if not line.startswith('#')]

        title, text, *_ = Entry.parse_text(entrytext)

        # Get the actual Entry object
        entry = Entry(title=title, text=text, timestamp=timestamp, tags=tags)

        # Write the journal entry to disk
        entry.save()

        # Append the entry to the journal database
        Journal.add_entry(entry)

        # Update the git log
        git.add([entry.path])
        git.commit("add-entry '%s'" % entry.title, author_date=str(timestamp))


def operate_on_entry(entry_id, show=False, delete=False):
    """Operate on the entry given the ID"""
    journal = Journal()
    entry_matches = [entry for entry in journal.entries
                     if entry.entry_id == entry_id]

    if not entry_matches:
        log.critical('Unable to find entry with ID %s', entry_id)

    old_entry = Entry.from_path(entry_matches[0].path)

    if show:
        print(old_entry.display())
        return

    if delete:
        prompt = "Deleting entry '%s'. Are you sure? [y/N] " % old_entry.title
        result = input(prompt)

        if result.lower() in ['y', 'yes']:
            Journal.delete_entry(old_entry)
            git.delete([old_entry.path])
            git.commit("delete-entry '%s'" % old_entry.title)

            print('Deleted journal entry', old_entry.title)
        else:
            print('Not deleting journal entry', old_entry.title)

        return

    # Edit
    with git.chdir(os.path.join(config.data_dir(), '.git')):
        with open('COMMIT_EDITMSG', 'w') as entry_text:
            entry_text.write(old_entry.title + '\n' + old_entry.text + '\n')

        editor('COMMIT_EDITMSG')

        new_entry = Entry.from_path('COMMIT_EDITMSG')

        if new_entry.title == old_entry.title and \
            new_entry.text == old_entry.text:

            # No change
            log.info('No change to entry')
            return

        # Entry has changed
        new_entry.save()
        Journal.delete_entry(old_entry)
        Journal.add_entry(new_entry)

        git.delete([old_entry.path])
        git.add([new_entry.path])

        if new_entry.title == old_entry.title:
            # No change to title
            commit_msg = "edit-entry '%s'" % new_entry.title
        else:
            commit_msg = "edit-entry '%s'->'%s'" % \
                (old_entry.title, new_entry.title)
        git.commit(commit_msg)

def list_or_display_entries(tags, display=False):
    """List or display all journal entries, filtered by tags"""
    journal = Journal()
    entries = journal.filter(tags)

    if display:
        # Display entries
        pager = os.popen('less -FRX', mode='w')
        for entry in entries:
            pager.write(Entry.from_path(entry.path).display())
        pager.close()

    else:
        # List entries
        fmt_str = '%-12s%-12s%s'
        print(fmt_str % ('ID', 'Date', 'Title'))
        print(terminal.horizontal_line())

        for entry in entries:
            print(fmt_str % (entry.entry_id, entry.date, entry.title))

def display_tags():
    """Display all tags saved in the database"""
    journal = Journal()
    tags = set()

    for entry in journal.entries:
        tags.update(entry.tags)

    print('\n'.join(sorted(tags)))

########################################################################
# CLI parser
########################################################################
JOURNAL_DESC = '''
The Overlord Journal Log allows you to keep an activity log. Entries are
automatically saved with the current timestamp, and you may add optional
tags to each entry to allow for filtering in the future. Tags may
contain the characters a-z, 0-9 and hyphen (-).
'''

def register_command(parser):
    """Register the journal command with the parent parser"""
    log.debug('Registering journal command')

    journal_cmd = parser.add_parser('journal', help='journal logging',
                                    description=JOURNAL_DESC)

    subcommand = journal_cmd.add_subparsers(title='Subcommands', dest='subcmd',
                                            metavar='')

    # Parent parsers for common arguments
    id_parser = argparse.ArgumentParser(add_help=False)
    id_parser.add_argument('id', nargs=1, help='entry id')

    tag_parser = argparse.ArgumentParser(add_help=False)
    tag_parser.add_argument('tags', nargs='*', help='tags')

    subcommand.add_parser('new', parents=[tag_parser],
                          help='add new journal entry with optional tags')
    subcommand.add_parser('list', parents=[tag_parser],
                          help='list all journal entries filtered by tags')
    subcommand.add_parser('delete', parents=[id_parser],
                          help='delete the entry by the given id')
    subcommand.add_parser('show', parents=[id_parser],
                          help='display the entry by the given id')
    subcommand.add_parser('edit', parents=[id_parser],
                          help='edit the entry by the given id')
    subcommand.add_parser('tags',
                          help='display all tags in the journal')
    subcommand.add_parser('display', parents=[tag_parser],
                          help='display all journal entries filtered by tags')
    subcommand.add_parser('retag', parents=[id_parser, tag_parser],
                          help='retag the journal entry with the given tags')

def handle_command(parser, args):
    """Handle journal command line"""
    if args.subcmd is None:
        parser.print_help()
    elif args.subcmd == 'new':
        log.debug('journal new ' + ' '.join(args.tags))
        create_entry(args.tags)
    elif args.subcmd in ['list', 'display']:
        log.debug('journal ' + args.subcmd + ' ' + ' '.join(args.tags))
        list_or_display_entries(args.tags, display=(args.subcmd == 'display'))
    elif args.subcmd == 'tags':
        log.debug('journal tags')
        display_tags()
    elif args.subcmd == 'show':
        log.debug('journal show ' + args.id[0])
        operate_on_entry(args.id[0], show=True)
    elif args.subcmd == 'delete':
        log.debug('journal delete ' + args.id[0])
        operate_on_entry(args.id[0], delete=True)
    elif args.subcmd == 'edit':
        log.debug('journal edit ' + args.id[0])
        operate_on_entry(args.id[0])
    else:
        log.error('not implemented - journal %s', args.subcmd)

########################################################################
# Backup and Restore Handlers
########################################################################
def export_data(backup_dict):
    """Export the Journal entries as a dictionary"""
    backup_dict.update(Journal().to_dict())

def import_data(journaldict):
    """Import the Journal entries from a dictionary"""
    return Journal.from_dict(journaldict)

########################################################################
# Init Handlers
########################################################################
def module_init(force=False):
    """Initialize the module"""
    Journal.module_init()
