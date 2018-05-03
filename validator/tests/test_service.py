import json

import pytest
import requests


@pytest.mark.parametrize(['data', 'errormsg'], [
    (None, 'JSON payload missing'),
    ('42', 'Payload must be a JSON object'),
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


@pytest.mark.parametrize(['data', 'errormsg'], [
    ('asdf', 'Data is not valid JSON'),
    ('42', 'Data must be a JSON object'),
])
def test_validate_malformed(testserver, data, errormsg):
    r = requests.post(testserver + 'v1/validate/',
            json={'data': data},
            headers={'Content-type': 'application/json'})
    assert r.status_code == 200, (r.status_code, r.content)
    assert r.json()['valid'] is False
    assert r.json()['message'] == errormsg


def _ensure_invalid(testserver, data, status_code: int, message: str):
    r = requests.post(testserver + 'v1/validate/',
            json={'data': data},
            headers={'Content-type': 'application/json'})
    assert r.status_code == status_code, (r.status_code, r.content)
    assert r.json()['valid'] is False
    assert r.json()['message'] == message


def test_validate_invalid_missing_api_version(testserver):
    _ensure_invalid(
        testserver,
        json.dumps({'a': 'b'}),
        200,
        'Data does not contain an "api" field',
    )


def test_validate_invalid_unknown_api_version(testserver):
    _ensure_invalid(
        testserver,
        json.dumps({'api': '0.4'}),
        200,
        'Unknown api version: "0.4"',
    )


def test_validate_invalid_missing_fields(testserver):
    with open('tests/data/missing_url_logo.json', 'r') as f:
        data = f.read()
        _ensure_invalid(
            testserver,
            data,
            200,
            '\'logo\' is a required property',
        )
