import socket
import os

variable_name = 'POD_IP'
pod_ip = ""
if variable_name in os.environ:
    pod_ip = os.environ[variable_name]

def get_ip_address():
    # Get the local machine's IP address
    host_name = socket.gethostname()
    ip_address = socket.gethostbyname(host_name)
    return ip_address

def start_server(port):
    # Create a socket object
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Bind the socket to a specific address and port
    server_socket.bind(('', port))

    # Listen for incoming connections
    server_socket.listen(1)
    print(f"Server listening on port {port}...")

    while True:
        # Accept a connection from a client
        client_socket, client_address = server_socket.accept()
        print(f"Connection from {client_address}")
        # Send the IP address to the client
        client_socket.sendall(pod_ip.encode())

        # Close the connection
        client_socket.close()

if __name__ == "__main__":
    # Specify the port on which the server will listen
    port_number = 12345

    # Start the server
    start_server(port_number)

