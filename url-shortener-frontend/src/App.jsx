import React, { useState } from "react";

function App() {
  const [url, setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setShortUrl("");

    try {
      const response = await fetch(
        "https://urlshortenerwith-go-production.up.railway.app/shorten",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ url }) // using state variable
        }
      );

      if (!response.ok) {
        const errData = await response.json();
        setError(errData.error || "Something went wrong");
        return;
      }

      const data = await response.json();
      setShortUrl(data.short_url); // match backend key
    } catch (err) {
      console.error("Error:", err);
      setError("Failed to fetch. Check backend server.");
    }
  };

  return (
    <div className="container">
      <h1 className="title">Go URL Shortener</h1>

      <form onSubmit={handleSubmit} className="form">
        <input
          type="text"
          placeholder="Enter your long URL..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
        />
        <button type="submit">Shorten</button>
      </form>

      {error && <p className="error">❌ {error}</p>}

      {shortUrl && (
        <div className="result">
          <p>✅ Shortened URL:</p>
          <a href={shortUrl} target="_blank" rel="noopener noreferrer">
            {shortUrl}
          </a>
        </div>
      )}
    </div>
  );
}

export default App;
