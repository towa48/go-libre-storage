import React from 'react';
import Button from '@material-ui/core/Button';
import './Sidebar.scss';

export function Sidebar() {
    return (
        <div class="sidebar sidebar_fixed">
            <div class="sidebar__buttons">
                <Button variant="contained"color="primary">Upload</Button>
                <Button variant="contained">Create</Button>
            </div>
            <div class="navigation sidebar__navigation">
                <div class="navigation__items">
                    <div class="navigation__item">
                        <a href="#">Files</a>
                    </div>
                </div>
            </div>
        </div>
    );
}
