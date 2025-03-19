
import './App.css';

import { BottomNavBar } from './BottomNavBar';

import * as cryptellation from 'cryptellation-gateway';

function App() {
  var client = new cryptellation.Client();

  return (
    <div>
      {BottomNavBar(client)}
    </div>
  );
}

export default App;
