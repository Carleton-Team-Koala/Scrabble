import { React, useState } from 'react';
import { Link } from "react-router-dom";
import '../css/Popup.css';

export default function Popup(props) {

    const [gameID, setGameID] = useState('');

    const handleSubmit = () => {
        const player = document.getElementById('playerName').value;
        sessionStorage.setItem('playerName', player);

        if (props.type === 'joinGame') {
            setGameID(document.getElementById('gameId').value);
            console.log(gameID);
            sessionStorage.setItem('gameId', gameID);
            console.log("storage ", sessionStorage.getItem('gameId'));
        }

        props.onSubmit();
    };

    return (props.trigger) ? (
        <div className='popup'>
            <button className='close-btn' onClick={() => props.setTrigger(false)}>Close</button>
            <div className='inner-container'>
                <h1>Enter Username:</h1>
                <input type='text' id='playerName'></input>

                {props.type === 'joinGame' && (
                    <div>
                        <h1>Enter Game ID:</h1>
                        <input type='text' id='gameId'></input>
                    </div>
                )}

                <Link to={`/play/${gameID}`}>
                    <button className='submit-btn' onClick={handleSubmit}>Submit</button>
                </Link>
            </div>
        </div>
    ) : "";
}
