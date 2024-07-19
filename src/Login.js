import React, {useState} from 'react';

function Login({loginSuccess}) {
const [response, setResponse] = useState('');  
const [password, setPassword] = useState('');
const [username, setUsername] = useState('');
const handleSubmit = (e) => {
        e.preventDefault();

        const data = { 
            username: username,
            password: password
        };

        fetch('https://pqlk51hogh.execute-api.us-east-2.amazonaws.com/login', {
            method: 'POST',
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
            .then(data =>  {
                setResponse(data.message)
                loginSuccess(username)
            })
            
            .catch(error => {
                console.error('Error:', error);
                setResponse('Incorrect username/password, or account does not exist, try again.');
            });
        
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
              <label>
                Username:
                <input 
                type = "text"
                value = {username}
                onChange = {(e) => setUsername(e.target.value)}
                />
                Password:
                <input
                type = "text"
                value = {password}
                onChange = {(e) => setPassword(e.target.value)}
                />
              </label>
                <button type="submit">Sign in</button>
            </form>
            
            <p>{response}</p>
        </div>
    );
}
export default Login