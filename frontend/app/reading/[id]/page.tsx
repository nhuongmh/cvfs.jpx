"use client";
import React, { useState, useEffect } from 'react';
import './ielts.css';
import { Article, Question, ArticleReading } from '@/models/article';
import { useSearchParams, usePathname } from 'next/navigation'
import { siteConfig } from "@/config/site";

function ArticleTestPage() {
  const [article, setArticle] = useState<Article | null>(null);
  const [articleReading, setArticleReading] = useState<ArticleReading | null>(null);
  const [timeRemaining, setTimeRemaining] = useState(60);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const searchParams = useSearchParams();
    const pathname = usePathname();
    if (pathname === null || searchParams === null) {
        return null
    }
    const articleId = pathname.split('/').pop() || searchParams.get('id') || "";

  // Fetch article and questions from API
  useEffect(() => {
    const fetchArticle = async () => {
      try {
        const url = `${siteConfig.private_url_prefix}/ie/article/${articleId}`;
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Failed to fetch article');
        }
        const data = await response.json();
        setArticle(data);
      } catch (err) {
        setError('Error fetching article: ' + (err as Error).message);
        console.error('Error fetching article:', err);
      }
    };

    const fetchQuestions = async () => {
      try {
        const url = `${siteConfig.private_url_prefix}/ie/article/${articleId}/reading`;
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Failed to fetch article reading');
        }
        const data = await response.json();
        setArticleReading(data);
      } catch (err) {
        setError('Error fetching questions: ' + (err as Error).message);
        console.error('Error fetching questions:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchArticle();
    fetchQuestions();
  }, []);

  const handleAnswerSelection = (questionIndex: number, answer: number) => {
    if (!articleReading) {
      console.error('Article reading data is not available');
      return;
    }
    if (questionIndex < 0 || questionIndex >= articleReading?.questions.length) return;
    if (answer < 0 || answer >= articleReading.questions[questionIndex].options.length) return;
    articleReading.questions[questionIndex].user_answer_option = answer;
  };

  const handleTextAnswerChange = (questionIndex: number, answer: string) => {
    if (!articleReading) {
      console.error('Article reading data is not available');
      return;
    }
    if (questionIndex < 0 || questionIndex >= articleReading.questions.length) return;
    articleReading.questions[questionIndex].user_answer_str = answer;
  };

  const handleSubmit = () => {
    if (!articleReading) {
      console.error('Article reading data is not available');
      return;
    }
    alert('Test submitted!');
    console.log('Answers:', articleReading.questions);
  };

  if (loading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Loading test content...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-container">
        <h2>Error</h2>
        <p>{error}</p>
        <button onClick={() => window.location.reload()}>Try Again</button>
      </div>
    );
  }

  // Function to format article content with proper paragraphs
  const formatContent = (content: string) => {
    if (!content) return [];
    return content.split('\n').map((paragraph, index) => 
      <p key={index}>{paragraph}</p>
    );
  };

  // Render the appropriate question type
  const renderQuestion = (question: Question, index: number) => {
    switch(question.type) {
      case 'multiple_choice':
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <div className="options">
              {question.options.map((option, optionIndex) => (
                <div className="option" key={optionIndex}>
                  <input 
                    type="radio" 
                    id={`q${index}_${optionIndex}`} 
                    name={`question${index}`} 
                    value={option}
                    checked={question.user_answer_option === optionIndex}
                    onChange={() => handleAnswerSelection(index, optionIndex)}
                  />
                  <label htmlFor={`q${index}_${optionIndex}`}>
                    <span className="option-letter">{String.fromCharCode(65 + optionIndex)}</span> {option}
                  </label>
                </div>
              ))}
            </div>
          </div>
        );
      
      case 'short_answer':
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <div className="short-answer">
              <input 
                type="text" 
                placeholder="Enter your answer"
                value={question.user_answer_str || ''}
                onChange={(e) => handleTextAnswerChange(index, e.target.value)}
              />
            </div>
          </div>
        );
      
      case 'true_false_not_given':
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <div className="options">
              {['TRUE', 'FALSE', 'NOT GIVEN'].map((option, optionIndex) => (
                <div className="option" key={optionIndex}>
                  <input 
                    type="radio" 
                    id={`q${index}_${option}`} 
                    name={`question${index}`} 
                    value={option}
                    checked={question.user_answer_option === optionIndex}
                    onChange={() => handleAnswerSelection(index, optionIndex)}
                  />
                  <label htmlFor={`q${index}_${option}`}>
                    <span className="option-letter">{option}</span>
                  </label>
                </div>
              ))}
            </div>
          </div>
        );
      
      case 'matching_headings':
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <div className="paragraph-container">
              <blockquote>{question.paragraph}</blockquote>
            </div>
            <div className="options">
              {question.headings.map((heading, headingIndex) => (
                <div className="option" key={headingIndex}>
                  <input 
                    type="radio" 
                    id={`q${index}_${headingIndex}`} 
                    name={`question${index}`} 
                    value={heading}
                    checked={question.user_answer_option == headingIndex}
                    onChange={() => handleAnswerSelection(index, headingIndex)}
                  />
                  <label htmlFor={`q${index}_${headingIndex}`}>
                    <span className="option-letter">{heading}</span>
                  </label>
                </div>
              ))}
            </div>
          </div>
        );
      
      default:
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <p className="error-message">Unknown question type</p>
          </div>
        );
    }
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <div className="logo">
          <img src="/Dharma_Initiative_logo.svg.png" alt="IELTS Online Tests" />
        </div>
        <div className="timer">
          <span>{timeRemaining} minutes remaining</span>
        </div>
        <div className="controls">
          <button className="btn-fullscreen">
            <i className="fas fa-expand"></i>
          </button>
          <button className="btn-submit" onClick={handleSubmit}>
            Submit <i className="fas fa-arrow-right"></i>
          </button>
        </div>
      </header>

      <div className="main-content">
        <div className="reading-section">

          <div className="passage-container">
            {article && (
              <>
          <div className="passage-header">
            <div className="passage-logo">
              <img src="/ielts-logo-small.png" alt="Reading Comprehensive Test" />
            </div>
            <div className="passage-title">
              <h2>{article.title}</h2>
            </div>
          </div>

          <div className="passage-content">
            <h2>{article.title}</h2>
            {article.image && (
              <div className="passage-image">
                <img src={article.image} alt={article.title} />
              </div>
            )}
            <div className="passage-text">
              {formatContent(article.content)}
            </div>
          </div>
              </>
            )}
          </div>

        </div>

        <div className="questions-section">
          <div className="questions-header">
            <h2>Questions 1-{articleReading?.questions.length}</h2>
            <p>Answer the questions according to the instructions.</p>
          </div>

          <div className="questions-container">
          {articleReading?.questions.map((question, index) => (
              renderQuestion(question, index)
            ))}
          </div>

        </div>
      </div>
    </div>
  );
}

export default ArticleTestPage;