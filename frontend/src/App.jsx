import { useEffect, useState } from 'react'
import './App.css'



function Header() {
  return (
  <div className="header">
    <h1>Travel Log</h1>
    <p>A record of my trip to edinburgh</p>
  </div>
  );
}

function App() {
  const [entries, setEntries] = useState([])
  const [editingEntry, setEditingEntry] = useState(null);
  const [formData, setFormData] = useState({
    place: "",
    comment: "",
    visited_at: ""
  });

  function startEdit(entry) {
    setEditingEntry(entry);

    setFormData({
      place: entry.place,
      comment: entry.comment,
      visited_at: entry.visited_at
    });
  }

  function deleteEntry(id) {
    fetch(`http://81.100.84.76:8080/api/entries/${id}`, {
      method: "DELETE"
    })
    .then(() => loadEntries());
  }

  function loadEntries () {
    fetch("http://81.100.84.76:8080/api/entries")
        .then(response => response.json())
        .then(data => setEntries(data))
        .catch(error => console.error(error));
  }

  function saveEdit() {
    fetch(`http://81.100.84.76:8080/api/entries/${editingEntry.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(formData)
    })
      .then(res => res.json())
      .then(() => {
        setEditingEntry(null);
        loadEntries(); // refresh list
      });
  }

  useEffect(() => {
      loadEntries()
    }, [] );

  return (
    <>
      <div className="container">
        <Header/>

        {editingEntry && (
          <div className="edit-form">
            <h3>Edit Entry</h3>

            <input
              value={formData.place}
              onChange={(e) =>
                setFormData({ ...formData, place: e.target.value })
              }
            />

            <input
              value={formData.comment}
              onChange={(e) =>
                setFormData({ ...formData, comment: e.target.value })
              }
            />

            <input
              type="datetime-local"
              value={new Date(formData.visited_at).toISOString().slice(0,16)}
              onChange={(e) =>
                setFormData({ ...formData, visited_at: e.target.value })
              }
            />
            <div className="button-row">
              <button onClick={saveEdit}>Save</button>
              <button onClick={() => setEditingEntry(null)}>Cancel</button>
            </div>
          </div>)}

        <div className="entries-grid">
          {entries.map(entry => (
            <div className="entry-card" key={entry.id}>
              <h3>{entry.place}</h3>
              <small>
                {new Date(entry.visited_at).toLocaleString("en-GB")}
              </small>
              <p>{entry.comment}</p>

              <div className="button-row">
                <button className="button" onClick={() => startEdit(entry)}>
                  Edit
                </button>
                <button className="button" onClick={() => deleteEntry(entry.id)}>
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>

      </div>
    </>
  )
}

export default App
