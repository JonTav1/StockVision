import React, {useState} from 'react';


function CreateUser() {
    const [response, setResponse] = useState('');  
    const [firstname, setFirstName] = useState('');
    const [lastname, setLastName] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const handleSubmit = (e) => {
        e.preventDefault();
        
        
        const data = { 
            username: username,
            firstname: firstname,
            lastname: lastname,
            password: password
        };
        console.log(data)
        
        fetch("http://localhost:8080/createuser", {
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
            .then(data => setResponse(data.message))
            
            .catch(error => {
                console.error('Error:', error);
                setResponse('Username already exists. Please choose a different one.');
            });
    };
    return (
        <div>
        <form onSubmit={handleSubmit}>
              <label>
                First Name
                <input 
                type = "text"
                value = {firstname}
                onChange = {(e) => setFirstName(e.target.value)}
                />
                Last Name
                <input
                type = "text"
                value = {lastname}
                onChange = {(e) => setLastName(e.target.value)}
                />
                Username
                <input
                type = "text"
                value = {username}
                onChange = {(e) => setUsername(e.target.value)}
                />
                Password
                <input
                type = "text"
                value = {password}
                onChange = {(e) => setPassword(e.target.value)}
                />
              </label>
                <button type="submit">Create User</button>
            </form>
            <h1>{response}</h1>
        </div>
    )
}
export default CreateUser