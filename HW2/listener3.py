from flask import Flask, jsonify
import subprocess

app = Flask(__name__)

@app.route('/')
def run_command():
    try:
        command = 'fw'
        result = subprocess.check_output(command, shell=True, text=True).strip()
        return jsonify({'result': result})

    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    # Run the Flask app on port 5000
    app.run(port=5000)
