import React, {useState} from 'react';
import {Link} from 'react-router-dom'

function SellStock({username}) {
    const [response, setResponse] = useState('');  
    const [ticker, setTicker] = useState('');
    const [amount, setAmount] = useState('');
    const [sellprice, setSellPrice] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        
        
        const data = { 
            username: username,
            ticker: ticker,
            amount: parseInt(amount),
            averageprice: parseInt(sellprice)
        };
        console.log(data)
        
        fetch("https://pqlk51hogh.execute-api.us-east-2.amazonaws.com/stocks/removestock", {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => setResponse(data.message))
            
            .catch(error => {
                console.error('Error:', error);
                setResponse('Incorrect username/password, or account does not exist, try again.');
            });
    };
    return (
        <div>
        <form onSubmit={handleSubmit}>
              <label>
                Ticker
                <input 
                type = "text"
                value = {ticker}
                onChange = {(e) => setTicker(e.target.value)}
                />
                Amount
                <input
                type = "text"
                value = {amount}
                onChange = {(e) => setAmount(e.target.value)}
                />
                Sell Price
                <input
                type = "text"
                value = {sellprice}
                onChange = {(e) => setSellPrice(e.target.value)}
                />
              </label>
                <button type="submit">Sell Stock</button>
            </form>
            <Link to="/deletestock">
                <button>Sell all of your shares instead?</button>
            </Link>
            <h1>{response}</h1>
        </div>
    )
}
export default SellStock