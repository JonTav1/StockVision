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
        <div className = "aisummary">
            <div>
            <form>
            <label>
                What is the Ticker symbol?
            <input 
                type = "text"
                value = {ticker}
                onChange = {(e) => setTicker(e.target.value)}
                />
            </label>
       
            
            <div>
                Is this a cryptocurrency?
                <button type = "button" onClick = {(e) => setChoice("crypto")}>Yes</button>
                <button type = "button" onClick = {(e) => setChoice("stock")}>No</button>
                
            </div>
            <br></br>
            <button onClick = {handleSubmit}>Submit Information</button>
            </form>
            </div>
            <p>{response}</p>
        </div>
        
    );
}

export default AiSummary; 