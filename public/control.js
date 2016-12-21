/* global Handlebars, tinyxhr, btoa */

// renderTemplate gets  the template at public/templates/<name>.hbs
// Renders the template with the provided <context>
// And replaces the resulting html into element with ID=<targetid>
// On error just prints to console and exits.
var renderTemplate = function (name, targetid, context, cb = function () {}) {
  var consumeTemplate = function (err, source, xhr) {
    if (err) {
      console.log('Could not get ' + name, err);
      return;
    }
    var template = Handlebars.compile(source);
    document.getElementById(targetid).innerHTML = template(context);
    cb();
  };
  tinyxhr('public/templates/' + name + '.hbs', consumeTemplate, 'GET');
};

var Control = function () {
  // The http basic auth string for Authorization header

  this.verbose = true;
  // tracks which page is currently displayed
  this.currentpage = '';
  // authstring is basic authentication string set by setauth
  this.authstring = '';
  // the logged in user
  this.user = {};

  // General methods
  this.log = function (message, importance = 0) {
    if (this.verbose || importance > 3) {
      console.log(message);
    }
  };

  // User callback methods
  // setauth gets user supplied name and password and tests it.
  this.setauth = function () {
    if (this.currentpage !== 'auth') {
      this.log('setauth method requires authpage', 4);
      return;
    }
    var un = document.getElementById('username').value;
    var pw = document.getElementById('password').value;
    this.authstring = 'Basic ' + btoa(un + ':' + pw);

    var consumeresponse = function (status, response, xhr) {
      if (xhr.status !== 200) {
        this.setpage('auth', {'error': xhr.responseText});
      } else {
        this.setuser(xhr.response);
      }
    };
    this.call(
        'auth',
        consumeresponse.bind(this)
    );
  };
  this.setuser = function (response) {
    this.user = JSON.parse(response);
    this.setpage('home');
  };
  this.patchuser = function () {
    var newuser = this.user;
    newuser.name = document.getElementById('username').value;
    newuser.email = document.getElementById('useremail').value;
    newuser.phone = document.getElementById('userphone').value;

    var consumeresponse = function (status, response, xhr) {
      if (xhr.status !== 200) {
        this.setpage('home', {'user': this.user, 'error': xhr.responseText});
      } else {
        this.setuser(xhr.response);
      }
    };

    this.call('users/' + this.user.id, consumeresponse.bind(this), 'PATCH', JSON.stringify(newuser), 'application/json');
  }


  // Rendering methods that manipulate the DOM
  this.setpage = function (page, context = {}) {
    this[page + 'page'](context);
    this.currentpage = page;
    this.log('page set to ' + this.currentpage, 0);
  };
  this.authpage = function (context) {
    var settriggers = function () {
      document.getElementById('login').addEventListener('click', this.setauth.bind(this), false);
    };
    renderTemplate('auth', 'container', context, settriggers.bind(this));
  };
  this.homepage = function (context) {
    var settriggers = function () {
      document.getElementById('patchuser').addEventListener('click', this.patchuser.bind(this), false);
    };
    context['user'] = this.user;
    renderTemplate('home', 'container', context, settriggers.bind(this));
  };
  // Ajax methods that callback to doorserver
  this.call = function (url, cb, method = 'GET', post = '', contenttype = '', headers = {}) {
    headers['Authorization'] = this.authstring;
    tinyxhr(url, cb, method, post, contenttype, headers);
  };
};

var Init = function () {
  var cnt = new Control();
  cnt.setpage('auth');
};
