import * as cryptellation from 'cryptellation-gateway';

import './BottomNavBar.css';
import { useEffect, useState } from 'react';

export function BottomNavBar(client : cryptellation.Client) {
    var [info, setInfo] = useState("unknown");

    useEffect(() => {
        client.getInfo().then((response) => {
            if (response.data) {
                setInfo(response.data.version!);
            }
        }).catch((error) => {
            console.error("Error fetching info:", error);
        });
    }, [client]);

    return (
        <nav className="navbar fixed-bottom navbar-expand-sm bottom-navbar">
            <div className="bottom-navbar-info">Cryptellation {info}</div>
        </nav>
    );
}