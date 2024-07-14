import React, { useState } from 'react';

function CheckStocks({ username }) {
    const [response, setResponse] = useState('');  
    const [stocks, setStocks] = useState([]);

    const handleSubmit = (e) => {
        e.preventDefault();
    
        const url = `http://localhost:8080/stocks?username=${username}`;
        fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('API response data:', data);
            
            if(data.stocks.length === 0){
                setResponse('You currently have no stocks.')
            }
            else{
            setStocks(data.stocks); // Update stocks state correctly
            setResponse('Stocks fetched successfully.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            setResponse('Failed to fetch stocks data. Please try again later.');
        });
    };

    return (
        <div>
            <div>
            <i>
                Press the Button below to check your stock portfolio
            </i>
            </div>
            <button onClick={handleSubmit}>Check Stocks</button>
            <h1>{response}</h1>
            <ul>
                {stocks.map((stock, index) => (
                    <li key={index}>
                        <strong>Ticker:</strong> {stock.Ticker}<br />
                        <strong>Amount:</strong> {stock.Amount}<br />
                        <strong>Average Price:</strong> {stock.AveragePrice}<br />
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default CheckStocks;