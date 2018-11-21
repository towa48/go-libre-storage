var defaultHeaders = {
    'Accept': 'application/json',
    'X-Requested-With': 'XMLHttpRequest'
}

var HttpService = function() {
}

HttpService.prototype.post = function(url, data, options) {
    options = options || {};
    options.method = 'POST';
    options.headers = defaultHeaders;
    options.body = JSON.stringify(data);

    options['Content-Type'] = 'application/json';

    return fetch(url, options);
}

export default new HttpService();