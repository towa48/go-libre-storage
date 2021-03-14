export const HttpStatus = {
    OK: 200,
    BAD_REQUEST: 400,
    NOT_FOUND: 404,
    SERVER_ERROR: 500
};

const defaultHeaders = {
    'Accept': 'application/json',
    'X-Requested-With': 'XMLHttpRequest'
}

class HttpService {
    post(url, data, options) {
        options = options || {};
        options.method = 'POST';
        options.headers = defaultHeaders;
        options.body = JSON.stringify(data);

        options['Content-Type'] = 'application/json';

        return fetch(url, options);
    }
}

export const $http = new HttpService();
