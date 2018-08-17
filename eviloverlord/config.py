"""Overlord Configuration"""

import os

def data_dir():
    """Path to overlord data"""
    overlord_dir = os.getenv('OVERLORD_DATA',
                             os.path.join(os.getenv('HOME'), '.overlord'))
    overlord_dir = os.path.normpath(overlord_dir)
    os.makedirs(overlord_dir, exist_ok=True)

    return overlord_dir
