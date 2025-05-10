"use client";
import React, { useState, useEffect, ReactElement } from 'react';
import './ielts.css';
import { Article, Question, ArticleReading } from '@/models/article';
import { useSearchParams, usePathname } from 'next/navigation'
import { siteConfig } from "@/config/site";
import { Button, ButtonGroup } from "@heroui/button";
import { RadioGroup, Radio } from "@heroui/radio";
import { Input, Form } from "@heroui/react";

function ArticleTestPage() {
  const [article, setArticle] = useState<Article | null>(null);
  const [articleReading, setArticleReading] = useState<ArticleReading | null>(null);
  const [timeRemaining, setTimeRemaining] = useState(60);
  const [loading, setLoading] = useState(true);
  const [textAnswer, setTextAnswer] = useState<string>('');
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

  const handleAnswerChange = (questionId: number, answer: string) => {
    if (!articleReading) {
      console.error('Article reading data is not available');
      return;
    }
    const question = articleReading.questions.find(q => q.id === questionId);
    if (question) {
      question.user_answer_str = answer;
    }
  };

  const handleSubmit = () => {

    if (!articleReading || !articleReading.questions) {
      console.error('Article reading data is not available');
      return;
    }

    for (const question of articleReading.questions) {
      if (!question.user_answer_str || question.user_answer_str.trim() === '') {
        alert(`Please answer all questions before submitting.`);
        return;
      }
    }

    let answers = new Map<number, string>();
    for (const question of articleReading.questions) {
      answers.set(question.id, question.user_answer_str);
    }
    console.log('Answers:', answers);
    console.log(JSON.stringify(Object.fromEntries(answers)))
    // const answers = articleReading.questions.reduce((acc, question) => {
    //   acc[question.id] = question.user_answer_str || '';

    //   return acc;
    // }, {} as Record<number, string>);

    const submitUrl = `${siteConfig.private_url_prefix}/ie/article/reading/${articleReading.id}/submit`;

    fetch(submitUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(Object.fromEntries(answers)),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Failed to submit answers');
        }
        return response.json();
      })
      .then((data) => {
        console.log('Submission successful:', data);

        // Update the articleReading object with the results
        if (articleReading) {
          articleReading.questions.forEach((question) => {
            const result = data.question_results.find(
              (res: any) => res.question_id === question.id
            );
            if (result) {
              question.correct = result.correct;
              question.correct_answer = result.answer;
            }
          });
        }

        // Display the final score
        alert(`Test submitted successfully! Your score: ${data.score}`);
        articleReading.score = data.score;
        setArticleReading({ ...articleReading }); // Trigger re-render
      })
      .catch((err) => {
        console.error('Error submitting answers:', err);
        alert('Error submitting answers. Please try again.');
      });

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
    const renderStatus = () => {
      if (question.correct !== undefined) {

        if (question.correct) {
          return <p className="pt-6  text-green-500 font-bold">Correct</p>;
        } else if (question.correct_answer) {
          return (
            <div className="pt-6 font-bold">
              <p className='text-red-500'>Incorrect</p>
              <p>Answer: {question.correct_answer}</p>
            </div>
          );
        }
      }
      return null;
    };

    switch (question.type) {
      case 'multiple_choice':
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <div className="options">
              <RadioGroup
                value={question.user_answer_str}
                onValueChange={(value) => handleAnswerChange(question.id, value)}
              >
                {question.options.map((option, optionIndex) => (
                  <Radio key={optionIndex} value={String.fromCharCode(65 + optionIndex)}>
                    <span className="option-letter">{String.fromCharCode(65 + optionIndex)}</span> {option}
                  </Radio>
                ))}
              </RadioGroup>
            </div>
            {renderStatus()}
          </div>
        );

      case 'short_answer':
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <div className="short-answer">
              <Input value={question.user_answer_str} onValueChange={(value) => handleAnswerChange(question.id, value)} />
            </div>
            {renderStatus()}
          </div>
        );

      case 'true_false_not_given':
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <div className="options">
              <RadioGroup
                value={question.user_answer_str}
                onValueChange={(value) => handleAnswerChange(question.id, value)}
              >
                {['TRUE', 'FALSE', 'NOT GIVEN'].map((option, optionIndex) => (
                  <Radio key={optionIndex} value={option}>
                    <span className="option-letter">{option}</span>
                  </Radio>
                ))}
              </RadioGroup>
            </div>
            {renderStatus()}
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
                    checked={question.user_answer_str == heading}
                    onChange={() => handleAnswerChange(question.id, heading)}
                  />
                  <label htmlFor={`q${index}_${headingIndex}`}>
                    <span className="option-letter">{heading}</span>
                  </label>
                </div>
              ))}
            </div>
            {renderStatus()}
          </div>
        );

      default:
        return (
          <div className="question" key={index}>
            <p className="question-number">{index + 1}. {question.question}</p>
            <p className="error-message">Unknown question type</p>
            {renderStatus()}
          </div>
        );
    }
  };

  return (
    <div className="app-container">
      <div className="app-header">
        <div className="logo">
          <img src="/images/Dharma_Initiative_logo.svg.png" alt="IELTS Online Tests" />
        </div>
        <div className="timer">
          <span>{timeRemaining} minutes remaining</span>
        </div>
        <div className="controls">
          <Button color="primary" onPress={handleSubmit}>
            Submit <i className="fas fa-arrow-right"></i>
          </Button>
        </div>
      </div>

      <div className="main-content">
        <div className="reading-section">

          <div className="passage-container">
            {article && (
              <>
                <div className="passage-header">
                  <div className="passage-logo">
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
          {articleReading?.score !== undefined && (
            <div className="score-display">
              <h3>Your Score: {articleReading.score}</h3>
            </div>
          )}
          <div className="questions-header">
            <h2>Questions 1-{articleReading?.questions.length}</h2>
            <p>Answer the questions according to the instructions.</p>
          </div>

          <div className="questions-container">
            <Form className="questions-form">
            {articleReading?.questions.map((question, index) => (
              renderQuestion(question, index)
            ))}
            </Form>
          </div>

        </div>
      </div>
    </div>
  );
}

export default ArticleTestPage;