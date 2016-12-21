// tinyxhr by Shimon Doodkin - licanse: public doamin - https://gist.github.com/4706967
//
// tinyxhr('http://site.com/ajaxaction',function (err,data,xhr){ if (err) console.log('goterr ',err,'status='+xhr.status); console.log(data)  });
// tinyxhr('http://site.com/ajaxaction',function (err,data,xhr){ if (err) console.log('goterr ',err,'status='+xhr.status); console.log(data)  },'POST','value1=1&value2=2');
// tinyxhr('http://site.com/ajaxaction.json',function (err,data,xhr){ if (err) console.log('goterr ',err,'status='+xhr.status); console.log(data); console.log(JSON.parse(data))  },'POST',JSON.stringify({value:1}),'application/javascript');
// cb - a callback function like: function (err,data,XMLHttpRequestObject){ if (err) throw err;   }

// 2016-12-15 aaron.disilvestro@nike.com Added the headers
function tinyxhr (url, cb, method = 'GET', post = '', contenttype, headers = {}) {
  var requestTimeout, xhr;
  try {
    xhr = new XMLHttpRequest();
  } catch (e) {
    try {
      xhr = new ActiveXObject('Msxml2.XMLHTTP');
    } catch (e) {
      if (console) console.log('tinyxhr: XMLHttpRequest not supported');
      return null;
    }
  }

  // Set timeout
  requestTimeout = setTimeout(function () {
    xhr.abort();
    cb(new Error('tinyxhr: aborted by a timeout'), '', xhr);
  }, 10000);

  // Prepare
  xhr.onreadystatechange = function () {
    if (xhr.readyState !== 4) return;
    clearTimeout(requestTimeout);
    cb(xhr.status !== 200 ? new Error('tinyxhr: server response status is ' + xhr.status) : false, xhr.responseText, xhr);
  };
  xhr.open(method.toUpperCase(), url, true);

  // Add headers
  for (var key in headers) {
    if (headers.hasOwnProperty(key)) {
      xhr.setRequestHeader(key, headers[key]);
    }
  }

  xhr.withCredentials = true;

  if (!post) {
    xhr.send();
  } else {
    xhr.setRequestHeader('Content-type', contenttype ? contenttype : 'application/x-www-form-urlencoded');
    xhr.send(post);
  }
}

// tinyxhr('/test',function (err,data,xhr){ if (err) console.log('goterr ',err); console.log(data)  });
