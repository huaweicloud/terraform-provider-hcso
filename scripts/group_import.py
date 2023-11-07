import argparse
import os

def insert_empty_line(input_file, output_file, pattern):
    with open(input_file, 'r') as infile, open(output_file, 'w') as outfile:
        found_match = False

        for line in infile:
            if pattern in line:
                found_match = True
            elif found_match:
                if ')' in line:
                    found_match = False
                else:
                    outfile.write('\n')
                    found_match = False

            outfile.write(line)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Insert empty lines before the first non-matching line in a file.')
    parser.add_argument('input_file', help='Input file path')
    #parser.add_argument('output_file', help='Output file path')
    parser.add_argument('pattern', help='Pattern to match')

    args = parser.parse_args()

    input_file = args.input_file
    output_file = args.input_file + '.tmp'
    pattern = args.pattern

    insert_empty_line(input_file, output_file, pattern)

    # 重命名输出文件以替换原始文件
    os.replace(output_file, input_file)

