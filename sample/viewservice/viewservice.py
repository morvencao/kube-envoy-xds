from flask import Flask
from flask import request
import socket
import os
import sys
import requests

app = Flask(__name__)

TRACE_HEADERS_TO_PROPAGATE = [
    'X-Ot-Span-Context',
    'X-Request-Id',

    # Zipkin headers
    'X-B3-TraceId',
    'X-B3-SpanId',
    'X-B3-ParentSpanId',
    'X-B3-Sampled',
    'X-B3-Flags',

    # Jaeger header (for native client)
    "uber-trace-id"
]

@app.route('/view/<id>')
def view(id):
    headers = {}
    for header in TRACE_HEADERS_TO_PROPAGATE:
        if header in request.headers:
            headers[header] = request.headers[header]
    # call model service from view service
    resp = requests.get("http://localhost:8200/data/" + id, headers=headers)
    out = resp.text.split('-')
    
    return ('Get data: <b>{}</b> from {} behind Envoy! <br>hostname: {} resolved<br>'
            'hostname: {}\n'.format(out[0], out[1], socket.gethostname(),
                                    socket.gethostbyname(socket.gethostname())))

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=9080, debug=True)
