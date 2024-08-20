import React, {useState} from 'react';


function DeleteStock({username}) {
    const [response, setResponse] = useState('');  
    const [ticker, setTicker] = useState('');
    

    const handleSubmit = (e) => {
        e.preventDefault();
        

        const data = { 
            username: username,
            ticker: ticker
        };

       console.log(data)
        fetch('https://pqlk51hogh.execute-api.us-east-2.amazonaws.com/stocks/deletestock', {
            method: 'DELETE',
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
            .then(data => setResponse(data.message)
        )
            
            .catch(error => {
                console.error('Error:', error);
                setResponse('Could not delete stock');
            });
    };
    return (
        <div className = "deletestock">
        <form onSubmit={handleSubmit}>
              <label>
                Ticker
                <input 
                type = "text"
                value = {ticker}
                onChange = {(e) => setTicker(e.target.value)}
                />
              </label>
                <button type="submit">Sell all stock</button>
            </form>
            <h1>{response}</h1>
        </div>
    )
}
export default DeleteStock