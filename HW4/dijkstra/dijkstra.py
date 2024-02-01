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
nodes2 = {}

class Node:
    def __init__(self, name: str, distance: float, neighbors: dict):
        self.name = name
        self.distance = distance
        self.neighbors = neighbors

def map(node: Node):
    if node.distance != float("inf"):
        for n, w in node.neighbors.items():
            print(f"{n} {w + node.distance}")
    print("{}:  {{\"Distance\": float(\"{}\"), \"AdjacencyList\": {}}}"
          .format(node.name, node.distance, node.neighbors))

def reduce(node_name: str, dists):
    if dists != None:
        nodes2[node_name].distance = min(min(dists), nodes2[node_name].distance)
    print("{}:  {{\"Distance\": float(\"{}\"), \"AdjacencyList\": {}}}"
        .format(node_name, nodes2[node_name].distance, nodes2[node_name].neighbors))

if __name__ == "__main__":
    count = 0
    while True:
        nodes = {}
        sys.stdout = sys.__stdout__
        with open("input.txt", 'r') as file:
            for line in file:
                node_name, values = line.split(':', 1)
                node_name = node_name.strip()
                values = eval(values.strip()) 
                node = Node(node_name, float(values['Distance']), values['AdjacencyList'])
                nodes[node_name] = node
            count += 1
        # all done in mapper process
        # mapper 
        sys.stdout = open("output.txt", "w")
        for _, node in nodes.items():
            map(node)
        # all done in reducer process
        # reducer
        distances = {}
        nodes2 = {}
        sys.stdout = sys.__stdout__
        with open("output.txt", 'r') as file:
            for line in file:
                values = line.strip().split(' ')
                if len(values) > 2:
                    node_name, values = line.split(':', 1)
                    node_name = node_name.strip()
                    values = eval(values.strip()) 
                    node = Node(node_name, float(values['Distance']), values['AdjacencyList'])
                    nodes2[node_name] = node
                else:
                    if values[0] in distances: 
                        distances[values[0]].append(float(values[1]))
                    else:
                        distances[values[0]] = [float(values[1])]
        sys.stdout = open("input.txt", "w")
        for n in nodes2:
            reduce(n, distances.get(n))
        if count >= len(nodes) - 1:
            break
