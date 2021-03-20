import React from 'react';
import { Header } from '../../components/Header/Header';
import { Sidebar } from '../../components/Sidebar/Sidebar';
import { Footer } from '../../components/Footer/Footer';
import 'normalize.css';
import './App.scss';

function App() {
    return (
        <div class="root">
            <Header />
            <Sidebar />
            <div class="root__content">
                <div class="root__content-inner root__content-inner_white">
                    Content
                </div>
                <Footer />
            </div>
        </div>
    );
}

export default App;
