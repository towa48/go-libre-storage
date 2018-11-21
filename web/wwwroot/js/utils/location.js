var LocationService = function() {
}

LocationService.prototype.url = function(val) {
    if (val) {
        window.location.href = val;
        return;
    }

    return window.location.href;
}

export default new LocationService();
