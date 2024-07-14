import React, {useState} from 'react';


function BuyStock({username}) {
    const [response, setResponse] = useState('');  
    const [ticker, setTicker] = useState('');
    const [amount, setAmount] = useState('');
    const [buyprice, setBuyPrice] = useState('');
    const handleSubmit = (e) => {
        e.preventDefault();
        
        
        const data = { 
            username: username,
            ticker: ticker,
            amount: parseInt(amount),
            averageprice: parseInt(buyprice)
        };
        console.log(data)
        
        fetch("http://localhost:8080/stocks/addstock", {
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
                Buy Price
                <input
                type = "text"
                value = {buyprice}
                onChange = {(e) => setBuyPrice(e.target.value)}
                />
              </label>
                <button type="submit">Purchase Stock</button>
            </form>
            <h1>{response}</h1>
        </div>
    )
}
export default BuyStock