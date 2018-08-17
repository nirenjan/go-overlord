"""Evil Overlord Modules handler"""

import sys

from . import log

def find_module(module_name, module_list=None):
    """Find the module given the module_list"""
    module_list = module_list or sys.modules

    return module_list.get('eviloverlord.' + module_name, None)

def find_module_function(module, function):
    """Find the function handle, given the module and function name"""
    return getattr(module, function, None)

def run_callback(callback_func, *args, **kwargs):
    """Run callback function in each of the Evil Overlord modules"""
    module_list = sys.modules.copy()
    for module_name in module_list:
        if not module_name.startswith('eviloverlord'):
            continue

        module = module_list[module_name]
        callback = find_module_function(module, callback_func)

        if callback is not None:
            log.debug('Calling %s.%s', module_name, callback_func)
            callback(*args, **kwargs)
