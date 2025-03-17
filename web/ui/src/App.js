
import './App.css';
import * as cryptellation from 'cryptellation-gateway';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Cryptellation</h1>
      </header>
    </div>
  );
}

var client = new cryptellation.Client();
client.info();

export default App;
