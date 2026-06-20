import { useEffect, useState } from 'react'
import './App.css'

function Header() {
  return <h1>Travel Log</h1>;
}



function App() {
  const [entries, setEntries] = useState([])

  useEffect(() => {
      fetch("http://81.100.84.76:8080/api/entries")
      .then(response => response.json())
      .then(data => setEntries(data))
      .catch(error => console.error(error));
    }, [] );

  return (
    <>
      <Header/>
      {entries.map(entry => (
        <div key={entry.id}>
          <h3>{entry.place}</h3>
          <p>{entry.comment}</p>
        </div>
      ))}
    </>
  )
}

export default App
