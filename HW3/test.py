import socket

def connect_to_server():
    # Server address (replace with the actual IP address or hostname)
    server_address = '192.168.49.2'
    server_port = 30624
    # Create a socket object
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        # Connect to the server
        client_socket.connect((server_address, server_port))
        print(f"Connected to {server_address}:{server_port}")
        # Receive data from the server
        data = client_socket.recv(1024)
        print(f"Received data: {data.decode()}")
    except Exception as e:
        print(f"Error: {e}")
    finally:
        # Close the socket
        client_socket.close()

if __name__ == "__main__":
    connect_to_server()

