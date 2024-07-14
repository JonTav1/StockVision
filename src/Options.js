import React, { useState } from 'react';
import NewStock from './NewStock'; 
import CheckStocks from './CheckStocks'; 
import DeleteStock from './DeleteStock'; 
import AiSummary from './AiSummary';
import BuyStock from './BuyStock';
import SellStock from './SellStock';
import Stock from './StockChart';
function Options({username}) {
    const [selectedOption, setSelectedOption] = useState(null);
    
    

    return (
        <div>
                <button onClick={() => setSelectedOption('checkStocks')}>Check Stocks</button>
                <button onClick={() => setSelectedOption('addStock')}>New Stock</button>
                <button onClick={() => setSelectedOption('removeStock')}>Delete Stock</button>
                <button onClick={() => setSelectedOption('aiSummary')}>Ai Stock Summary</button>
                <button onClick={() => setSelectedOption('buystock')}>Buy Stock</button>
                <button onClick={() => setSelectedOption('sellstock')}>Sell Stock</button>
                <button onClick={() => setSelectedOption('stockchart')}>Stock Chart</button>

            {selectedOption === 'checkStocks' && <CheckStocks username={username} />}
            {selectedOption === 'addStock' && <NewStock username={username} />}
            {selectedOption === 'removeStock' && <DeleteStock username={username} />}
            {selectedOption === 'aiSummary' && <AiSummary/>}  
            {selectedOption === 'buystock' && <BuyStock username={username}/>} 
            {selectedOption === 'sellstock' && <SellStock username={username}/>} 
            {selectedOption === 'stockchart' && <Stock/>} 



            
        </div>
        
        
    );
}

export default Options;