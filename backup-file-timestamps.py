#!/usr/bin/python

# Utility script for saving and restore the modification times for all files in a tree

from __future__ import print_function

import argparse
import json
import os
import sys
import time

def collect_file_attrs(path):
    dirs = os.walk(path)
    file_attrs = {}
    for (dirpath, dirnames, filenames) in dirs:
        files = dirnames + filenames
        for file in files:
            path = os.path.join(dirpath, file)
            file_attrs[path] = {
                'mtime' : os.path.getmtime(path)
            }
    return file_attrs

def apply_file_attrs(attrs):
    for path in sorted(attrs):
        attr = attrs[path]
        if os.path.lexists(path):
            mtime = attr['mtime']

            mtime_changed = os.path.getmtime(path) != mtime

            if mtime_changed:
                print('Updating mtime for %s' % path, file=sys.stderr)
                os.utime(path, (mtime, mtime))
        else:
            print('Skipping non-existent file %s' % path, file=sys.stderr)

def dir_path(string):
    if os.path.isdir(string):
        return string
    else:
        print("error: - invalid path to folder provided")
        sys.exit(1)

def main():
    ATTR_FILE_NAME = '.saved-file-timestamps'
    parser = argparse.ArgumentParser('Save / Restore timestamps for all files in a directory tree \nexamples:\n-save \"C:\\myfolder\"\n-restore \"C:\\myfolder\"')
    parser.add_argument('path', nargs='?', type=dir_path, help='Path to the directory (uses current by default)')
    parser.add_argument('-save', type=dir_path, help='Save the timestamps of files in the directory tree (used by default)')
    parser.add_argument('-restore', type=dir_path, help='Restore saved file timestamps in the directory tree')
    args = parser.parse_args()

    if args.restore:
        filepath = os.path.join(args.restore, ATTR_FILE_NAME)
        if 'raw_input' in vars(__builtins__):
            inp = raw_input
        else:
            inp = input
        if not os.path.exists(filepath):
            print('Timestamps file \'%s\' not found' % filepath, file=sys.stderr)
            sys.exit(1)
        if inp("- Are you sure you want to restore timestamps to previously saved ones? (Y/N)").lower() == 'y':
            attr_file = open(filepath, 'r')
            attrs = json.load(attr_file)
            apply_file_attrs(attrs)
        else:
            print("- Aborted")
            sys.exit(1)
    else:
        if args.save:
            path = args.save
        elif args.path:
            path = args.path
        else:
            path = os.path.abspath(os.getcwd())
        filepath = os.path.join(path, ATTR_FILE_NAME)
        print('Saving timestamps to: \'%s\'' % filepath)
        attr_file = open(filepath, 'w')
        attrs = collect_file_attrs(path)
        json.dump(attrs, attr_file)
        print('Save complete !')
    time.sleep(3)
if __name__ == '__main__':
    main()
