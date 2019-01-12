from flask import Flask
from flask import request
import socket
import os
import sys
import requests

app = Flask(__name__)

@app.route('/view/<id>')
def view(id):
   # call model service from view service
   resp = requests.get("http://localhost:8200/data?id=" + id)
   
   return ('Get data: {} from behind Envoy! hostname: {} resolved'
          'hostname: {}\n'.format(resp.text, socket.gethostname(),
                                  socket.gethostbyname(socket.gethostname())))

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=9080, debug=True)
