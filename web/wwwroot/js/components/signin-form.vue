<script>
import routes from '../const/routes';
import httpStatus from '../const/httpStatus';
import $http from '../utils/http';
import $location from '../utils/location';

export default {
    data: function() {
        return {
            username: null,
            password: null
        }
    },
    methods: {
        signIn: function() {
            var data = {
                username: this.username,
                password: this.password,
                rememberMe: false
            };

            $http.post(routes.signIn, data).then(function(resp) {
                switch(resp.status) {
                    case httpStatus.ok:
                        resp.json().then(function(respData) {
                            $location.url(respData.Url || '/');
                        });
                        break;
                    case httpStatus.badRequest:
                        resp.json().then(function(respData) {
                            console.log(respData.Error);
                            // TODO
                        });
                        break;
                    case httpStatus.serverError:
                        // TODO
                        break;
                }
            }).catch(function(err) {
                // TODO
                console.log(err);
            })
        }
    }
}
</script>