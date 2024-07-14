import React from 'react';
import Plot from 'react-plotly.js';

class Stock extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      stockSymbol: '',
      stockChartXValues: [],
      stockChartYValues: []
    };
    
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({ stockSymbol: event.target.value });
  }

  handleSubmit(event) {
    event.preventDefault();
    this.fetchStock();
  }

  fetchStock() {
    const pointerToThis = this;
    let StockSymbol = this.state.stockSymbol;
    const API_KEY = ''
    let API_Call = `https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=${StockSymbol}&outputsize=compact&apikey=${API_KEY}`;
    let stockChartXValuesFunction = [];
    let stockChartYValuesFunction = [];

    fetch(API_Call)
      .then(response => response.json())
      .then(data => {
        if (data['Time Series (Daily)']) {
          for (let key in data['Time Series (Daily)']) {
            stockChartXValuesFunction.push(key);
            stockChartYValuesFunction.push(data['Time Series (Daily)'][key]['1. open']);
          }

          pointerToThis.setState({
            stockChartXValues: stockChartXValuesFunction,
            stockChartYValues: stockChartYValuesFunction
          });
        } else {
          console.error('Error fetching data:', data);
        }
      })
      .catch(error => console.error('Error fetching stock data:', error));
  }

  render() {
    return (
      <div>
        <h1>Stock Market</h1>
        <form onSubmit={this.handleSubmit}>
          <label>
            What stock chart would you like to see?
            <input 
              type="text" 
              value={this.state.stockSymbol} 
              onChange={this.handleChange} 
              placeholder="Enter stock symbol" 
            />
          </label>
          <button type="submit">Submit</button>
        </form>
        <Plot
          data={[
            {
              x: this.state.stockChartXValues,
              y: this.state.stockChartYValues,
              type: 'scatter',
              mode: 'lines+markers',
              marker: { color: 'red' },
            }
          ]}
          layout={{ width: 720, height: 440, title: `Stock Chart for ${this.state.stockSymbol.toUpperCase()}` }}
        />
      </div>
    );
  }
}

export default Stock;