import React from 'react';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import AddIcon from '@material-ui/icons/Add';
import './Sidebar.scss';

export function Sidebar() {
    return (
        <div class="sidebar sidebar_fixed">
            <div class="sidebar__buttons">
                <Box
                    sx={{
                        '& button': {
                        mb: 2,
                        },
                    }}>
                    <Button
                        variant="contained"
                        fullWidth
                        startIcon={<CloudUploadIcon />}>Upload</Button>
                    <Button
                        variant="outlined"
                        fullWidth
                        startIcon={<AddIcon />}>Create</Button>
                </Box>
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
