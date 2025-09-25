import { useState } from "react";

export default function ShortenForm() {
  const [url, setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    try {
      const res = await fetch("http://localhost:8080/shorten", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });
      const data = await res.json();
      setShortUrl(data.short_url);
      setUrl("");
    } catch (err) {
      alert("Failed to shorten URL");
    }
  }

  return (
    <div className="card">
      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <input
            className="input"
            type="url"
            placeholder="Enter a long URL..."
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            required
          />
          <button className="btn btn-primary" type="submit">Shorten</button>
        </div>
      </form>

      {shortUrl && (
        <div className="short-box mt-4">
          <a className="short-link" href={shortUrl} target="_blank" rel="noreferrer">{shortUrl}</a>
          <button className="btn btn-secondary small" onClick={() => navigator.clipboard.writeText(shortUrl)}>Copy</button>
        </div>
      )}
    </div>
  );
}
