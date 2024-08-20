import React from 'react'
import {Link} from 'react-router-dom'

function Navbar({onLogout}){


    return(
    
        <nav>
            <div className="nav-items" >      
                <Link to="/">Home</Link>
                <Link to="/buystock">Buy Stock</Link>
                <Link to="/sellstock">Sell Stock</Link>
                <Link to="/checkstocks">Check Stocks</Link>
                <Link to="/aisummary">AI Stock Sentiment</Link>
                <Link to="/chart">Stock Chart</Link>
            
            </div>
            <button onClick={() => {
                    console.log("Logout button clicked"); // Debugging log
                    onLogout();
                }}>Logout</button>
           
        </nav>
    
    )
}
export default Navbar