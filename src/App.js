import React, { useState } from 'react';
import Options from './Options'
import Login from './Login'
import CreateUser from './CreateUser';
function App() {
    const [LoggedIn, setLoggedIn] = useState(false)
    const [username, setUsername] = useState("")
    const [selectedOption, setSelectedOption] = useState("") 
    const successfulLogin = (username) => {
      setLoggedIn(true);
      setUsername(username)
    };
    return (
      <div>
          <h1>Stock Portfolio Organizer</h1>
          <button onClick={() => setSelectedOption('login')}>Login</button>
          <button onClick={() => setSelectedOption('createuser')}>Create User</button>

          {selectedOption === 'createuser' && <CreateUser />}
          {!LoggedIn && selectedOption === 'login' && <Login loginSuccess={successfulLogin} />}
          
          {LoggedIn ? (
              <Options username={username} />
          ) : (
              selectedOption !== 'login' && <Login loginSuccess={successfulLogin} />
          )}
      </div>
  );
}
export default App;