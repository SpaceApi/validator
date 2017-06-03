import json

from bottle import get, post, run, request, response, redirect, abort, error
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


@error(400)
def error400(error):
    response.set_header('Content-Type', 'application/json')
    return json.dumps({'detail': error.body})


@error(404)
def error404(error):
    response.set_header('Content-Type', 'application/json')
    return json.dumps({'detail': 'Page not found'})


@get('/')
def root():
    redirect('/v1/')


@get('/v1/')
def index():
    return {
        'version': __version__,
        'description': 'Space API Validator API',
        'usage': 'Send a POST request in JSON format to /v1/validate/.',
    }


@post('/v1/validate/')
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
    run(host='127.0.0.1', port=6767)
