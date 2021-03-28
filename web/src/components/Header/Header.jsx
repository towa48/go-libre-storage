import React from 'react';
import Avatar from '@material-ui/core/Avatar';
import PersonIcon from '@material-ui/icons/Person';
import './Header.scss';

export function Header() {
    return (
        <div class="header header_fixed">
            <div class="header__inner header__inner_left">
                STORAGE
            </div>
            <div class="header__inner header__inner_right">
                <Avatar>
                    <PersonIcon />
                </Avatar>
            </div>
        </div>
    );
}
