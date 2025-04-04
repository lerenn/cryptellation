
import './App.css';

import { BottomNavBar } from './BottomNavBar';

import { Client } from 'cryptellation-gateway';

function App() {
  const client = new Client({
    baseUrl: 'https://localhost:8080',
  });

  return (
    <div>
      {BottomNavBar(client)}
    </div>
  );
}

export default App;
