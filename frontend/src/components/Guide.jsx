import React, { useState, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism';
import remarkGfm from 'remark-gfm';

const components = {
  code({ node, inline, className, children, ...props }) {
    const match = /language-(\w+)/.exec(className || '');
    return !inline && match ? (
      <SyntaxHighlighter
        language={match[1]}
        PreTag="div"
        {...props}
      >
        {String(children).replace(/\n$/, '')}
      </SyntaxHighlighter>
    ) : (
      <code className={className} {...props}>
        {children}
      </code>
    );
  },
  h1: props => <h1 className="text-3xl font-bold my-4" {...props} />,
  h2: props => <h2 className="text-2xl font-semibold my-3" {...props} />,
  h3: props => <h3 className="text-xl font-semibold my-2" {...props} />,
  h4: props => <h4 className="text-lg font-semibold my-2" {...props} />,
  p: props => <p className="my-2" {...props} />,
  ul: props => <ul className="list-disc list-inside my-2" {...props} />,
  ol: props => <ol className="list-decimal list-inside my-2" {...props} />,
};

const Guide = ({ guideName }) => {
  const [content, setContent] = useState('');
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchGuide = async () => {
      try {
        const response = await fetch(`/guides/${guideName}.mdx`);
        if (!response.ok) {
          throw new Error('Failed to fetch guide');
        }
        const text = await response.text();
        setContent(text);
      } catch (err) {
        setError('Failed to load guide');
        console.error(err);
      }
    };

    fetchGuide();
  }, [guideName]);

  if (error) {
    return (
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong className="font-bold">Error: </strong>
        <span className="block sm:inline">{error}</span>
      </div>
    );
  }

  if (!content) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <article className="prose prose-sm sm:prose lg:prose-lg xl:prose-xl">
        <ReactMarkdown 
          components={components}
          remarkPlugins={[remarkGfm]}
        >
          {content}
        </ReactMarkdown>
      </article>
    </div>
  );
};

export default Guide;
