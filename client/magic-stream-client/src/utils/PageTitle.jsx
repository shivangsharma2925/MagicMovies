import { useEffect } from "react";

const DEFAULT_TITLE =
  "Magic Stream — Discover Movies and Watch Trailers";

const PageTitle = ({ title }) => {
  useEffect(() => {
    document.title = title
      ? `${title} — Magic Stream`
      : DEFAULT_TITLE;

    return () => {
      document.title = DEFAULT_TITLE;
    };
  }, [title]);

  return null;
};

export default PageTitle;