#!/usr/bin/env python3

import json
from csv import DictReader

# TODO: make item # equal the number for each book
	# This way we can assume each item is on their specific page

def dumpPage(page, pageNumber, path):
	with open(path + str(pageNumber) + '.json', 'w') as f:
		json.dump(page, f, ensure_ascii=False, indent=4)
	

def main():
	INPATH = 'full_catalog.csv'
	PAGESIZE = 100
	OUTPATH = 'pages/'
	with open(INPATH) as f:
		reader = DictReader(f)
		nextPage = {
			'items': {}
		}
		for i, entry in enumerate(reader):
			nextPage['items'][i+1] = entry # Gutenberg starts indexing at 1, we start indexing at 0 :/
			if (i+1) % PAGESIZE == 0:
				dumpPage(nextPage, (i+1)//PAGESIZE, OUTPATH)
				nextPage = {
					'items': {}
				}
		else:
			if len(nextPage['items']) > 0:
				dumpPage(nextPage, (i+99)//PAGESIZE, OUTPATH)
		

if __name__ == '__main__':
	main()
