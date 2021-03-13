import React, { PropTypes } from 'react';

export default class LoginForm extends React.Component {
    static propTypes = {
        login: PropTypes.string,
        password: PropTypes.string,
        inProgress: PropTypes.bool,
        onChangeLogin: PropTypes.func,
        onChangePassword: PropTypes.func,
        onSubmit: PropTypes.func,
    };

    static defaultProps = {
        login: '',
        password: '',
    };

    constructor(props) {
        super(props);

        this.changeLogin = this.changeLogin.bind(this);
        this.changePassword = this.changePassword.bind(this);
        this.submitForm = this.submitForm.bind(this);
    }

    changeLogin(event) {
        this.props.onChangeLogin(event.target.value);
    }

    changePassword(event) {
        this.props.onChangePassword(event.target.value);
    }

    submitForm(login, password) {
        return (event) => {
            event.preventDefault();
            this.props.onSubmit(login, password);
        }
    }

    render() { 
        const login = this.props.login;
        const password = this.props.password;
        return (
        <form onSubmit={this.submitForm(login, password)}>
            <h2 class="text-center">Login</h2>
            <div class="form-group">
                <input
                    type="text"
                    class="form-control"
                    name="username"
                    placeholder="Username"
                    required="required"
                    value={login}
                    onChange={this.changeLogin} />
            </div>
            <div class="form-group">
                <input
                    type="password"
                    class="form-control"
                    name="password"
                    placeholder="Password"
                    required="required"
                    value={password}
                    onChange={this.changePassword} />
            </div>
            <div class="form-group">
                <button
                    type="submit"
                    class="btn btn-primary btn-lg btn-block"
                    disabled={this.props.inProgress}>
                    Sign in
                </button>
            </div>
        </form>
        );
    }
}
