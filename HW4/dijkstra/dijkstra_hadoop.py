from mrjob.job import MRJob
from mrjob.step import MRStep
from mrjob.protocol import RawProtocol

class Node:
    def __init__(self, name: str, distance: float, neighbors: dict):
        self.name = name
        self.distance = distance
        self.neighbors = neighbors

class Dijkstra(MRJob):
    OUTPUT_PROTOCOL = RawProtocol

    def mapper(self, _, line):
        node_name, values = line.split(':', 1)
        node_name = node_name.strip()
        values = eval(values.strip()) 
        node = Node(node_name, float(values['Distance']), values['AdjacencyList'])
        if node.distance != float("inf"):
            for n, w in node.neighbors.items():
                yield (str(n), w + node.distance)
        yield (node_name, [node.distance, str(node.neighbors)])

    def reducer(self, node_name, values):
        neighbors = "{}"
        min_distance = float('inf')
        for value in values:
            if isinstance(value, list):
                min_distance = min(value[0], min_distance)
                # if len(value[1]) > 3:
                neighbors = value[1]
            else:
                min_distance = min(float(value), min_distance)
        yield f"{node_name}:", "{{\"Distance\": float(\"{}\"), \"AdjacencyList\": {}}}"\
               .format(min_distance, neighbors)

    def steps(self):
        return [
            MRStep(mapper=self.mapper,
                   reducer=self.reducer
                   ),
        ]
if __name__ == "__main__":
    Dijkstra.run()
