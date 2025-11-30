import { useEffect, useMemo, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { toast } from "sonner";
import { ListDocs, ReadDoc } from "../../wailsjs/go/main/App";
import type { config } from "../../wailsjs/go/models";
import "./Docs.css";

type DocsProvider = config.DocsProvider;
type DocsPage = config.DocsPage;

export default function Docs() {
  const [providers, setProviders] = useState<DocsProvider[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedProvider, setSelectedProvider] = useState<string>("");
  const [selectedPage, setSelectedPage] = useState<DocsPage | null>(null);
  const [content, setContent] = useState<string>("");
  const [reading, setReading] = useState(false);

  useEffect(() => {
    const load = async () => {
      try {
        const data = await ListDocs();
        setProviders(data || []);
        if (data && data.length > 0) {
          setSelectedProvider(data[0].provider);
        }
      } catch (err) {
        console.error(err);
        toast.error("Failed to load docs index");
      } finally {
        setLoading(false);
      }
    };
    void load();
  }, []);

  const pagesForProvider = useMemo(() => {
    const prov = providers.find((p) => p.provider === selectedProvider);
    return prov?.pages ?? [];
  }, [providers, selectedProvider]);

  const handleSelectPage = async (page: DocsPage) => {
    setSelectedPage(page);
    setContent("");
    setReading(true);
    try {
      const text = await ReadDoc(page.provider, page.slug, "md");
      setContent(text);
    } catch (err) {
      console.error(err);
      toast.error("Failed to load document");
    } finally {
      setReading(false);
    }
  };

  if (loading) {
    return (
      <div className="docs-container">
        <div className="docs-loading">Loading docs index…</div>
      </div>
    );
  }

  if (providers.length === 0) {
    return (
      <div className="docs-container">
        <div className="docs-empty">
          <h2>No docs synced yet</h2>
          <p>Run scripts/sync-docs.sh to fetch provider documentation.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="docs-container">
      <div className="docs-sidebar">
        <div className="docs-provider-header">
          <span>Provider</span>
        </div>
        <select
          className="docs-provider-select"
          value={selectedProvider}
          onChange={(e) => {
            setSelectedProvider(e.target.value);
            setSelectedPage(null);
            setContent("");
          }}
        >
          {providers.map((p) => (
            <option key={p.provider} value={p.provider}>
              {p.provider} ({p.date})
            </option>
          ))}
        </select>

        <div className="docs-page-list">
          {pagesForProvider.map((p) => (
            <button
              key={`${p.provider}-${p.slug}`}
              type="button"
              className={`docs-page-item ${
                selectedPage?.provider === p.provider &&
                selectedPage?.slug === p.slug
                  ? "active"
                  : ""
              }`}
              onClick={() => handleSelectPage(p)}
            >
              {p.title || p.slug}
            </button>
          ))}
        </div>
      </div>

      <div className="docs-content">
        {reading && <div className="docs-loading">Loading…</div>}
        {!reading && !selectedPage && (
          <div className="docs-placeholder">
            <h2>Select a doc</h2>
            <p>Choose a provider and document from the list.</p>
          </div>
        )}
        {!reading && selectedPage && content && (
          <div className="docs-markdown">
            {content.trim().startsWith("<") ? (
              <div
                className="docs-html"
                dangerouslySetInnerHTML={{ __html: content }}
              />
            ) : (
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {content}
              </ReactMarkdown>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
