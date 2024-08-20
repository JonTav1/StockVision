import React, { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import NewStock from './NewStock'; 
import CheckStocks from './CheckStocks'; 
import DeleteStock from './DeleteStock'; 
import AiSummary from './AiSummary';
import BuyStock from './BuyStock';
import SellStock from './SellStock';
import StockChart from './StockChart';
import Login from './Login';
import CreateUser from './CreateUser';
import Navbar from './Navbar';
import Home from './Home';
import Options from './Options';
import './styles.css';
function App() {
    const [LoggedIn, setLoggedIn] = useState(false);
    const [username, setUsername] = useState("");

    useEffect(() => {
      const storedUsername = localStorage.getItem("username");
      console.log("Stored username:", storedUsername); // Debugging log
      if (storedUsername) {
          setUsername(storedUsername);
          setLoggedIn(true);
      }
  }, []);

    const handleLoginSuccess = (username) => {
      console.log("Login successful: ", username);
      setUsername(username);
      setLoggedIn(true);
      localStorage.setItem("username", username);
      console.log("localStorage after login:", localStorage.getItem("username"));
  };

  const handleLogout = () => {
    setLoggedIn(false);
    setUsername("");
    localStorage.removeItem("username");
};

    return (
        <div className = "app">
            <BrowserRouter>
                <Navbar onLogout={handleLogout} />
                <Routes>
                    {!LoggedIn ? (
                        <>
                            <Route path="/login" element={<Login loginSuccess={handleLoginSuccess} />} />
                            <Route path="/createuser" element={<CreateUser loginSuccess={handleLoginSuccess} />} />
                            <Route path="*" element={<Navigate to="/login" />} />
                        </>
                    ) : (
                        <>
                            <Route path="/" element={<Home />} />
                            <Route path="/newstock" element={<NewStock username={username} />} />
                            <Route path="/checkstocks" element={<CheckStocks username={username} />} />
                            <Route path="/buystock" element={<BuyStock username={username} />} />
                            <Route path="/sellstock" element={<SellStock username={username} />} />
                            <Route path="/aisummary" element={<AiSummary />} />
                            <Route path="/chart" element={<StockChart />} />
                            <Route path="/newstock" element={<NewStock />} username = {username}/>
                            <Route path = "/deleteStock" element={<DeleteStock username = {username}/>}/>
                            <Route path="*" element={<Navigate to="/" />} />
                            
                        </>
                    )}
                </Routes>
            </BrowserRouter>
        </div>
    );
}

export default App;