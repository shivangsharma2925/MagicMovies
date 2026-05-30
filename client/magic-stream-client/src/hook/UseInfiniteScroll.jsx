import { useEffect, useRef } from "react";

const useInfiniteScroll = (
  fetchNextPage,
  hasNextPage
) => {
  const observerRef = useRef();

  useEffect(() => {
    const currentTarget = observerRef.current;

    const observer = new IntersectionObserver(
      (entries) => {
        if (
          entries[0].isIntersecting &&
          hasNextPage
        ) {
          fetchNextPage();
        }
      },
      {
        rootMargin: "100px",
        threshold: 0.1,
      }
    );

    if (currentTarget) {
      observer.observe(currentTarget);
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget);
      }
    };
  }, [fetchNextPage, hasNextPage]);

  return observerRef;
};

export default useInfiniteScroll;