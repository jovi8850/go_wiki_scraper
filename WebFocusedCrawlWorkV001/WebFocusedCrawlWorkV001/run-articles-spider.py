# This is a program for running the articles spider
# Run from the WebFocusedCrawlWork directory with this command:
# python run-articles-spider.py

import scrapy  # object-oriented framework for crawling and scraping
import os  # operating system commands
import time  # for timing the crawl process
import sys

# Record start time
start_time = time.time()

# make directory for storing complete html code for web page
page_dirname = 'wikipages'
if not os.path.exists(page_dirname):
    os.makedirs(page_dirname)

# function for walking and printing directory structure
def list_all(current_directory):
    for root, dirs, files in os.walk(current_directory):
        level = root.replace(current_directory, '').count(os.sep)
        indent = ' ' * 4 * (level)
        print('{}{}/'.format(indent, os.path.basename(root)))
        subindent = ' ' * 4 * (level + 1)
        for f in files:
            print('{}{}'.format(subindent, f))

# initial directory should have this form (except for items beginning with .):
#    TOP-LEVEL-DIRECTORY-FOR-SCRAPY-WORK
#        RUN-SCAPY-JOB-NAME.py
#        scrapy.cfg
#        DIRECTORY-FOR-SCRAPY
#            __init__.py
#            items.py
#            pipelines.py
#            settings.py
#            spiders
#                __init__.py
#                FIRST-SCRAPY-SPIDER.py
#                SECOND-SCRAPY-SPIDER.py

# examine the directory structure
print("Current directory structure:")
current_directory = os.getcwd()
list_all(current_directory)

# list the available spiders
print('\nScrapy spider names:\n')
os.system('scrapy list')

# decide upon the desired format for exporting output: 
# such as csv, JSON, XML, or jl for JSON lines

# run the scraper exporting results as a comma-delimited text file items.csv
# os.system('scrapy crawl quotes -o items.csv')

# run the scraper exporting results as a JSON text file items.json
# this gives a JSON array 
# os.system('scrapy crawl quotes -o items.json')

# for JSON lines we use this command
print("\nStarting crawl process...")
crawl_command = 'scrapy crawl articles-spider -o items.jl'
result = os.system(crawl_command)

if result == 0:
    print('\nJSON lines successfully written to items.jl')
else:
    print(f'\nCrawl process failed with exit code: {result}')
    sys.exit(1)

# Calculate and display runtime
end_time = time.time()
total_runtime = end_time - start_time
print(f"\nTotal crawl runtime: {total_runtime:.2f} seconds ({total_runtime/60:.2f} minutes)")

# Validate output file
if os.path.exists('items.jl'):
    file_size = os.path.getsize('items.jl')
    print(f"Output file size: {file_size} bytes")
    
    # Count lines in the output file
    with open('items.jl', 'r', encoding='utf-8') as f:
        line_count = sum(1 for line in f)
    print(f"Number of JSON lines: {line_count}")
else:
    print("Warning: items.jl file was not created")

# run the scraper exporting results as a dictionary XML text file items.xml
# os.system('scrapy crawl quotes -o items.xml')