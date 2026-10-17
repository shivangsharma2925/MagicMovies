import { useState } from "react";
import { Modal } from "react-bootstrap";
import useAxiosPrivate from "../hook/UseAxiosPrivate";
import "../styling/AI_related.css";

const AskAiBtn = () => {
  const EXAMPLES = [
    "Give me details for the Fast & Furious series",
    "Movie Info of John Wick Chapter 2",
    "All Marvel Spider-Man movies",
    "Christopher Nolan's top 5 movies",
  ];

  const [show, setShow] = useState(false);
  const [prompt, setPrompt] = useState("");
  const [loading, setLoading] = useState(false);
  const [showExamples, setShowExamples] = useState(true);
  const [result, setResult] = useState(null); // { ids: [{}] } | { error: string }

  const axiosPrivate = useAxiosPrivate();

  const handleClose = () => {
    setShow(false);
    setPrompt("");
    setShowExamples(true);
    setResult(null);
  };

  const handleGenerate = async () => {
    if (!prompt.trim()) return;
    setLoading(true);
    setResult(null);

    try {
      const res = await axiosPrivate.post("/ai/imdb", {
        prompt: prompt.trim(),
      });

      const data = res.data;

      if (data.error || !data.ids?.length) {
        setResult({ error: true });
      } else {
        setResult({ ids: data.ids });
        setShowExamples(false);
      }
    } catch {
      setResult({ error: true });
    } finally {
      setLoading(false);
    }
  };

  const handleUseIds = async () => {
    try {
      const movieDetails = result.ids.map((movie) => ({
        movie_name: movie.movie_name,
        movie_year: String(movie.year ?? ""),
      }));

      await axiosPrivate.post("/addmovieai", {
        movie_details: movieDetails,
      });

      handleClose();
    } catch (error) {
      setResult((prev) => ({
        ...prev,
        addError: true,
      }));
      console.log(error);
    } finally {
      setLoading(false);
      setShowExamples(true);
    }
  };

  return (
    <>
      <button
        className="btn-ask-ai"
        onClick={(e) => {
          e.preventDefault();
          setShow(true);
        }}
      >
        <span className="ask-ai-icon">✦</span> Ask AI
      </button>

      <Modal show={show} onHide={handleClose} centered className="ai-modal">
        <Modal.Body className="ai-modal-body">
          {/* Header */}
          <div className="ai-modal-header">
            <div className="ai-modal-title">
              <span className="ask-ai-icon">✦</span> AI Movie(s) Finder
            </div>
            <button className="ai-modal-close" onClick={handleClose}>
              ×
            </button>
          </div>
          <p className="ai-modal-sub">
            Describe the movies you want and AI will find the movie details for
            you.
          </p>

          {/* Examples */}
          {showExamples && (
            <div id="example-section">
              <p className="ai-section-label">Try an example</p>
              <div className="ai-examples">
                {EXAMPLES.map((ex) => (
                  <button
                    key={ex}
                    className="ai-example-chip"
                    onClick={() => setPrompt(ex)}
                  >
                    {ex}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Prompt input */}
          <textarea
            className="ai-prompt-input"
            rows={3}
            placeholder="e.g. Give me details for the Mission Impossible series..."
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            onKeyDown={(e) =>
              e.key === "Enter" &&
              !e.shiftKey &&
              (e.preventDefault(), handleGenerate())
            }
          />

          {/* Generate button */}
          <button
            className="btn-ai-generate"
            onClick={handleGenerate}
            disabled={loading || !prompt.trim()}
          >
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                />
                Thinking...
              </>
            ) : (
              <>
                <span className="ask-ai-icon">✦</span> Bring Details
              </>
            )}
          </button>

          {/* Result */}
          {result?.ids && (
            <div className="ai-result-box">
              <p className="ai-section-label">
                Found {result.ids.length} movie
                {result.ids.length > 1 ? "s" : ""}
              </p>
              <div className="ai-movie-list">
                {result.ids.map((id) => (
                  <div key={id.imdbid} className="ai-movie-row">
                    <div className="ai-movie-info">
                      <span className="ai-movie-name">
                        {id.movie_name || "Unknown"}
                      </span>
                      <span className="ai-movie-year">{id.year}</span>
                    </div>
                    <span className="ai-id-pill">
                      <a
                        href={`https://www.google.com/search?q=site%3Aimdb.com%2Ftitle%2F+${id.movie_name}+${id.year}&btnI=1`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="ai-verify-link"
                      >
                        Verify?
                      </a>
                    </span>
                  </div>
                ))}
                <div className="ai-context-box">
                  Not Satisfied ? Try adding more Context like
                  <em> Year, Genre, Cast, etc</em>
                </div>
              </div>
              <button className="btn-ai-use" onClick={handleUseIds}>
                Add these movies →
              </button>
            </div>
          )}

          {result?.error && (
            <div className="ai-error-box">
              Uh oh — couldn't find valid IMDB IDs for that. Try rephrasing,
              e.g. <em>"IMDB IDs for the Avengers movies"</em>.
            </div>
          )}
        </Modal.Body>
      </Modal>
    </>
  );
};

export default AskAiBtn;
