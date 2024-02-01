import sys

lines = [
    "A:  {\"Distance\": 0, \"AdjacencyList\": {\"B\": 10, \"C\": 5}}",
    "B:  {\"Distance\": float('inf'), \"AdjacencyList\": {\"C\": 2, \"D\": 1}}",
    "C:  {\"Distance\": float('inf'), \"AdjacencyList\": {\"E\": 2, \"B\": 3, \"D\": 9}}",
    "D:  {\"Distance\": float('inf'), \"AdjacencyList\": {\"E\": 4}}",
    "E:  {\"Distance\": float('inf'), \"AdjacencyList\": {\"A\": 7, \"D\": 6}}",
    "F:  {\"Distance\": float('inf'), \"AdjacencyList\": {\"D\": 6}}"
]

lines2 = ["B 10.0",
"C 5.0", 
"C inf",
"D inf",
"E inf",
"B inf",
"D inf"]

nodes = {}

class Node:
    def __init__(self, name: str, distance: float, neighbors: dict):
        self.name = name
        self.distance = distance
        self.neighbors = neighbors

def map(node: Node):
    print(f"{node.name} {node.distance}")
    if node.distance != float("inf"):
        for n, w in node.neighbors.items():
            print(f"{n} {w + node.distance}")

def reduce(node_name: str, dists):
    #nodes[node_name].distance = min(min(dists), nodes[node_name].distance)
    #print(f"{node_name} {nodes[node_name].distance}")
    print(f"{node_name} {min(dists)}")

if __name__ == "__main__":
    with open("shortest_path.txt", 'w') as file:
            pass
    with open("input.txt", 'r') as file:
        for line in file:
            node_name, value = line.split(':', 1)
            node_name = node_name.strip()
            value = eval(value.strip()) 
            node = Node(node_name, float(value['Distance']), value['AdjacencyList'])
            nodes[node_name] = node
    for i in range(len(nodes) - 1):
        # all done in mapper process
        # update shortest_path
        sys.stdout = sys.__stdout__
        with open("shortest_path.txt", 'r') as file:
            for line in file:
                values = line.strip().split(' ')
                nodes[values[0]].distance = min(nodes[values[0]].distance, float(values[1]))
        # mapper 
        sys.stdout = open("output.txt", "w")
        for _, node in nodes.items():
            map(node)
        # all done in reducer process
        # reducer
        distances = {}
        sys.stdout = sys.__stdout__
        with open("output.txt", 'r') as file:
            for line in file:
                values = line.strip().split(' ')
                if values[0] in distances: # this should be fixed
                    distances[values[0]].append(float(values[1]))
                else:
                    distances[values[0]] = [float(values[1])]
        sys.stdout = open("shortest_path.txt", "w")
        for n, dists in distances.items():
            reduce(n, dists)
