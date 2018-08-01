# Overlord

Overlord is a command-line based personal assistant. It can take notes, maintain
a journal, set reminders and more.

## Requirements

Overlord requires Git 1.7.3 or later and Bash 4 or later.

## Usage

Overlord is used from the command line. It creates a folder in your $HOME called
~/.overlord/ and saves all application data there.

Overlord expects a command to execute.

### Commands

TODO: List all commands

## Installation

1. Clone the Overlord repository.
2. Add the path to your PATH environment variable. Alternatively, you can create
   a symbolic link to *<overlord-path>*/overlord from a folder in your PATH.
2. Make sure your name and e-mail address is set in the global git configuration
3. Run `overlord init` and it will automatically take care of the rest.

## Contributing

1. Fork it ( https://github.com/nirenjan/overlord/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
