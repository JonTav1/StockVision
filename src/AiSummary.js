import React, { useState } from 'react';


function AiSummary() {
    const [response, setResponse] = useState('');  
    const [ticker, setTicker] = useState('')
    const [choice, setChoice] = useState('')

    const handleSubmit = (e) => {
        e.preventDefault();
    
        const url = `https://pqlk51hogh.execute-api.us-east-2.amazonaws.com/stocks/aisummary/${ticker}/${choice}`;
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
            setResponse(data)
        })
        .catch(error => {
            console.error('Error:', error);
            setResponse('Failed to fetch stocks data. Please try again later.');
        });
    };

    return (
        <div>
            <div>
            <form onSubmit = {handleSubmit}>
            <label>
                What is the Ticker symbol?
            <input 
                type = "text"
                value = {ticker}
                onChange = {(e) => setTicker(e.target.value)}
                />
            </label>
            </form>
            </div>
            <div>
                Is this a cryptocurrency?
                <button onClick = {(e) => setChoice("crypto")}>Yes</button>
                <button onClick = {(e) => setChoice("stock")}>No</button>
            </div>
            
         
            <p>{response}</p>
        </div>
        
    );
}

export default AiSummary; 