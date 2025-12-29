// frontend/src/components/ErrorFallback.tsx
import React from 'react';
import './ErrorFallback.css';

interface ErrorFallbackProps {
  error: Error;
  onRetry: () => void;
}

const ErrorFallback: React.FC<ErrorFallbackProps> = ({ error, onRetry }) => {
  return (
    <div className="error-fallback">
      <h2>Something went wrong.</h2>
      <p>We're sorry, but an unexpected error occurred.</p>
      <details>
        <summary>Error Details</summary>
        <pre>{error.message}</pre>
      </details>
      <button onClick={onRetry}>Try Again</button>
    </div>
  );
};

export default ErrorFallback;
