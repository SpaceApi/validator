import json
try:
    from json import JSONDecodeError
except ImportError:
    # Python <3.5 compat
    JSONDecodeError = ValueError

import bottle
from bottle import request, response, redirect, abort, error
from jsonschema.exceptions import SchemaError

import validation

__version__ = '1.0.0'


SCHEMATA = {
    '0.8': 'schema/8.json',
    '0.9': 'schema/9.json',
    '0.11': 'schema/11.json',
    '0.12': 'schema/12.json',
    '0.13': 'schema/13.json',
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


def invalid_payload(message: str) -> dict:
    """
    Return a response body indicating that the validation was not successful.
    """
    return {'valid': False, 'message': message}


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
    try:
        data = request.json
    except JSONDecodeError:
        abort(400, 'Request data is not valid JSON')

    # Validate
    if data is None:
        abort(400, 'JSON payload missing')
    if 'data' not in data:
        abort(400, 'Payload does not contain a "data" field')
    try:
        data = json.loads(data['data'])
    except JSONDecodeError:
        return invalid_payload('Data is not valid JSON')
    if 'api' not in data:
        return invalid_payload('Data does not contain an "api" field')
    version = data['api']
    if version not in SCHEMATA:
        return invalid_payload('Unknown api version: "%s"' % version)

    # Do validation of submitted endpoint
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
