import { useState } from "react";

const App = () => {
  const [url, setUrl] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState("");

  const baseUrl = window.location.origin;

  const handleSubmit = async (e: React.SubmitEvent) => {
    e.preventDefault();

    setError("");
    setCode("");

    const res = await fetch("/api/shorten", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ url })
    });

    if (!res.ok) {
      setError("Failed to shorten URL");
      return;
    }

    const data = await res.json();

    setUrl("");
    setCode(data.code);
  }

  return (
    <div>
      <h1>URL SHORT</h1>
      <form onSubmit={handleSubmit}>
        <input type="text" placeholder="Enter URL" value={url} onChange={(e) => setUrl(e.target.value)} />
        <button type="submit">Shorten</button>
      </form>

      {error ? <p>{error}</p> : null}
      {code ? <p>Short URL: <a href={`${baseUrl}/${code}`} target="_blank" rel="noopener noreferrer">{`${baseUrl}/${code}`}</a></p> : null}
    </div >
  );
}

export default App
