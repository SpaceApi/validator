import json
from typing import Tuple

import jsonschema
from jsonschema.exceptions import ValidationError


def validate(schema_path: str, data: dict) -> Tuple[bool, str]:
    """
    Validate the data with the specified schema.

    Return a tuple where the first value indicates whether the data is valid,
    and the second value contains an error message for the case where the data
    is invalid.

    """
    with open(schema_path, 'r') as f:
        schema = json.loads(f.read())
    # TODO: Cache schema and instantiate a Draft4Validator directly
    try:
        jsonschema.validate(data, schema, cls=jsonschema.Draft4Validator)
    except ValidationError as e:
        return False, e.message
    return True, None
