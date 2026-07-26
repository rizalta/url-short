import { useState } from "react";
import "./index.css"

const App = () => {
  const [url, setUrl] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  const baseUrl = window.location.origin;
  const shortUrl = code ? `${baseUrl}/${code}` : "";

  const handleSubmit = async (e: React.SubmitEvent) => {
    e.preventDefault();
    if (!url.trim()) return;

    setError("");
    setCode("");
    setLoading(true);

    try {
      const res = await fetch("/api/shorten", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ url: url.trim() })
      });

      const data = await res.json();

      if (!res.ok) {
        setError(data.error || "Something went wrong");
        return;
      }

      setCode(data.code);
      setUrl("");
    } catch {
      setError("Something went wrong")
    } finally {
      setLoading(false);
    }
  }

  const handleCopy = () => {
    if (!shortUrl) return;
    navigator.clipboard.writeText(shortUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <div className="container">
      <div className="header">
        <h1>URL-SHORTZ</h1>
        <p>Shorten ur urlz</p>
      </div>
      <form onSubmit={handleSubmit} className="form-group">
        <input
          type="text"
          className="input-field"
          placeholder="Enter ur long url..."
          value={url} onChange={(e) => setUrl(e.target.value)} />
        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? "..." : "Shorten"}
        </button>
      </form>

      {error && (
        <div className="error-message">{error}</div>
      )}

      {shortUrl && (
        <div className="result-card">
          <a
            href={shortUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="result-link"
          >
            {shortUrl}
          </a>
          <button
            onClick={handleCopy}
            className={`btn-copy ${copied ? "copied" : ""}`}
          >
            {copied ? "✓ Copied!" : "📋 Copy"}
          </button>
        </div>
      )}
    </div >
  );
}

export default App
