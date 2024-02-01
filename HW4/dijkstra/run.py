import subprocess

# Define a custom key function to extract the first number in each row
def get_first_number(row):
    return int(row.split(':')[0])

def print_sorted(file_path):
    with open(file_path, "r") as file:
        data = file.readlines()
    sorted_data = sorted(data, key=get_first_number)
    for item in sorted_data:
        print(item.strip()) 

def count_lines(filename):
    with open(filename, 'r') as file:
        return sum(1 for line in file)

def run_dijkstra(n):
    script_to_execute = './dijkstra_hadoop.py'
    input_file = './graph.txt'
    output_file = './graph2.txt'
    for i in range(n):
        if i % 2 == 0:
            command = ['python3', script_to_execute, input_file]
            with open(output_file, 'w') as output:
                subprocess.run(command, stdout=output)
        else:
            command = ['python3', script_to_execute, output_file]
            with open(input_file, 'w') as output:
                subprocess.run(command, stdout=output)


def main():
    input_file = 'graph.txt'
    nodes_count = count_lines(input_file)
    run_dijkstra(nodes_count)
    print_sorted(input_file)

if __name__ == "__main__":
    main()
