import React, { PropTypes } from 'react';

export default class LoginForm extends React.Component {
    static propTypes = {
        login: PropTypes.string,
        password: PropTypes.string,
        inProgress: PropTypes.bool,
        isPasswordShown: PropTypes.bool,
        onChangeLogin: PropTypes.func,
        onChangePassword: PropTypes.func,
        onSubmit: PropTypes.func,
        error: PropTypes.string,
    };

    static defaultProps = {
        login: '',
        password: '',
        inProgress: false,
        isPasswordShown: false,
        error: null
    };

    constructor(props) {
        super(props);

        this.changeLogin = this.changeLogin.bind(this);
        this.changePassword = this.changePassword.bind(this);
        this.submitForm = this.submitForm.bind(this);
        this.togglePasswordVisibility = this.togglePasswordVisibility.bind(this);
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

    togglePasswordVisibility() {
        this.props.onTogglePasswordVisibility();
    }

    render() { 
        const login = this.props.login;
        const password = this.props.password;
        const isPasswordShown = this.props.isPasswordShown;
        const error = this.props.error;
        const errorClass = error ? ' is-invalid' : '';

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
                    type={isPasswordShown ? 'text' : 'password'}
                    class={`form-control${errorClass}`}
                    name="password"
                    spellcheck="false"
                    autocapitalize="off"
                    placeholder="Password"
                    required="required"
                    value={password}
                    onChange={this.changePassword} />
                <div class="invalid-feedback">Can't find user with specified login and password.</div>
            </div>
            <div class="form-group">
                <div class="form-check">
                    <input
                        type="checkbox"
                        class="form-check-input"
                        id="showPassword"
                        defaultChecked={isPasswordShown}
                        onChange={this.togglePasswordVisibility} />
                    <label class="form-check-label" for="showPassword">Show password</label>
                </div>
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
