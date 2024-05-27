import csv
import random


def create_receipt_code():
    return "".join(random.sample('0123456789', 8))


phone_number = '18000000000'

num_lines = (20 * 1024 * 1024) // 23

with open('phone_receipt.csv', 'w') as csv_file:
    writer = csv.writer(csv_file)

    for _ in range(num_lines):
        writer.writerow([phone_number, create_receipt_code()])

print(f"Generate {num_lines} lines of data to phone_receipt.csv.")
