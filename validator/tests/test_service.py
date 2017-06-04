import json

import pytest
import requests


@pytest.mark.parametrize(['data', 'errormsg'], [
    (None, 'JSON payload missing'),
    ('[]', 'Payload does not contain a "data" field'),
    ('{}', 'Payload does not contain a "data" field'),
    ('asfd', 'Request data is not valid JSON'),
])
def test_validate_validation(testserver, data, errormsg):
    r = requests.post(testserver + 'v1/validate/',
            data=data,
            headers={'Content-type': 'application/json'})
    assert r.status_code == 400, (r.status_code, r.content)
    assert r.json()['detail'] == errormsg, '%r != %r' % (r.json()['detail'], errormsg)


def test_validate_valid(testserver):
    with open('tests/data/minimal.json', 'r') as f:
        data = f.read()
    r = requests.post(testserver + 'v1/validate/',
            json={'data': data},
            headers={'Content-type': 'application/json'})
    assert r.status_code == 200, (r.status_code, r.content)
    assert r.json()['valid'] is True
    assert r.json()['message'] is None


def test_validate_malformed(testserver):
    data = 'asdf'
    r = requests.post(testserver + 'v1/validate/',
            json={'data': data},
            headers={'Content-type': 'application/json'})
    assert r.status_code == 200, (r.status_code, r.content)
    assert r.json()['valid'] is False
    assert r.json()['message'] == 'Data is not valid JSON'


def test_validate_invalid_missing_api_version(testserver):
    data = json.dumps({'a': 'b'})
    r = requests.post(testserver + 'v1/validate/',
            json={'data': data},
            headers={'Content-type': 'application/json'})
    assert r.status_code == 200, (r.status_code, r.content)
    assert r.json()['valid'] is False
    assert r.json()['message'] == 'Data does not contain an "api" field'


def test_validate_invalid_unknown_api_version(testserver):
    data = json.dumps({'api': '0.4'})
    r = requests.post(testserver + 'v1/validate/',
            json={'data': data},
            headers={'Content-type': 'application/json'})
    assert r.status_code == 200, (r.status_code, r.content)
    assert r.json()['valid'] is False
    assert r.json()['message'] == 'Unknown api version: "0.4"'


def test_validate_invalid_missing_fields(testserver):
    with open('tests/data/missing_url_logo.json', 'r') as f:
        data = f.read()
    r = requests.post(testserver + 'v1/validate/',
            json={'data': data},
            headers={'Content-type': 'application/json'})
    assert r.status_code == 200, (r.status_code, r.content)
    assert r.json()['valid'] is False
    assert r.json()['message'] == '\'logo\' is a required property'
