var Validator = (function() {
  var validator_target = 'http://localhost:6767/';
  var textarea_selector = '#content';
  var version_selector = '#version';
  var submit_selector = '#submit_validate';

  var req = undefined;

  prepareNewRequest();

  document.querySelector(submit_selector).onclick = validate;

  req.onreadystatechange = function(event) {
    console.log(event);

    if (req.readyState === 4) {
      prepareNewRequest();
    }
  };

  return {
    req: req,
    validate: validate,
    getValue: getValue,
  }

  function validate() {
    var content = getValue(textarea_selector);
    var version = getValue(version_selector);

    req.send({ version: version, schema: content, });
  }

  function getValue(selector) {
    var element = document.querySelector(selector);

    return element.value || element.innerHTML || '';
  }

  function prepareNewRequest() {
    req = new XMLHttpRequest();

    req.open('POST', validator_target, true);
    req.setRequestHeader('Content-Type', 'application/json');

    req.onreadystatechange = console.log;
  }
})();
