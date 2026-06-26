import { useEffect, useState } from 'react'
import './App.css'

function Header({onAdd}) {
  return (
  <div className="header">
    <h1>Travel Log</h1>
    <p>A record of my trip to edinburgh</p>
    <button className="button" onClick={onAdd}>
      Add New Entry
    </button>
  </div>
  
  );
}

function EntryForm({ title, formData, setFormData, onSave, onCancel }) {
  return (
    <div className="edit-form">
      <h3>{title}</h3>

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
        value={new Date(formData.visited_at).toISOString().slice(0, 16)}
        onChange={(e) =>
          setFormData({ ...formData, visited_at: e.target.value })
        }
      />

      <div className="button-row">
        <button onClick={onSave}>Save</button>
        <button onClick={onCancel}>Cancel</button>
      </div>
    </div>
  );
}

function App() {
  const [entries, setEntries] = useState([])
  const [editingEntry, setEditingEntry] = useState(null);
  const [addingEntry, setAddingEntry] = useState(false);
  const [formData, setFormData] = useState({
    place: "",
    comment: "",
    visited_at: ""
  });

  function startEdit(entry) {
    setAddingEntry(false)

    setFormData({
      place: entry.place,
      comment: entry.comment,
      visited_at: entry.visited_at
    });
    
    setEditingEntry(entry);
  }

  function startCreate() {
    setEditingEntry(null);

    setFormData({
      place: "",
      comment: "",
      visited_at: new Date().toISOString()
    });
    
    setAddingEntry(true);
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

  function saveNewEntry() {
    fetch(`http://81.100.84.76:8080/api/entries`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(formData).padEnd(3, "00Z")
    })
      .then(res => res.json())
      .then(() => {
        setAddingEntry(false);
        loadEntries(); // refresh list
      });
  }

  useEffect(() => {
      loadEntries()
    }, [] );

  return (
    <>
      <div className="container">
        <Header onAdd={startCreate}/>
        {editingEntry && (
          <EntryForm
            title="Edit Entry"
            formData={formData}
            setFormData={setFormData}
            onSave={saveEdit}
            onCancel={() => setEditingEntry(null)}
          />
        )}

        {addingEntry && (
          <EntryForm
            title="New Entry"
            formData={formData}
            setFormData={setFormData}
            onSave={saveNewEntry}
            onCancel={() => setAddingEntry(false)}
          />
        )}

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
