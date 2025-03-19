import * as cryptellation from 'cryptellation-gateway';

import './BottomNavBar.css';

export function BottomNavBar(client : cryptellation.Client) {
    var info = client.info();

    return (
        <nav className="navbar fixed-bottom navbar-expand-sm bottom-navbar">
            <ul className="navbar-nav mr-auto">
                <li className="nav-item">
                    <div className="bottom-navbar-info">Cryptellation {info}</div>
                </li>
            </ul>
        </nav>
    );
}