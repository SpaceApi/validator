import json

import bottle
from bottle import request, response, redirect, abort, error
from jsonschema.exceptions import SchemaError

import validation

__version__ = '1.0.0'


SCHEMATA = {
    8: 'schema/8.json',
    9: 'schema/9.json',
    11: 'schema/11.json',
    12: 'schema/12.json',
    13: 'schema/13.json',
}


app = bottle.app()


def _add_cors_headers(response):
    response.headers['Access-Control-Allow-Origin'] = '*'
    response.headers['Access-Control-Allow-Methods'] = 'GET, POST, PUT, OPTIONS'
    response.headers['Access-Control-Allow-Headers'] = \
            'Origin, Accept, Content-Type, X-Requested-With, X-CSRF-Token'


class EnableCors:
    """
    Set CORS headers.

    Taken from https://stackoverflow.com/a/17262900/284318

    """
    name = 'enable_cors'
    api = 2

    def apply(self, fn, context):
        def _enable_cors(*args, **kwargs):
            _add_cors_headers(response)
            if bottle.request.method != 'OPTIONS':
                # Actual request (not a preflight)
                # reply with the actual response
                return fn(*args, **kwargs)

        return _enable_cors


app.install(EnableCors())


@error(400)
def error400(error):
    _add_cors_headers(response)
    response.set_header('Content-Type', 'application/json')
    return json.dumps({'detail': error.body})


@error(404)
def error404(error):
    _add_cors_headers(response)
    response.set_header('Content-Type', 'application/json')
    return json.dumps({'detail': 'Page not found'})


@app.route('/', method=['GET', 'OPTIONS'])
def root():
    redirect('/v1/')


@app.route('/v1/', method=['GET', 'OPTIONS'])
def index():
    return {
        'version': __version__,
        'description': 'Space API Validator API',
        'usage': 'Send a POST request in JSON format to /v1/validate/.',
    }


@app.route('/v1/validate/', method=['POST', 'OPTIONS'])
def validate():
    data = request.json

    # Validate
    if not data:
        abort(400, 'JSON payload missing')
    if 'version' not in data:
        abort(400, 'Payload does not contain a "version" field')
    if 'data' not in data:
        abort(400, 'Payload does not contain a "data" field')
    try:
        version = int(data['version'])
    except ValueError:
        abort(400, 'Invalid version "%s": Not an integer' % data['version'])
    if version not in SCHEMATA:
        abort(400, 'Unknown version: "%s"' % version)
    try:
        data = json.loads(data['data'])
    except json.JSONDecodeError:
        abort(400, 'Data is not valid JSON')  # TODO: regular response

    # Do validation
    try:
        valid, message = validation.validate(schema_path=SCHEMATA[version], data=data)
    except SchemaError:
        abort(500, 'Invalid schema on server! Please contact one of the admins.')
    return {
        'valid': valid,
        'message': message,
    }


if __name__ == '__main__':
    app.run(host='127.0.0.1', port=6767)
