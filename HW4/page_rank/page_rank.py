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

changed = True
G = 0
a = 0.85

class Node:
    def __init__(self, name: str, page_rank: float, neighbors: dict):
        self.name = name
        self.page_rank = page_rank
        self.neighbors = neighbors

def map(node: Node):
    for n in node.neighbors:
        print(f"{n} {(1-a)*node.page_rank/len(node.neighbors)}")
    print("{}:  {{\"PageRank\": {}, \"AdjacencyList\": {}}}"
          .format(node.name, node.page_rank, node.neighbors))

def reduce(node_name: str, ranks: [float]):
    x = False
    summ = a/G
    if ranks != None:
        summ += sum(ranks)
        if nodes2[node_name].page_rank != summ:
            x = True
    nodes2[node_name].page_rank = summ
    print("{}:  {{\"PageRank\": {}, \"AdjacencyList\": {}}}"
        .format(node_name, summ, nodes2[node_name].neighbors))
    return x

if __name__ == "__main__":
    count = 0
    while changed:
        changed = False
        nodes = {}
        sys.stdout = sys.__stdout__
        with open("input.txt", 'r') as file:
            for line in file:
                node_name, values = line.split(':', 1)
                node_name = node_name.strip()
                values = eval(values.strip()) 
                node = Node(node_name, float(values['PageRank']), values['AdjacencyList'])
                nodes[node_name] = node
            count += 1
        # all done in mapper process
        # mapper 
        sys.stdout = open("output.txt", "w")
        for _, node in nodes.items():
            map(node)
        # all done in reducer process
        # reducer
        ranks = {}
        nodes2 = {}
        sys.stdout = sys.__stdout__
        with open("output.txt", 'r') as file:
            for line in file:
                values = line.strip().split(' ')
                if len(values) > 2:
                    node_name, values = line.split(':', 1)
                    node_name = node_name.strip()
                    values = eval(values.strip()) 
                    node = Node(node_name, float(values['PageRank']), values['AdjacencyList'])
                    nodes2[node_name] = node
                else:
                    if values[0] in ranks: 
                        ranks[values[0]].append(float(values[1]))
                    else:
                        ranks[values[0]] = [float(values[1])]
        G = len(nodes2)
        sys.stdout = open("input.txt", "w")
        for n in nodes2:
            changed = changed | reduce(n, ranks.get(n))
        if count >= len(nodes) - 1:
            break
