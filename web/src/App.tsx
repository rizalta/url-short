import { useState } from "react";

const App = () => {
  const [url, setUrl] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.SubmitEvent) => {
    e.preventDefault();
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
    setCode(data.short_url)
  }

  return (
    <div>
      <h1>URL SHORT</h1>
      <form onSubmit={handleSubmit}>
        <input type="text" placeholder="Enter URL" value={url} onChange={(e) => setUrl(e.target.value)} />
        <button type="submit">Shorten</button>
      </form>

      {error ? <p>{error}</p> : null}
      {code ? <p>Short URL: {`${window.location.origin}/${code}`}</p> : null}
    </div >
  );
}

export default App
