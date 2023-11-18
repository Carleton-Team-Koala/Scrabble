import React, { useState, useEffect } from 'react';
import { baseURL } from "./Welcome";
import Board from "../components/Board";
import ActionPanel from "../components/ActionPanel";
import Infoboard from "../components/Infoboard";
import Tile from '../components/Tile';
import Rules from './Rules';
import '../css/Game.css';
import '../css/Rules.css';
import "../css/App.css";

const gameID = sessionStorage.getItem('gameId');
const playerName = sessionStorage.getItem('playerName');

function initializeTiles(hand) { // initialize tiles for the board and hand
  return Array.from({ length: hand.length }, (_, i) => ({
    id: i,
    letter: hand[i] === 'BLANK' ? '' : hand[i],
    position: 'ActionPanel', // initial position
  }));
};

export default function Game() {

  /**
   * Main controller component for the actual game.
   * Controls the display of the information to the user.
   * Receives hand and tilebag from the initialization.
   */

  const [hand, setHand] = useState(['BLANK', 'B', 'C', 'D', 'E', 'A', 'G']); // array of letters, gets rendered in the hand
  const [tilebag, setTilebag] = useState({
    'A': 0, 'B': 0, 'C': 0, 'D': 0, 'E': 0, 'F': 0, 'G': 0,
    'H': 0, 'I': 0, 'J': 0, 'K': 0, 'L': 0, 'M': 0, 'N': 0,
    'O': 0, 'P': 0, 'Q': 0, 'R': 0, 'S': 0, 'T': 0, 'U': 0,
    'V': 0, 'W': 0, 'X': 0, 'Y': 0, 'Z': 0, 'BLANK': 0
  }); // array of letters, gets rendered in the tilebag
  const [scoredLetters, setScoredLetters] = useState({}); // {cellKey: letter}, letters returned by server go here
  const [letterUpdates, setLetterUpdates] = useState({}); // {id: [cellKey, letter]}, gets sent to server on submit
  const [tiles, setTiles] = useState(initializeTiles(hand)); // array of tiles, gets rendered on the board and hand
  const [isRulesOpen, setIsRulesOpen] = useState(false);

  async function getGame() {
    const url = baseURL + "/getgamestate/" + gameID + "/";
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json"
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      // Check if the necessary data is present
      if (data.gameState && data.gameState.Players && data.gameState.Players[playerName]) {
        setHand(data.gameState.Players[playerName].hand);
        setTilebag(data.gameState.LetterDistribution);
      } else {
        console.error("Unexpected data format: ", data);
      }
    } catch (error) {
      alert(`An error occurred: ${error.message}`);
      console.error("Error: ", error);
    }
  }

  // useEffect to fetch game state and set tiles on component mount
  useEffect(() => {
    getGame().then(() => {
      console.log("getGame finished, now setting tiles");
      // This will trigger the second useEffect hook
    });
  }, []);

  // useEffect to update tiles whenever hand changes
  useEffect(() => {
    console.log("Hand updated, now updating tiles");
    setTiles(initializeTiles(hand));
  }, [hand]);

  /**
    * Handles the event when a tile is dropped onto the board.
    * Sets letterUpdates and tiles states.
    *
    * @param {number} id - The ID of the tile that was dropped.
    * @param {string} cellKey - The key of the cell where the tile was dropped.
    * @param {string} letter - The letter on the tile that was dropped.
    */
  function handleTileDrop(id, cellKey, letter) {
    id = Number(id);

    setLetterUpdates(prevState => ({
      ...prevState,
      [id]: [cellKey, letter]
    }));

    setTiles(prevTiles =>
      prevTiles.map(tile =>
        tile.id === id ? { ...tile, position: 'Board' } : tile
      )
    );
  };

  const shuffle = () => {
    let indices = [0, 1, 2, 3, 4, 5, 6]
    let tilesCopy = [...tiles];
    for (let i = 0; i < 7; i++) {
      let loc = Math.floor(Math.random() * indices.length);
      tilesCopy[i] = tiles[indices[loc]];
      indices = indices.filter(value => value !== indices[loc]);
    }
    setTiles(tilesCopy);
  }

  /**
   * Reset tiles on the board back into the hand.
   */
  const reset = () => {
    setTiles(prevTiles =>
      prevTiles.map(tile => ({ ...tile, position: 'ActionPanel' }))
    );
    // console.log("before: ", letterUpdates);
    setLetterUpdates({});
    console.log("after: ", letterUpdates);
  }

  /**
   * Refreshes all the tiles in the hand with a random set of 7 tiles
   */
  const refresh = () => {
    let url = baseURL + "refreshhand/" + gameID + "/";

    setTiles(prevTiles =>
      prevTiles.map(tile => ({ ...tile, position: 'ActionPanel' }))
    );
    setLetterUpdates({});

    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ playerName: playerName })
    })
      .then(response => response.json())
      .then(data => {
        // processing the server response
        console.log(data);
        setHand(data);
      })
      .catch(error => {
        alert(error);
        console.log("Error: ", error);
      })
  }

  /** 
   * Parses the board returned by the server.
  */
  function parseBoard(board) {

    let newpos = {};

    for (let i = 0; i < 15; i++) {
      for (let j = 0; j < 15; j++) {
        if (board[i][j] !== '') {
          let cellKey = `${i}-${j}`;
          newpos[cellKey] = board[i][j];
        }
      }
    }

    setScoredLetters(prevScoredLetters => ({
      ...prevScoredLetters,
      ...newpos
    }));

  };

  /**
   * Parses all the updates returned by the server.
   * If the move is valid, it updates the game state accordingly.
   * Otherwise it reverts all the moves.
   */
  function parseOwnUpdates(updates) {

    if (updates.valid) { // if move is valid, update the game state
      const updatesState = updates.gameState;
      parseBoard(updatesState.Board);
      setHand(updatesState.Players[playerName].hand);
      setTilebag(updatesState.LetterDistribution);
    }
    else { // else revert all the moves
      setTiles(prevTiles =>
        prevTiles.map(tile => ({ ...tile, position: 'ActionPanel' }))
      );
      alert(updates.message); // display the error message
    };

  };

  /**
   * Submits the current state of the game board.
   *
   * This function iterates over the `letterUpdates` object, which contains the updates to the letters on the board.
   * For each update, it extracts the location and the letter and adds them to the `data` array.
   * The location is split into x and y coordinates, which are converted to numbers.
   * The `data` array is then ready to be sent to the server or processed further.
   */
  const submit = () => {
    let data = []
    for (const [key, value] of Object.entries(letterUpdates)) {
      let locs = value[0].split("-");
      data.push({ letter: value[1], xLoc: Number(locs[1]), yLoc: Number(locs[0]) });
    };
    setLetterUpdates({});
    const url = baseURL + gameID + "/updategame/";
    // const data = JSON.stringify({ playerName: player, updates: tilePositions })
    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ playerName: playerName, updates: data })
    })
      .then(response => response.json())
      .then(data => {
        // processing the server response
        console.log(data);
        parseOwnUpdates(data);
      })
      .catch(error => {
        alert(error);
        console.log("Error: ", error);
      })
  };

  return (
    <div className='App'>
      <div className="board-score">
        <Board
          letterUpdates={letterUpdates}
          onTileDrop={handleTileDrop}
          scoredLetters={scoredLetters}
        />
        <Infoboard
          tilebag={tilebag}
          p1_score={0}
          p2_score={0}
        />
      </div>
      <ActionPanel
        tilesAp={tiles.map(tile => {
          if (tile.position === 'ActionPanel') {
            return <Tile key={tile.id} letter={tile.letter} id={tile.id} />;
          } else {
            return <div key={tile.id} className="tile-placeholder"></div>;
          }
        })}
        shuffle={shuffle}
        submit={submit}
        reset={reset}
        refresh={refresh}
      />
      <Rules isRulesOpen={isRulesOpen} setIsRulesOpen={setIsRulesOpen}></Rules>
    </div>
  );
};