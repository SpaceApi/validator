import pytest
import requests


@pytest.mark.parametrize(['data', 'errormsg'], [
    (None, 'JSON payload missing'),
    ('[]', 'Payload does not contain a "data" field'),
    ('{}', 'Payload does not contain a "data" field'),
    ('asfd', 'Request data is not valid JSON'),
    ('{"data": "asfd"}', 'Data is not valid JSON'),
    ('{"data": "{\\"a\\": \\"pi\\"}"}', 'Data does not contain an "api" field'),
    ('{"data": "{\\"api\\": \\"0.4\\"}"}', 'Unknown api version: "0.4"'),
])
def test_validate_validation(testserver, data, errormsg):
    r = requests.post(testserver + 'v1/validate/',
            data=data,
            headers={'Content-type': 'application/json'})
    assert r.status_code == 400, (r.status_code, r.content)
    assert r.json()['detail'] == errormsg, '%r != %r' % (r.json()['detail'], errormsg)
