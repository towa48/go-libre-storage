class LocationService {
    url(val) {
        if (val) {
            window.location.href = val;
            return;
        }

        return window.location.href;
    }
}

export const $location = new LocationService();
