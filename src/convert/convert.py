#!/usr/bin/env python
# coding: utf-8

import sys
import argparse

data_dir = None

def parse(name):
    print("Page processing now")
    initial = 'INSERT INTO `{0}` VALUES '.format(name)
    input_name = data_dir + 'jawiki-20150512-{0}.sql'.format(name)
    output_name = data_dir + '{0}.txt'.format(name)
    with open(input_name, 'r') as r:
        with open(output_name, 'w') as w:
            for line in r:
                if not line.startswith(initial):
                    continue
                # 先頭のINITALと末尾の`);`も省く
                split = line[len(initial):-3].split('),')
                for row in split:
                    w.write('p' + row[1:] + "\n")

def split_parse(name, splitlen):
    print("Pagelinks processing now")
    initial = 'INSERT INTO `{0}` VALUES '.format(name)
    input_name = data_dir + 'jawiki-20150512-{0}.sql'.format(name)
    output_name = data_dir + '{0}-'.format(name)
    count = 0
    at = 0
    w = None

    with open(input_name, 'r') as r: 
        while True:
            # can't encode some lines with utf-8
            try:
                line = r.readline()
                if not line:
                    if w: w.close()
                    break
                if not line.startswith(initial):
                    continue
                # eliminate INITAL and `);`
                split = line[len(initial):-3].split('),')
                for row in split:
                    if count % splitlen == 0:
                        if w: w.close()
                        w = open(output_name + str(at) + '.txt', 'w')
                        at += 1
                    w.write('l' + row[1:] + "\n")
                    count += 1
            except UnicodeDecodeError as e:
                print(e)
                continue

def main(args):
    if not args.dir: return

    global data_dir
    data_dir = args.dir
    if data_dir[-1] != '/':
        data_dir += '/'
    if args.page: parse('page')
    if args.pagelinks: split_parse('pagelinks', 10000000)

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-d', '--dir')
    parser.add_argument('-p', '--page', action='store_true')
    parser.add_argument('-pl', '--pagelinks', action='store_true')
    args = parser.parse_args()
    main(args)

