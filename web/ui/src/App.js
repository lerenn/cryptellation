
import './App.css';
import * as cryptellation from 'cryptellation';

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
console.log(client.info());

export default App;
